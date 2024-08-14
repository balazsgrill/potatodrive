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

	"github.com/rs/zerolog/log"

	"github.com/balazsgrill/potatodrive/win"
	"github.com/balazsgrill/potatodrive/win/cfapi"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
	"golang.org/x/sys/windows"
)

type VirtualizationInstance struct {
	rootPath    string
	shortprefix string
	longprefix  string
	fs          afero.Fs

	connectionKey cfapi.CF_CONNECTION_KEY
	lock          sync.Mutex
	watcher       *fsnotify.Watcher
}

func StartProjecting(rootPath string, filesystem afero.Fs) (win.Virtualization, error) {
	instance := &VirtualizationInstance{
		rootPath: rootPath,
		fs:       filesystem,
	}

	instance.longprefix = toLongPath(rootPath)
	instance.shortprefix = toShortPath(rootPath)

	return instance, instance.start()
}

func (instance *VirtualizationInstance) start() error {
	var registration cfapi.CF_SYNC_REGISTRATION
	registration.ProviderName = win.GetPointer("PotatoDrive")
	registration.ProviderVersion = win.GetPointer("0.1")
	registration.StructSize = uint32(unsafe.Sizeof(registration))
	var policies cfapi.CF_SYNC_POLICIES
	policies.StructSize = uint32(unsafe.Sizeof(policies))
	policies.Hydration.Primary = cfapi.CF_HYDRATION_POLICY_FULL
	policies.Hydration.Modifier = cfapi.CF_HYDRATION_POLICY_MODIFIER_AUTO_DEHYDRATION_ALLOWED
	policies.Population.Primary = cfapi.CF_POPULATION_POLICY_ALWAYS_FULL
	policies.InSync = cfapi.CF_INSYNC_POLICY_TRACK_ALL
	policies.HardLink = cfapi.CF_HARDLINK_POLICY_NONE
	policies.PlaceholderManagement = cfapi.CF_PLACEHOLDER_MANAGEMENT_POLICY_DEFAULT
	log.Print("Registering sync root")
	hr := cfapi.CfRegisterSyncRoot(win.GetPointer(instance.rootPath), &registration, &policies, cfapi.CF_REGISTER_FLAG_NONE)
	if hr != 0 {
		return win.ErrorByCode(hr)
	}

	callbacks := &cfapi.Callbacks{
		FetchData: instance.fetchData,
		//FetchPlaceholders: instance.fetchPlaceholders,
		//DeleteCompletion: instance.deleteCompletion,
	}

	log.Print("Connecting sync root")
	hr = cfapi.CfConnectSyncRoot(win.GetPointer(instance.rootPath), callbacks.CreateCallbackTable(), uintptr(unsafe.Pointer(instance)), cfapi.CF_CONNECT_FLAG_REQUIRE_FULL_FILE_PATH, &instance.connectionKey)

	err := win.ErrorByCode(hr)
	if err != nil {
		return err
	}

	instance.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	instance.watcher.Add(instance.rootPath)
	go instance.watch()

	err = instance.PerformSynchronization()
	if err != nil {
		log.Printf("Initial synchronization failed %v", err)
	}
	return nil
}

func (instance *VirtualizationInstance) readRemoteHash(remotepath string) ([]byte, error) {
	hashpath := instance.path_hashFile(remotepath)
	exists, err := afero.Exists(instance.fs, hashpath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	return afero.ReadFile(instance.fs, hashpath)
}

func getFileNameFromIdentity(info *cfapi.CF_CALLBACK_INFO) string {
	name := unsafe.Slice((*byte)(unsafe.Pointer(info.FileIdentity)), info.FileIdentityLength)
	return string(name)
}

func getPlaceholder(f fs.FileInfo) cfapi.CF_PLACEHOLDER_CREATE_INFO {
	var placeholder cfapi.CF_PLACEHOLDER_CREATE_INFO
	filename := f.Name()
	log.Print(filename)
	placeholder.RelativeFileName = win.GetPointer(filename)
	placeholder.FsMetadata.BasicInfo = toBasicInfo(f)
	identity := []byte(filename)
	placeholder.FileIdentity = uintptr(unsafe.Pointer(&identity[0]))
	log.Printf("Identity address %x", placeholder.FileIdentity)
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
		return win.ErrorByCode(hr)
	}

	hr = cfapi.CfUnregisterSyncRoot(win.GetPointer(instance.rootPath))
	return win.ErrorByCode(hr)
}

func (instance *VirtualizationInstance) PerformSynchronization() error {
	err := instance.syncRemoteToLocal()
	if err != nil {
		return err
	}
	return instance.syncLocalToRemote()
}

func (instance *VirtualizationInstance) streamLocalToRemote(filename string) error {
	file, err := os.Open(instance.path_remoteToLocal(filename))
	if err != nil {
		return err
	}
	defer file.Close()
	data := make([]byte, 1024*1024)
	targetfile, err := instance.fs.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer targetfile.Close()

	hash := md5.New()
	for {
		n, err := file.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		_, err = hash.Write(data[:n])
		if err != nil {
			return err
		}
		_, err = targetfile.Write(data[:n])
		if err != nil {
			return err
		}
	}

	return afero.WriteFile(instance.fs, instance.path_hashFile(filename), hash.Sum(nil), 0666)
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
	return &placeholderInfo, win.ErrorByCode(hr)
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
		log.Printf("deleteCompletion: remove %s failed: %v", remotepath, err)
	}
}

func (instance *VirtualizationInstance) deleteCompletion(info *cfapi.CF_CALLBACK_INFO, data *cfapi.CF_CALLBACK_PARAMETERS_DeleteCompletion) uintptr {
	instance.lock.Lock()
	defer instance.lock.Unlock()
	filename := instance.callback_getRemoteFilePath(info)
	log.Printf("deleteCompletion: %s", filename)
	//hashfilename := instance.path_hashFile(filename)

	err := instance.fs.Remove(filename)
	if err != nil {
		log.Printf("deleteCompletion: remove %s failed: %v", filename, err)
	}
	/*
		err = instance.fs.Remove(hashfilename)
		if err != nil {
			log.Printf("deleteCompletion: remove %s failed: %v", hashfilename, err)
		}*/
	return 0
}
