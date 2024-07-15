package projfero

import (
	"encoding/binary"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"syscall"
	"unsafe"

	"C"

	"github.com/balazsgrill/projfs"
	"github.com/google/uuid"
	"github.com/spf13/afero"
)
import (
	"fmt"
	"path/filepath"
	"strings"
)

type VirtualizationInstance struct {
	rootPath        string
	fs              afero.Fs
	_instanceHandle projfs.PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT
	enumerations    map[syscall.GUID]*enumerationSession
}

type Virtualization interface {
	io.Closer
	PerformSynchronization() error
}

type enumerationSession struct {
	searchstr uintptr
	countget  int
	sentcount int
	wildcard  bool
}

func (instance *VirtualizationInstance) Close() error {
	if instance._instanceHandle == 0 {
		return errors.New("not started")
	}
	projfs.PrjStopVirtualizing(instance._instanceHandle)
	instance._instanceHandle = 0
	log.Println("Stopped virtualization")
	return nil
}

func StartProjecting(rootPath string, filesystem afero.Fs) (Virtualization, error) {
	instance := &VirtualizationInstance{
		enumerations: make(map[syscall.GUID]*enumerationSession),
	}
	return instance, instance.start(rootPath, filesystem)
}

func (instance *VirtualizationInstance) start(rootPath string, filesystem afero.Fs) error {
	if instance._instanceHandle != 0 {
		return errors.New("already started")
	}
	instance.rootPath = rootPath
	instance.fs = filesystem

	id, err := instance.ensureVirtualizationFolderExists()
	if err != nil {
		return err
	}

	hr := projfs.PrjMarkDirectoryAsPlaceholder(rootPath, "", nil, id)
	if hr != 0 {
		log.Printf("Error marking directory as placeholder: %s", projfs.ErrorByCode(hr))
		return projfs.ErrorByCode(hr)
	}
	log.Printf("Starting virtualization of '%s' (%v)", rootPath, *id)
	options := &projfs.PRJ_STARTVIRTUALIZING_OPTIONS{
		NotificationMappings: &projfs.PRJ_NOTIFICATION_MAPPING{
			NotificationBitMask: projfs.PRJ_NOTIFY_NEW_FILE_CREATED | projfs.PRJ_NOTIFY_FILE_OVERWRITTEN | projfs.PRJ_NOTIFY_FILE_HANDLE_CLOSED_FILE_DELETED | projfs.PRJ_NOTIFY_FILE_HANDLE_CLOSED_FILE_MODIFIED,
			NotificationRoot:    projfs.GetPointer(""),
		},
		NotificationMappingsCount: 1,
		PoolThreadCount:           4,
		ConcurrentThreadCount:     4,
	}
	hr = projfs.PrjStartVirtualizing(rootPath, instance.get_callbacks(), instance, options, &instance._instanceHandle)
	err = projfs.ErrorByCode(hr)
	if err != nil {
		log.Printf("Error starting virtualization: %s", err)
		return err
	}
	return instance.syncRemoteToLocal()
}

func (instance *VirtualizationInstance) PerformSynchronization() error {
	err := instance.syncRemoteToLocal()
	if err != nil {
		return err
	}
	return instance.syncLocalToRemote()
}

func (instance *VirtualizationInstance) syncRemoteToLocal() error {
	return afero.Walk(instance.fs, "", func(path string, remoteinfo fs.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if remoteinfo.IsDir() {
			return nil
		}
		localpath := instance.rootPath + "\\" + path
		var localstate projfs.PRJ_FILE_STATE
		hr := projfs.PrjGetOnDiskFileState(localpath, &localstate)
		if hr != 0 {
			return projfs.ErrorByCode(hr)
		}

		if (localstate | (projfs.PRJ_FILE_STATE_FULL & projfs.PRJ_FILE_STATE_HYDRATED_PLACEHOLDER)) != 0 {
			// check if remote is newer
			localinfo, _ := os.Stat(localpath)
			if localinfo.ModTime().UTC().Unix() < remoteinfo.ModTime().UTC().Unix() {
				log.Printf("Updating local file '%s'", path)
				var placeholderInfo projfs.PRJ_PLACEHOLDER_INFO
				FillInPlaceholderInfo(&placeholderInfo, remoteinfo)
				//err = projfs.ErrorByCode(projfs.PrjWritePlaceholderInfo(instance._instanceHandle, path, &placeholderInfo, uint32(unsafe.Sizeof(placeholderInfo))))
				err = instance.UpdateFileIfNeeded(path, &placeholderInfo, uint32(unsafe.Sizeof(placeholderInfo)), projfs.PRJ_UPDATE_ALLOW_DIRTY_METADATA|projfs.PRJ_UPDATE_ALLOW_DIRTY_DATA)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (instance *VirtualizationInstance) syncLocalToRemote() error {
	return filepath.Walk(instance.rootPath, func(localpath string, localinfo fs.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}

		path := strings.TrimPrefix(localpath, instance.rootPath)
		path = strings.TrimPrefix(path, "\\")
		if localinfo.IsDir() {
			return instance.fs.MkdirAll(path, 0777)
		}
		if strings.HasPrefix(path, ".") {
			return nil
		}

		var localstate projfs.PRJ_FILE_STATE
		hr := projfs.PrjGetOnDiskFileState(localpath, &localstate)
		if hr != 0 {
			return projfs.ErrorByCode(hr)
		}

		if (localstate | (projfs.PRJ_FILE_STATE_FULL & projfs.PRJ_FILE_STATE_HYDRATED_PLACEHOLDER)) != 0 {
			// check if local is newer
			remoteinfo, err := instance.fs.Stat(path)
			if os.IsNotExist(err) {
				// new local file, remote does not exist
				log.Printf("Uploading file '%s'", path)
				return instance.streamLocalToRemote(path)
			}
			if err != nil {
				return err
			}
			// info from walk return modification time of 2185 TODO: why?
			localinfo, err := os.Stat(localpath)
			if err != nil {
				return err
			}
			localmodtime := localinfo.ModTime()
			localtime := localmodtime.Unix()
			remotetime := remoteinfo.ModTime().Unix()
			if localtime > remotetime {
				log.Printf("Updating remote file '%s'", path)
				return instance.streamLocalToRemote(path)
			}
		}
		return nil
	})
}

func (instance *VirtualizationInstance) getVirtualizationInfoFileName() string {
	return instance.rootPath + "\\.virtualization"
}

func bytesToGuid(b []byte) *syscall.GUID {
	return &syscall.GUID{
		Data1: binary.LittleEndian.Uint32(b[0:4]),
		Data2: binary.LittleEndian.Uint16(b[4:6]),
		Data3: binary.LittleEndian.Uint16(b[6:8]),
		Data4: ([8]byte)(b[8:16]),
	}
}

func (instance *VirtualizationInstance) ensureVirtualizationFolderExists() (*syscall.GUID, error) {
	err := os.MkdirAll(instance.rootPath, 0777)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(instance.getVirtualizationInfoFileName()); errors.Is(err, os.ErrNotExist) {
		uuid, _ := uuid.NewRandom()
		id := bytesToGuid(uuid[:])
		err = os.WriteFile(instance.getVirtualizationInfoFileName(), uuid[:], 0666)
		if err != nil {
			return nil, err
		}
		return id, nil
	}

	bytes, err := os.ReadFile(instance.getVirtualizationInfoFileName())
	if err != nil {
		return nil, err
	}
	if len(bytes) != 16 {
		return nil, errors.New("invalid virtualization info file")
	}

	return bytesToGuid(bytes), nil
}

func (instance *VirtualizationInstance) get_callbacks() *projfs.PRJ_CALLBACKS {
	return &projfs.PRJ_CALLBACKS{
		NotificationCallback:              instance.Notify,
		QueryFileNameCallback:             instance.QueryFileName,
		CancelCommandCallback:             instance.CancelCommand,
		StartDirectoryEnumerationCallback: instance.StartDirectoryEnumeration,
		GetDirectoryEnumerationCallback:   instance.GetDirectoryEnumeration,
		EndDirectoryEnumerationCallback:   instance.EndDirectoryEnumeration,
		GetPlaceholderInfoCallback:        instance.GetPlaceholderInfo,
		GetFileDataCallback:               instance.GetFileData,
	}
}

func (instance *VirtualizationInstance) UpdateFileIfNeeded(relativePath string, placeholderInfo *projfs.PRJ_PLACEHOLDER_INFO, length uint32, updateFlags projfs.PRJ_UPDATE_TYPES) error {
	var failureReason projfs.PRJ_UPDATE_FAILURE_CAUSES
	err := projfs.ErrorByCode(projfs.PrjUpdateFileIfNeeded(instance._instanceHandle, relativePath, placeholderInfo, length, updateFlags, &failureReason))
	if err != nil {
		err = fmt.Errorf("UpdateFileIfNeeded failed: %w (reason: %d)", err, failureReason)
	}
	return err
}

func returncode(err error) uintptr {
	if err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func (instance *VirtualizationInstance) Notify(callbackData *projfs.PRJ_CALLBACK_DATA, IsDirectory bool, notification projfs.PRJ_NOTIFICATION, destinationFileName uintptr, operationParameters *projfs.PRJ_NOTIFICATION_PARAMETERS) uintptr {
	// operation is done on file system
	filename := callbackData.GetFilePathName()
	log.Printf("Notify: %t %d %d '%s', %d", IsDirectory, callbackData.CommandId, notification, filename, *operationParameters)
	switch notification {

	case projfs.PRJ_NOTIFICATION_NEW_FILE_CREATED:
		if IsDirectory {
			return returncode(instance.fs.Mkdir(filename, 0777))
		} else {
			_, err := instance.fs.Create(filename)
			if err != nil {
				return returncode(err)
			}
			return returncode(err)
		}
	case projfs.PRJ_NOTIFICATION_FILE_HANDLE_CLOSED_FILE_MODIFIED, projfs.PRJ_NOTIFICATION_FILE_OVERWRITTEN:
		if !IsDirectory {
			return returncode(instance.streamLocalToRemote(filename))
		}
	case projfs.PRJ_NOTIFICATION_FILE_HANDLE_CLOSED_FILE_DELETED:
		return returncode(instance.fs.Remove(filename))
	}
	return 0
}

func (instance *VirtualizationInstance) streamLocalToRemote(filename string) error {
	data, err := os.ReadFile(instance.rootPath + "\\" + filename)
	if err != nil {
		return err
	}
	return afero.WriteFile(instance.fs, filename, data, 0666)
}

func (instance *VirtualizationInstance) QueryFileName(callbackData *projfs.PRJ_CALLBACK_DATA) uintptr {
	log.Printf("QueryFileName: '%s'", callbackData.GetFilePathName())
	return 0
}

func (instance *VirtualizationInstance) CancelCommand(callbackData *projfs.PRJ_CALLBACK_DATA) uintptr {
	return 0
}

func (instance *VirtualizationInstance) StartDirectoryEnumeration(callbackData *projfs.PRJ_CALLBACK_DATA, enumerationId *syscall.GUID) uintptr {
	log.Printf("StartDirectoryEnumeration: '%v'", *enumerationId)
	instance.enumerations[*enumerationId] = &enumerationSession{
		searchstr: 0,
		countget:  0,
		sentcount: 0,
		wildcard:  false,
	}
	return 0
}

func (instance *VirtualizationInstance) EndDirectoryEnumeration(callbackData *projfs.PRJ_CALLBACK_DATA, enumerationId *syscall.GUID) uintptr {
	log.Printf("EndDirectoryEnumeration: '%v'", *enumerationId)
	instance.enumerations[*enumerationId] = nil
	return 0
}

func (instance *VirtualizationInstance) GetDirectoryEnumeration(callbackData *projfs.PRJ_CALLBACK_DATA, enumerationId *syscall.GUID, searchExpression uintptr, dirEntryBufferHandle projfs.PRJ_DIR_ENTRY_BUFFER_HANDLE) uintptr {
	filepath := callbackData.GetFilePathName()
	first := instance.enumerations[*enumerationId].countget == 0
	restart := callbackData.Flags&projfs.PRJ_CB_DATA_FLAG_ENUM_RESTART_SCAN != 0

	session, ok := instance.enumerations[*enumerationId]
	if !ok {
		return uintptr(syscall.EINVAL)
	}
	log.Printf("GetDirectoryEnumeration (%t, %t, %d) %s", first, restart, session.sentcount, filepath)

	if restart || first {
		session.sentcount = 0
		if searchExpression != 0 {
			session.searchstr = searchExpression
			session.wildcard = projfs.PrjDoesNameContainWildCards(searchExpression)
		} else {
			session.searchstr = 0
			session.wildcard = false
		}
	}
	instance.enumerations[*enumerationId].countget++

	files, err := afero.ReadDir(instance.fs, filepath)
	if err != nil {
		log.Printf("Error reading directory %s: %s", filepath, err)
		return uintptr(syscall.EIO)
	}

	for _, file := range files[session.sentcount:] {
		if session.searchstr != 0 {
			match := projfs.PrjFileNameMatch(file.Name(), session.searchstr)
			if !match {
				continue
			}
		}
		dirEntry := toBasicInfo(file)
		projfs.PrjFillDirEntryBuffer(file.Name(), &dirEntry, dirEntryBufferHandle)
		session.sentcount += 1
	}
	log.Printf("Sent %d entries", session.sentcount)
	return 0
}

func toBasicInfo(file fs.FileInfo) projfs.PRJ_FILE_BASIC_INFO {
	ftime := syscall.NsecToFiletime(file.ModTime().UnixNano())
	return projfs.PRJ_FILE_BASIC_INFO{
		IsDirectory:    file.IsDir(),
		FileSize:       file.Size(),
		CreationTime:   ftime,
		LastAccessTime: ftime,
		LastWriteTime:  ftime,
		ChangeTime:     ftime,
		FileAttributes: 0,
	}
}

func getVersionInfo(basicInfo *projfs.PRJ_FILE_BASIC_INFO) projfs.PRJ_PLACEHOLDER_VERSION_INFO {
	result := projfs.PRJ_PLACEHOLDER_VERSION_INFO{
		ProviderID: [projfs.PRJ_PLACEHOLDER_ID_LENGTH]byte{0, 0x1},
		ContentID:  [projfs.PRJ_PLACEHOLDER_ID_LENGTH]byte{0},
	}

	version := uint64(basicInfo.LastWriteTime.Nanoseconds())
	binary.LittleEndian.PutUint64(result.ContentID[:], version)
	log.Printf("Version: %d %v", version, result.ContentID)
	return result
}

func FillInPlaceholderInfo(data *projfs.PRJ_PLACEHOLDER_INFO, fileinfo fs.FileInfo) {
	data.FileBasicInfo = toBasicInfo(fileinfo)
	data.VersionInfo = getVersionInfo(&data.FileBasicInfo)
}

func (instance *VirtualizationInstance) GetPlaceholderInfo(callbackData *projfs.PRJ_CALLBACK_DATA) uintptr {
	var data projfs.PRJ_PLACEHOLDER_INFO
	filename := callbackData.GetFilePathName()
	log.Printf("GetPlaceholderInfo %s", filename)
	stat, err := instance.fs.Stat(filename)
	if os.IsNotExist(err) {
		return uintptr(0x80070002)
	}
	if err != nil {
		log.Printf("Error getting placeholder info for %s: %s", filename, err)
		return uintptr(syscall.EIO)
	}
	FillInPlaceholderInfo(&data, stat)
	return projfs.PrjWritePlaceholderInfo(instance._instanceHandle, callbackData.GetFilePathName(), &data, uint32(unsafe.Sizeof(data)))
}

func (instance *VirtualizationInstance) GetFileData(callbackData *projfs.PRJ_CALLBACK_DATA, byteOffset uint64, length uint32) uintptr {
	filename := callbackData.GetFilePathName()
	log.Printf("GetFileData %s", filename)
	file, err := instance.fs.Open(filename)
	if err != nil {
		log.Printf("Error opening file %s: %s", filename, err)
		return uintptr(syscall.EIO)
	}
	defer file.Close()
	buffer := make([]byte, length)
	_, err = file.ReadAt(buffer, int64(byteOffset))
	if err != nil {
		log.Printf("Error reading file %s: %s", filename, err)
		return uintptr(syscall.EIO)
	}
	return projfs.PrjWriteFileData(instance._instanceHandle, &callbackData.DataStreamId, &buffer[0], byteOffset, length)
}
