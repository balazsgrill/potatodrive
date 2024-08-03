//go:build windows

package projfs

import "syscall"
import "C"

type PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT syscall.Handle
type PRJ_DIR_ENTRY_BUFFER_HANDLE syscall.Handle
type PRJ_COMPLETE_COMMAND_TYPE uint32
type PRJ_EXT_INFO_TYPE uint32
type PRJ_FILE_STATE uint32
type PRJ_NOTIFICATION uint32
type PRJ_NOTIFY_TYPES uint32
type PRJ_PLACEHOLDER_ID uint32
type PRJ_STARTVIRTUALIZING_FLAGS uint32
type PRJ_UPDATE_FAILURE_CAUSES uint32
type PRJ_UPDATE_TYPES uint32
type PRJ_CALLBACK_DATA_FLAGS uint32

type PRJ_VIRTUALIZATION_INSTANCE_INFO struct {
	InstanceID     syscall.GUID
	WriteAlignment uint32
}

// https://learn.microsoft.com/en-us/windows/win32/api/projectedfslib/ns-projectedfslib-prj_complete_command_extended_parameters
type PRJ_COMPLETE_COMMAND_EXTENDED_PARAMETERS struct {
	CommandType uint32
	Data        uint32
}

const (
	PRJ_COMPLETE_COMMAND_TYPE_NOTIFICATION PRJ_COMPLETE_COMMAND_TYPE = 1
	PRJ_COMPLETE_COMMAND_TYPE_ENUMERATION  PRJ_COMPLETE_COMMAND_TYPE = 2

	PRJ_EXT_INFO_TYPE_SYMLINK PRJ_EXT_INFO_TYPE = 1

	PRJ_FILE_STATE_PLACEHOLDER          PRJ_FILE_STATE = 0x00000001
	PRJ_FILE_STATE_HYDRATED_PLACEHOLDER PRJ_FILE_STATE = 0x00000002
	PRJ_FILE_STATE_DIRTY_PLACEHOLDER    PRJ_FILE_STATE = 0x00000004
	PRJ_FILE_STATE_FULL                 PRJ_FILE_STATE = 0x00000008
	PRJ_FILE_STATE_TOMBSTONE            PRJ_FILE_STATE = 0x00000010

	PRJ_NOTIFICATION_FILE_OPENED                        PRJ_NOTIFICATION = 0x00000002
	PRJ_NOTIFICATION_NEW_FILE_CREATED                   PRJ_NOTIFICATION = 0x00000004
	PRJ_NOTIFICATION_FILE_OVERWRITTEN                   PRJ_NOTIFICATION = 0x00000008
	PRJ_NOTIFICATION_PRE_DELETE                         PRJ_NOTIFICATION = 0x00000010
	PRJ_NOTIFICATION_PRE_RENAME                         PRJ_NOTIFICATION = 0x00000020
	PRJ_NOTIFICATION_PRE_SET_HARDLINK                   PRJ_NOTIFICATION = 0x00000040
	PRJ_NOTIFICATION_FILE_RENAMED                       PRJ_NOTIFICATION = 0x00000080
	PRJ_NOTIFICATION_HARDLINK_CREATED                   PRJ_NOTIFICATION = 0x00000100
	PRJ_NOTIFICATION_FILE_HANDLE_CLOSED_NO_MODIFICATION PRJ_NOTIFICATION = 0x00000200
	PRJ_NOTIFICATION_FILE_HANDLE_CLOSED_FILE_MODIFIED   PRJ_NOTIFICATION = 0x00000400
	PRJ_NOTIFICATION_FILE_HANDLE_CLOSED_FILE_DELETED    PRJ_NOTIFICATION = 0x00000800
	PRJ_NOTIFICATION_FILE_PRE_CONVERT_TO_FULL           PRJ_NOTIFICATION = 0x00001000

	PRJ_NOTIFY_NONE                               PRJ_NOTIFY_TYPES = 0x00000000
	PRJ_NOTIFY_SUPPRESS_NOTIFICATIONS             PRJ_NOTIFY_TYPES = 0x00000001
	PRJ_NOTIFY_FILE_OPENED                        PRJ_NOTIFY_TYPES = 0x00000002
	PRJ_NOTIFY_NEW_FILE_CREATED                   PRJ_NOTIFY_TYPES = 0x00000004
	PRJ_NOTIFY_FILE_OVERWRITTEN                   PRJ_NOTIFY_TYPES = 0x00000008
	PRJ_NOTIFY_PRE_DELETE                         PRJ_NOTIFY_TYPES = 0x00000010
	PRJ_NOTIFY_PRE_RENAME                         PRJ_NOTIFY_TYPES = 0x00000020
	PRJ_NOTIFY_PRE_SET_HARDLINK                   PRJ_NOTIFY_TYPES = 0x00000040
	PRJ_NOTIFY_FILE_RENAMED                       PRJ_NOTIFY_TYPES = 0x00000080
	PRJ_NOTIFY_HARDLINK_CREATED                   PRJ_NOTIFY_TYPES = 0x00000100
	PRJ_NOTIFY_FILE_HANDLE_CLOSED_NO_MODIFICATION PRJ_NOTIFY_TYPES = 0x00000200
	PRJ_NOTIFY_FILE_HANDLE_CLOSED_FILE_MODIFIED   PRJ_NOTIFY_TYPES = 0x00000400
	PRJ_NOTIFY_FILE_HANDLE_CLOSED_FILE_DELETED    PRJ_NOTIFY_TYPES = 0x00000800
	PRJ_NOTIFY_FILE_PRE_CONVERT_TO_FULL           PRJ_NOTIFY_TYPES = 0x00001000
	PRJ_NOTIFY_USE_EXISTING_MASK                  PRJ_NOTIFY_TYPES = 0xFFFFFFFF

	PRJ_PLACEHOLDER_ID_LENGTH PRJ_PLACEHOLDER_ID = 128

	PRJ_FLAG_NONE                    PRJ_STARTVIRTUALIZING_FLAGS = 0x00000000
	PRJ_FLAG_USE_NEGATIVE_PATH_CACHE PRJ_STARTVIRTUALIZING_FLAGS = 0x00000001

	PRJ_UPDATE_FAILURE_CAUSE_NONE           PRJ_UPDATE_FAILURE_CAUSES = 0x00000000
	PRJ_UPDATE_FAILURE_CAUSE_DIRTY_METADATA PRJ_UPDATE_FAILURE_CAUSES = 0x00000001
	PRJ_UPDATE_FAILURE_CAUSE_DIRTY_DATA     PRJ_UPDATE_FAILURE_CAUSES = 0x00000002
	PRJ_UPDATE_FAILURE_CAUSE_TOMBSTONE      PRJ_UPDATE_FAILURE_CAUSES = 0x00000004
	PRJ_UPDATE_FAILURE_CAUSE_READ_ONLY      PRJ_UPDATE_FAILURE_CAUSES = 0x00000008

	PRJ_UPDATE_NONE                 PRJ_UPDATE_TYPES = 0x00000000
	PRJ_UPDATE_ALLOW_DIRTY_METADATA PRJ_UPDATE_TYPES = 0x00000001
	PRJ_UPDATE_ALLOW_DIRTY_DATA     PRJ_UPDATE_TYPES = 0x00000002
	PRJ_UPDATE_ALLOW_TOMBSTONE      PRJ_UPDATE_TYPES = 0x00000004
	PRJ_UPDATE_RESERVED1            PRJ_UPDATE_TYPES = 0x00000008
	PRJ_UPDATE_RESERVED2            PRJ_UPDATE_TYPES = 0x00000010
	PRJ_UPDATE_ALLOW_READ_ONLY      PRJ_UPDATE_TYPES = 0x00000020
	PRJ_UPDATE_MAX_VAL              PRJ_UPDATE_TYPES = 0x00000040

	PRJ_CB_DATA_FLAG_ENUM_RESTART_SCAN        PRJ_CALLBACK_DATA_FLAGS = 0x00000001
	PRJ_CB_DATA_FLAG_ENUM_RETURN_SINGLE_ENTRY PRJ_CALLBACK_DATA_FLAGS = 0x00000002
)

type PRJ_CALLBACK_DATA struct {
	Size                           uint32
	Flags                          PRJ_CALLBACK_DATA_FLAGS
	NamespaceVirtualizationContext PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT
	CommandId                      int32
	FileId                         syscall.GUID
	DataStreamId                   syscall.GUID
	FilePathName                   uintptr
	VersionInfo                    *PRJ_PLACEHOLDER_VERSION_INFO
	TriggeringProcessId            *uint32
	TriggeringProcessImageFileName uintptr
	InstanceContext                uintptr
}

func (data *PRJ_CALLBACK_DATA) GetFilePathName() string {
	return GetString(data.FilePathName)
}

type PRJ_CALLBACKS struct {
	StartDirectoryEnumerationCallback PRJ_START_DIRECTORY_ENUMERATION_CB
	EndDirectoryEnumerationCallback   PRJ_END_DIRECTORY_ENUMERATION_CB
	GetDirectoryEnumerationCallback   PRJ_GET_DIRECTORY_ENUMERATION_CB
	GetPlaceholderInfoCallback        PRJ_GET_PLACEHOLDER_INFO_CB
	GetFileDataCallback               PRJ_GET_FILE_DATA_CB
	QueryFileNameCallback             PRJ_QUERY_FILE_NAME_CB
	NotificationCallback              PRJ_NOTIFICATION_CB
	CancelCommandCallback             PRJ_CANCEL_COMMAND_CB
}

func (p *PRJ_CALLBACKS) to_raw() *PRJ_CALLBACKS_raw {
	return &PRJ_CALLBACKS_raw{
		StartDirectoryEnumerationCallback: syscall.NewCallback(p.StartDirectoryEnumerationCallback),
		EndDirectoryEnumerationCallback:   syscall.NewCallback(p.EndDirectoryEnumerationCallback),
		GetDirectoryEnumerationCallback:   syscall.NewCallback(p.GetDirectoryEnumerationCallback),
		GetPlaceholderInfoCallback:        syscall.NewCallback(p.GetPlaceholderInfoCallback),
		GetFileDataCallback:               syscall.NewCallback(p.GetFileDataCallback),
		QueryFileNameCallback:             syscall.NewCallback(p.QueryFileNameCallback),
		NotificationCallback:              syscall.NewCallback(p.NotificationCallback),
		CancelCommandCallback:             syscall.NewCallback(p.CancelCommandCallback),
	}
}

type PRJ_CALLBACKS_raw struct {
	StartDirectoryEnumerationCallback uintptr
	EndDirectoryEnumerationCallback   uintptr
	GetDirectoryEnumerationCallback   uintptr
	GetPlaceholderInfoCallback        uintptr
	GetFileDataCallback               uintptr
	QueryFileNameCallback             uintptr
	NotificationCallback              uintptr
	CancelCommandCallback             uintptr
}

type PRJ_CANCEL_COMMAND_CB func(*PRJ_CALLBACK_DATA) uintptr
type PRJ_END_DIRECTORY_ENUMERATION_CB func(*PRJ_CALLBACK_DATA, *syscall.GUID) uintptr

type PRJ_EXTENDED_INFO struct {
	InfoType       PRJ_EXT_INFO_TYPE
	NextInfoOffset uint32
	TargetName     uintptr
}

type PRJ_FILE_BASIC_INFO struct {
	IsDirectory    bool
	FileSize       int64
	CreationTime   syscall.Filetime
	LastAccessTime syscall.Filetime
	LastWriteTime  syscall.Filetime
	ChangeTime     syscall.Filetime
	FileAttributes uint32
}

type PRJ_START_DIRECTORY_ENUMERATION_CB func(callbackData *PRJ_CALLBACK_DATA, enumerationId *syscall.GUID) uintptr
type PRJ_GET_DIRECTORY_ENUMERATION_CB func(callbackData *PRJ_CALLBACK_DATA, enumerationId *syscall.GUID, searchExpression uintptr, dirEntryBufferHandle PRJ_DIR_ENTRY_BUFFER_HANDLE) uintptr
type PRJ_GET_FILE_DATA_CB func(callbackData *PRJ_CALLBACK_DATA, byteOffset uint64, length uint32) uintptr
type PRJ_GET_PLACEHOLDER_INFO_CB func(callbackData *PRJ_CALLBACK_DATA) uintptr
type PRJ_NOTIFICATION_CB func(callbackData *PRJ_CALLBACK_DATA, IsDirectory bool, notification PRJ_NOTIFICATION, destinationFileName uintptr, operationParameters *PRJ_NOTIFICATION_PARAMETERS) uintptr
type PRJ_QUERY_FILE_NAME_CB func(callbackData *PRJ_CALLBACK_DATA) uintptr

type PRJ_NOTIFICATION_MAPPING struct {
	NotificationBitMask PRJ_NOTIFY_TYPES
	NotificationRoot    uintptr
}

type PRJ_NOTIFICATION_PARAMETERS PRJ_NOTIFY_TYPES
type PRJ_PLACEHOLDER_INFO struct {
	FileBasicInfo PRJ_FILE_BASIC_INFO
	EaInformation struct {
		EaBufferSize    uint32
		OffsetToFirstEa uint32
	}
	SecurityInformation struct {
		SecurityBufferSize         uint32
		OffsetToSecurityDescriptor uint32
	}
	StreamsInformationstruct struct {
		StreamsInfoBufferSize   uint32
		OffsetToFirstStreamInfo uint32
	}
	VersionInfo  PRJ_PLACEHOLDER_VERSION_INFO
	VariableData *uint8
}

type PRJ_PLACEHOLDER_VERSION_INFO struct {
	ProviderID [PRJ_PLACEHOLDER_ID_LENGTH]byte
	ContentID  [PRJ_PLACEHOLDER_ID_LENGTH]byte
}

type PRJ_STARTVIRTUALIZING_OPTIONS struct {
	Flags                     PRJ_STARTVIRTUALIZING_FLAGS
	PoolThreadCount           uint32
	ConcurrentThreadCount     uint32
	NotificationMappings      *PRJ_NOTIFICATION_MAPPING
	NotificationMappingsCount uint32
}
