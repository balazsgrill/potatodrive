package filesystem

import (
	"crypto/md5"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/rs/zerolog"

	"github.com/balazsgrill/potatodrive/core"
	"github.com/balazsgrill/potatodrive/core/cfapi"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
	"golang.org/x/sys/windows"
)

type VirtualizationInstance struct {
	zerolog.Logger
	rootPath         string
	shortprefix      string
	longprefix       string
	fs               afero.Fs
	remoteCacheState core.RemoteStateCache

	connectionKey cfapi.CF_CONNECTION_KEY
	lock          sync.Mutex
	watcher       *fsnotify.Watcher
	callbacks     core.FileStateCallbacks
}

// SetStateCallbacks implements core.Virtualization.
func (instance *VirtualizationInstance) SetStateCallbacks(callbacks core.FileStateCallbacks) {
	instance.callbacks = callbacks
}

func StartProjecting(rootPath string, filesystem afero.Fs, logger zerolog.Logger) (core.Virtualization, error) {
	instance := &VirtualizationInstance{
		Logger:           logger,
		rootPath:         rootPath,
		fs:               filesystem,
		remoteCacheState: core.HashFilesRemotely(filesystem),
	}

	instance.longprefix = core.ToLongPath(rootPath)
	instance.shortprefix = core.ToShortPath(rootPath)

	return instance, instance.start()
}

func (instance *VirtualizationInstance) start() error {
	callbacks := &cfapi.Callbacks{
		FetchData: instance.fetchData,
		//FetchPlaceholders: instance.fetchPlaceholders, // using always_full
		//DeleteCompletion: instance.deleteCompletion,   // replaced by fswatch
	}

	instance.Logger.Print("Connecting sync root")
	hr := cfapi.CfConnectSyncRoot(core.GetPointer(instance.rootPath), callbacks.CreateCallbackTable(), uintptr(unsafe.Pointer(instance)), cfapi.CF_CONNECT_FLAG_REQUIRE_FULL_FILE_PATH, &instance.connectionKey)

	err := core.ErrorByCode(hr)
	if err != nil {
		return err
	}

	instance.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	instance.watcher.Add(instance.rootPath)
	go instance.watch()

	return nil
}

func getFileNameFromIdentity(info *cfapi.CF_CALLBACK_INFO) string {
	name := unsafe.Slice((*byte)(unsafe.Pointer(info.FileIdentity)), info.FileIdentityLength)
	return string(name)
}

func getPlaceholder(f fs.FileInfo) cfapi.CF_PLACEHOLDER_CREATE_INFO {
	var placeholder cfapi.CF_PLACEHOLDER_CREATE_INFO
	filename := f.Name()
	placeholder.RelativeFileName = core.GetPointer(filename)
	placeholder.FsMetadata.BasicInfo = toBasicInfo(f)
	identity := []byte(filename)
	placeholder.FileIdentity = uintptr(unsafe.Pointer(&identity[0]))
	placeholder.FileIdentityLength = uint32(len(identity))
	if !f.IsDir() {
		placeholder.FsMetadata.FileSize = int64(f.Size())
		placeholder.Flags = cfapi.CF_PLACEHOLDER_CREATE_FLAG_DISABLE_ON_DEMAND_POPULATION | cfapi.CF_PLACEHOLDER_CREATE_FLAG_MARK_IN_SYNC
	} else {
		placeholder.FsMetadata.FileSize = 0
	}
	return placeholder
}

func toBasicInfo(file fs.FileInfo) cfapi.FILE_BASIC_INFO {
	ftime := syscall.NsecToFiletime(file.ModTime().UnixNano())
	var attributes int32
	if file.IsDir() {
		attributes |= syscall.FILE_ATTRIBUTE_DIRECTORY
	} else {
		attributes |= syscall.FILE_ATTRIBUTE_NORMAL
	}
	return cfapi.FILE_BASIC_INFO{
		CreationTime:   ftime,
		LastAccessTime: ftime,
		LastWriteTime:  ftime,
		ChangeTime:     ftime,
		FileAttributes: attributes,
	}
}

func (instance *VirtualizationInstance) getOperationInfo(info *cfapi.CF_CALLBACK_INFO) cfapi.CF_OPERATION_INFO {
	operation := cfapi.CF_OPERATION_INFO{}
	operation.StructSize = uint32(unsafe.Sizeof(operation))
	operation.ConnectionKey = instance.connectionKey
	operation.TransferKey = info.TransferKey
	operation.CorrelationVector = info.CorrelationVector
	operation.RequestKey = info.RequestKey
	return operation
}

func (instance *VirtualizationInstance) transferPlaceholders(info *cfapi.CF_CALLBACK_INFO, parameters *cfapi.CF_OPERATION_PARAMETERS_TransferPlaceholders) uintptr {
	operation := instance.getOperationInfo(info)
	operation.Type = cfapi.CF_OPERATION_TYPE_TRANSFER_PLACEHOLDERS
	return cfapi.CfExecute(&operation, uintptr(unsafe.Pointer(parameters)))
}

func (instance *VirtualizationInstance) transferData(info *cfapi.CF_CALLBACK_INFO, parameters *cfapi.CF_OPERATION_PARAMETERS_TransferData) uintptr {
	operation := instance.getOperationInfo(info)
	operation.Type = cfapi.CF_OPERATION_TYPE_TRANSFER_DATA
	return cfapi.CfExecute(&operation, uintptr(unsafe.Pointer(parameters)))
}

func (instance *VirtualizationInstance) Close() error {
	if instance.connectionKey == 0 {
		return errors.New("not started")
	}

	instance.watcher.Close()
	hr := cfapi.CfDisconnectSyncRoot(instance.connectionKey)
	if hr != 0 {
		return core.ErrorByCode(hr)
	}

	return nil
}

func (instance *VirtualizationInstance) PerformSynchronization() error {
	err := instance.syncRemoteToLocal()
	if err != nil {
		return err
	}
	uploads, err := instance.syncLocalToRemote()
	for _, path := range uploads {
		localpath := instance.path_remoteToLocal(path)
		instance.FileUploading(localpath, 0)
		instance.Logger.Info().Msgf("Updating remote file '%s'", path)
		err = instance.streamLocalToRemote(path)
		if err != nil {
			instance.FileError(localpath, err)
			return err
		} else {
			instance.FileDone(localpath)
		}
	}
	return err
}

func (instance *VirtualizationInstance) setInSync(localpath string) error {
	instance.Logger.Info().Msgf("Set in-sync '%s'", localpath)
	placeholderstate, err := getPlaceholderState(localpath)
	if err != nil {
		return err
	}

	var handle syscall.Handle
	hr := cfapi.CfOpenFileWithOplock(core.GetPointer(localpath), cfapi.CF_OPEN_FILE_FLAG_WRITE_ACCESS|cfapi.CF_OPEN_FILE_FLAG_EXCLUSIVE, &handle)
	if hr != 0 {
		return core.ErrorByCode(hr)
	}
	defer cfapi.CfCloseHandle(handle)

	insync := (placeholderstate & cfapi.CF_PLACEHOLDER_STATE_IN_SYNC) != 0
	isaplacehoder := (placeholderstate & cfapi.CF_PLACEHOLDER_STATE_PLACEHOLDER) != 0

	if !isaplacehoder {
		fileinfo, err := os.Stat(localpath)
		if err != nil {
			return err
		}
		placeholder := getPlaceholder(fileinfo)

		// setting in-sync staate only works if it's a placeholder
		hr = cfapi.CfConvertToPlaceholder(handle, placeholder.FileIdentity, placeholder.FileIdentityLength, cfapi.CF_CONVERT_FLAG_NONE, 0, 0)
		if hr != 0 {
			return core.ErrorByCode(hr)
		}
	}
	if !insync {
		// updating a placeholder only works if it is marked as in-sync
		hr = cfapi.CfSetInSyncState(handle, cfapi.CF_IN_SYNC_STATE_IN_SYNC, cfapi.CF_SET_IN_SYNC_FLAG_NONE, nil)
		if hr != 0 {
			return core.ErrorByCode(hr)
		}
	}
	return nil
}

func (instance *VirtualizationInstance) localHash(remotepath string) ([]byte, error) {
	localpath := instance.path_remoteToLocal(remotepath)
	// only calculate hash if file is available on local disk
	localstate, err := getPlaceholderState(localpath)
	if err != nil {
		return nil, err
	}
	if (localstate | (cfapi.CF_PLACEHOLDER_STATE_IN_SYNC)) == 0 {
		return nil, nil
	}
	hash := md5.New()
	f, err := os.Open(localpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = io.Copy(hash, f)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func getPlaceholderInfo(localpath string) (*cfapi.CF_PLACEHOLDER_BASIC_INFO, error) {
	localpathstr, err := syscall.UTF16PtrFromString(localpath)
	if err != nil {
		return nil, err
	}
	fileHandle, err := syscall.CreateFile(localpathstr, syscall.GENERIC_READ, syscall.FILE_SHARE_READ, nil, syscall.OPEN_EXISTING, // existing file only
		syscall.FILE_ATTRIBUTE_NORMAL|syscall.FILE_FLAG_OVERLAPPED,
		0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(fileHandle)
	var placeholderInfo cfapi.CF_PLACEHOLDER_BASIC_INFO
	var ReturnedLength uint32
	hr := cfapi.CfGetPlaceholderInfo(fileHandle, cfapi.CF_PLACEHOLDER_INFO_BASIC, uintptr(unsafe.Pointer(&placeholderInfo)), uint32(unsafe.Sizeof(placeholderInfo)), &ReturnedLength)
	return &placeholderInfo, core.ErrorByCode(hr)
}

const FileAttributeTagInfo uint32 = 9

func getPlaceholderState(localpath string) (cfapi.CF_PLACEHOLDER_STATE, error) {
	localpathstr, err := syscall.UTF16PtrFromString(localpath)
	if err != nil {
		return cfapi.CF_PLACEHOLDER_STATE_INVALID, err
	}
	var finddata windows.Win32finddata
	findhandle, err := windows.FindFirstFile(localpathstr, &finddata)
	if err != nil {
		return cfapi.CF_PLACEHOLDER_STATE_INVALID, err
	}
	defer windows.FindClose(findhandle)

	result := cfapi.CfGetPlaceholderStateFromFindData(uintptr(unsafe.Pointer(&finddata)))
	return cfapi.CF_PLACEHOLDER_STATE(result), nil
}

func (instance *VirtualizationInstance) handleDeletion(localpath string) {
	instance.lock.Lock()
	defer instance.lock.Unlock()
	parentpath := filepath.Dir(localpath)
	remoteparent := instance.path_localToRemote(parentpath)
	remotepath := remoteparent + "/" + filepath.Base(localpath)
	remotepath = strings.TrimPrefix(remotepath, "/")
	err := instance.fs.RemoveAll(remotepath)
	if err != nil {
		instance.Logger.Printf("deleteCompletion: remove %s failed: %v", remotepath, err)
	}
}

func (instance *VirtualizationInstance) deleteCompletion(info *cfapi.CF_CALLBACK_INFO, data *cfapi.CF_CALLBACK_PARAMETERS_DeleteCompletion) uintptr {
	instance.lock.Lock()
	defer instance.lock.Unlock()
	filename := instance.callback_getRemoteFilePath(info)
	instance.Logger.Printf("deleteCompletion: %s", filename)
	//hashfilename := instance.path_hashFile(filename)

	err := instance.fs.Remove(filename)
	if err != nil {
		instance.Logger.Printf("deleteCompletion: remove %s failed: %v", filename, err)
	}
	/*
		err = instance.fs.Remove(hashfilename)
		if err != nil {
			instance.Logger.Printf("deleteCompletion: remove %s failed: %v", hashfilename, err)
		}*/
	return 0
}
