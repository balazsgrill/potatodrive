//go:build windows

package cfapi

import "syscall"

type CF_CONNECTION_KEY int64
type CF_TRANSFER_KEY int64
type CF_REQUEST_KEY int64

// from https://github.com/microsoft/CorrelationVector-Go/blob/master/correlationvector/correlationvector.go
// CorrelationVector represents a lightweight vector for identifying and measuring causality.
type CorrelationVector struct {
	BaseVector  string
	Extension   int32
	Version     Version
	IsImmutable bool
}

// Version represents a version of the correlation vector protocol.
type Version int

type CF_PROCESS_INFO struct {
	StructSize    uint32
	ProcessId     uint32
	ImagePath     uintptr
	PackageName   uintptr
	ApplicationId uintptr
}

type CF_CALLBACK_INFO struct {
	StructSize             uint32
	ConnectionKey          CF_CONNECTION_KEY
	CallbackContext        uintptr
	VolumeGuidName         uintptr
	VolumeDosName          uintptr
	VolumeSerialNumber     uint32
	SyncRootFileId         int64
	SyncRootIdentity       uintptr
	SyncRootIdentityLength uint32
	FileId                 int64
	FileSize               int64
	FileIdentity           uintptr
	FileIdentityLength     uint32
	NormalizedPath         uintptr
	TransferKey            CF_TRANSFER_KEY
	PriorityHint           uint8
	CorrelationVector      *CorrelationVector
	ProcessInfo            *CF_PROCESS_INFO
	RequestKey             CF_REQUEST_KEY
}

type CF_CALLBACK_TYPE uint32

const (
	CF_CALLBACK_TYPE_FETCH_DATA                   CF_CALLBACK_TYPE = 0
	CF_CALLBACK_TYPE_VALIDATE_DATA                CF_CALLBACK_TYPE = 1
	CF_CALLBACK_TYPE_CANCEL_FETCH_DATA            CF_CALLBACK_TYPE = 2
	CF_CALLBACK_TYPE_FETCH_PLACEHOLDERS           CF_CALLBACK_TYPE = 3
	CF_CALLBACK_TYPE_CANCEL_FETCH_PLACEHOLDERS    CF_CALLBACK_TYPE = 4
	CF_CALLBACK_TYPE_NOTIFY_FILE_OPEN_COMPLETION  CF_CALLBACK_TYPE = 5
	CF_CALLBACK_TYPE_NOTIFY_FILE_CLOSE_COMPLETION CF_CALLBACK_TYPE = 6
	CF_CALLBACK_TYPE_NOTIFY_DEHYDRATE             CF_CALLBACK_TYPE = 7
	CF_CALLBACK_TYPE_NOTIFY_DEHYDRATE_COMPLETION  CF_CALLBACK_TYPE = 8
	CF_CALLBACK_TYPE_NOTIFY_DELETE                CF_CALLBACK_TYPE = 9
	CF_CALLBACK_TYPE_NOTIFY_DELETE_COMPLETION     CF_CALLBACK_TYPE = 10
	CF_CALLBACK_TYPE_NOTIFY_RENAME                CF_CALLBACK_TYPE = 11
	CF_CALLBACK_TYPE_NOTIFY_RENAME_COMPLETION     CF_CALLBACK_TYPE = 12
	CF_CALLBACK_TYPE_NONE                         CF_CALLBACK_TYPE = 0xffffffff
)

type CF_CALLBACK_REGISTRATION struct {
	Type     CF_CALLBACK_TYPE
	Callback uintptr
}

type CF_CALLBACK_CANCEL_FLAGS uint32

type CF_CALLBACK_PARAMETERS_Cancel struct {
	ParamSize uint32
	Flags     CF_CALLBACK_CANCEL_FLAGS
	FetchData struct {
		FileOffset int64
		Length     int64
	}
}

type CF_CALLBACK_FETCH_DATA_FLAGS uint32

const (
	CF_CALLBACK_FETCH_DATA_FLAG_NONE               CF_CALLBACK_FETCH_DATA_FLAGS = 0x00000000
	CF_CALLBACK_FETCH_DATA_FLAG_RECOVERY           CF_CALLBACK_FETCH_DATA_FLAGS = 0x00000001
	CF_CALLBACK_FETCH_DATA_FLAG_EXPLICIT_HYDRATION CF_CALLBACK_FETCH_DATA_FLAGS = 0x00000002
)

type CF_CALLBACK_DEHYDRATION_REASON uint32

const (
	CF_CALLBACK_DEHYDRATION_REASON_NONE              CF_CALLBACK_DEHYDRATION_REASON = 0
	CF_CALLBACK_DEHYDRATION_REASON_USER_MANUAL       CF_CALLBACK_DEHYDRATION_REASON = 1
	CF_CALLBACK_DEHYDRATION_REASON_SYSTEM_LOW_SPACE  CF_CALLBACK_DEHYDRATION_REASON = 2
	CF_CALLBACK_DEHYDRATION_REASON_SYSTEM_INACTIVITY CF_CALLBACK_DEHYDRATION_REASON = 3
	CF_CALLBACK_DEHYDRATION_REASON_SYSTEM_OS_UPGRADE CF_CALLBACK_DEHYDRATION_REASON = 4
)

type CF_CALLBACK_PARAMETERS_FetchData struct {
	ParamSize uint32
	Flags     CF_CALLBACK_FETCH_DATA_FLAGS
	// TODO, reordered fields?
	RequiredLength     int64
	RequiredFileOffset int64
	OptionalLength     int64
	OptionalFileOffset int64

	LastDehydrationTime   int64
	LastDehydrationReason CF_CALLBACK_DEHYDRATION_REASON
}

type CF_CALLBACK_VALIDATE_DATA_FLAGS uint32

const (
	CF_CALLBACK_VALIDATE_DATA_FLAG_NONE               CF_CALLBACK_VALIDATE_DATA_FLAGS = 0x00000000
	CF_CALLBACK_VALIDATE_DATA_FLAG_EXPLICIT_HYDRATION CF_CALLBACK_VALIDATE_DATA_FLAGS = 0x00000001
)

type CF_CALLBACK_PARAMETERS_ValidateData struct {
	Flags              CF_CALLBACK_VALIDATE_DATA_FLAGS
	RequiredFileOffset int64
	RequiredLength     int64
}

type CF_CALLBACK_FETCH_PLACEHOLDERS_FLAGS uint32

const (
	CF_CALLBACK_FETCH_PLACEHOLDERS_FLAG_NONE CF_CALLBACK_FETCH_PLACEHOLDERS_FLAGS = 0x00000000
)

type CF_CALLBACK_PARAMETERS_FetchPlaceholders struct {
	ParamSize uint32
	Flags     CF_CALLBACK_FETCH_PLACEHOLDERS_FLAGS
	Pattern   uintptr
}

type CF_CALLBACK_OPEN_COMPLETION_FLAGS uint32

const (
	CF_CALLBACK_OPEN_COMPLETION_FLAG_NONE                    CF_CALLBACK_OPEN_COMPLETION_FLAGS = 0x00000000
	CF_CALLBACK_OPEN_COMPLETION_FLAG_PLACEHOLDER_UNKNOWN     CF_CALLBACK_OPEN_COMPLETION_FLAGS = 0x00000001
	CF_CALLBACK_OPEN_COMPLETION_FLAG_PLACEHOLDER_UNSUPPORTED CF_CALLBACK_OPEN_COMPLETION_FLAGS = 0x00000002
)

type CF_CALLBACK_PARAMETERS_OpenCompletion struct {
	ParamSize uint32
	Flags     CF_CALLBACK_OPEN_COMPLETION_FLAGS
}

type CF_CALLBACK_CLOSE_COMPLETION_FLAGS uint32

const (
	CF_CALLBACK_CLOSE_COMPLETION_FLAG_NONE    CF_CALLBACK_CLOSE_COMPLETION_FLAGS = 0x00000000
	CF_CALLBACK_CLOSE_COMPLETION_FLAG_DELETED CF_CALLBACK_CLOSE_COMPLETION_FLAGS = 0x00000001
)

type CF_CALLBACK_PARAMETERS_CloseCompletion struct {
	ParamSize uint32
	Flags     CF_CALLBACK_CLOSE_COMPLETION_FLAGS
}

type CF_CALLBACK_DEHYDRATE_FLAGS uint32

const (
	CF_CALLBACK_DEHYDRATE_FLAG_NONE       CF_CALLBACK_DEHYDRATE_FLAGS = 0x00000000
	CF_CALLBACK_DEHYDRATE_FLAG_BACKGROUND CF_CALLBACK_DEHYDRATE_FLAGS = 0x00000001
)

type CF_CALLBACK_PARAMETERS_Dehydrate struct {
	ParamSize uint32
	Flags     CF_CALLBACK_DEHYDRATE_FLAGS
	Reason    CF_CALLBACK_DEHYDRATION_REASON
}

type CF_CALLBACK_DEHYDRATE_COMPLETION_FLAGS uint32

const (
	CF_CALLBACK_DEHYDRATE_COMPLETION_FLAG_NONE       CF_CALLBACK_DEHYDRATE_COMPLETION_FLAGS = 0x00000000
	CF_CALLBACK_DEHYDRATE_COMPLETION_FLAG_BACKGROUND CF_CALLBACK_DEHYDRATE_COMPLETION_FLAGS = 0x00000001
	CF_CALLBACK_DEHYDRATE_COMPLETION_FLAG_DEHYDRATED CF_CALLBACK_DEHYDRATE_COMPLETION_FLAGS = 0x00000002
)

type CF_CALLBACK_PARAMETERS_DehydrateCompletion struct {
	ParamSize uint32
	Flags     CF_CALLBACK_DEHYDRATE_COMPLETION_FLAGS
	Reason    CF_CALLBACK_DEHYDRATION_REASON
}

type CF_CALLBACK_RENAME_COMPLETION_FLAGS uint32

const (
	CF_CALLBACK_RENAME_COMPLETION_FLAG_NONE CF_CALLBACK_RENAME_COMPLETION_FLAGS = 0x00000000
)

type CF_CALLBACK_PARAMETERS_RenameCompletion struct {
	ParamSize  uint32
	Flags      CF_CALLBACK_RENAME_COMPLETION_FLAGS
	SourcePath uintptr
}

type CF_CALLBACK_RENAME_FLAGS uint32

const (
	CF_CALLBACK_RENAME_FLAG_NONE            CF_CALLBACK_RENAME_FLAGS = 0x00000000
	CF_CALLBACK_RENAME_FLAG_IS_DIRECTORY    CF_CALLBACK_RENAME_FLAGS = 0x00000001
	CF_CALLBACK_RENAME_FLAG_SOURCE_IN_SCOPE CF_CALLBACK_RENAME_FLAGS = 0x00000002
	CF_CALLBACK_RENAME_FLAG_TARGET_IN_SCOPE CF_CALLBACK_RENAME_FLAGS = 0x00000004
)

type CF_CALLBACK_PARAMETERS_Rename struct {
	ParamSize  uint32
	Flags      CF_CALLBACK_RENAME_FLAGS
	TargetPath uintptr
}

type CF_CALLBACK_DELETE_COMPLETION_FLAGS uint32

const (
	CF_CALLBACK_DELETE_COMPLETION_FLAG_NONE CF_CALLBACK_DELETE_COMPLETION_FLAGS = 0x00000000
)

type CF_CALLBACK_PARAMETERS_DeleteCompletion struct {
	ParamSize uint32
	Flags     CF_CALLBACK_DELETE_COMPLETION_FLAGS
}

type CF_CALLBACK_DELETE_FLAGS uint32

const (
	CF_CALLBACK_DELETE_FLAG_NONE         CF_CALLBACK_DELETE_FLAGS = 0x00000000
	CF_CALLBACK_DELETE_FLAG_IS_DIRECTORY CF_CALLBACK_DELETE_FLAGS = 0x00000001
	CF_CALLBACK_DELETE_FLAG_IS_UNDELETE  CF_CALLBACK_DELETE_FLAGS = 0x00000002
)

type CF_CALLBACK_PARAMETERS_Delete struct {
	ParamSize uint32
	Flags     CF_CALLBACK_DELETE_FLAGS
}

type CF_SYNC_REGISTRATION struct {
	StructSize             uint32
	ProviderName           uintptr
	ProviderVersion        uintptr
	SyncRootIdentity       uintptr
	SyncRootIdentityLength uint32
	FileIdentity           uintptr
	FileIdentityLength     uint32
	ProviderId             syscall.GUID
}

type CF_SYNC_POLICIES struct {
	StructSize            uint32
	Hydration             CF_HYDRATION_POLICY
	Population            CF_POPULATION_POLICY
	InSync                CF_INSYNC_POLICY
	HardLink              CF_HARDLINK_POLICY
	PlaceholderManagement CF_PLACEHOLDER_MANAGEMENT_POLICY
}

type CF_HYDRATION_POLICY_PRIMARY_USHORT int16

const (
	CF_HYDRATION_POLICY_PARTIAL     CF_HYDRATION_POLICY_PRIMARY_USHORT = 0x0000
	CF_HYDRATION_POLICY_PROGRESSIVE CF_HYDRATION_POLICY_PRIMARY_USHORT = 0x0001
	CF_HYDRATION_POLICY_FULL        CF_HYDRATION_POLICY_PRIMARY_USHORT = 0x0002
	CF_HYDRATION_POLICY_ALWAYS_FULL CF_HYDRATION_POLICY_PRIMARY_USHORT = 0x0003
)

type CF_HYDRATION_POLICY_MODIFIER_USHORT int16

const (
	CF_HYDRATION_POLICY_MODIFIER_NONE                         CF_HYDRATION_POLICY_MODIFIER_USHORT = 0x0000
	CF_HYDRATION_POLICY_MODIFIER_VALIDATION_REQUIRED          CF_HYDRATION_POLICY_MODIFIER_USHORT = 0x0001
	CF_HYDRATION_POLICY_MODIFIER_STREAMING_ALLOWED            CF_HYDRATION_POLICY_MODIFIER_USHORT = 0x0002
	CF_HYDRATION_POLICY_MODIFIER_AUTO_DEHYDRATION_ALLOWED     CF_HYDRATION_POLICY_MODIFIER_USHORT = 0x0004
	CF_HYDRATION_POLICY_MODIFIER_ALLOW_FULL_RESTART_HYDRATION CF_HYDRATION_POLICY_MODIFIER_USHORT = 0x0008
)

type CF_HYDRATION_POLICY struct {
	Primary  CF_HYDRATION_POLICY_PRIMARY_USHORT
	Modifier CF_HYDRATION_POLICY_MODIFIER_USHORT
}

type CF_POPULATION_POLICY_PRIMARY_USHORT int16

const (
	CF_POPULATION_POLICY_PARTIAL     CF_POPULATION_POLICY_PRIMARY_USHORT = 0x0000
	CF_POPULATION_POLICY_FULL        CF_POPULATION_POLICY_PRIMARY_USHORT = 0x0002
	CF_POPULATION_POLICY_ALWAYS_FULL CF_POPULATION_POLICY_PRIMARY_USHORT = 0x0003
)

type CF_POPULATION_POLICY_MODIFIER_USHORT int16

const (
	CF_POPULATION_POLICY_MODIFIER_NONE CF_POPULATION_POLICY_MODIFIER_USHORT = 0x0000
)

type CF_POPULATION_POLICY struct {
	Primary  CF_POPULATION_POLICY_PRIMARY_USHORT
	Modifier CF_POPULATION_POLICY_MODIFIER_USHORT
}

type CF_INSYNC_POLICY uint32

const (
	CF_INSYNC_POLICY_NONE                               CF_INSYNC_POLICY = 0x00000000
	CF_INSYNC_POLICY_TRACK_FILE_CREATION_TIME           CF_INSYNC_POLICY = 0x00000001
	CF_INSYNC_POLICY_TRACK_FILE_READONLY_ATTRIBUTE      CF_INSYNC_POLICY = 0x00000002
	CF_INSYNC_POLICY_TRACK_FILE_HIDDEN_ATTRIBUTE        CF_INSYNC_POLICY = 0x00000004
	CF_INSYNC_POLICY_TRACK_FILE_SYSTEM_ATTRIBUTE        CF_INSYNC_POLICY = 0x00000008
	CF_INSYNC_POLICY_TRACK_DIRECTORY_CREATION_TIME      CF_INSYNC_POLICY = 0x00000010
	CF_INSYNC_POLICY_TRACK_DIRECTORY_READONLY_ATTRIBUTE CF_INSYNC_POLICY = 0x00000020
	CF_INSYNC_POLICY_TRACK_DIRECTORY_HIDDEN_ATTRIBUTE   CF_INSYNC_POLICY = 0x00000040
	CF_INSYNC_POLICY_TRACK_DIRECTORY_SYSTEM_ATTRIBUTE   CF_INSYNC_POLICY = 0x00000080
	CF_INSYNC_POLICY_TRACK_FILE_LAST_WRITE_TIME         CF_INSYNC_POLICY = 0x00000100
	CF_INSYNC_POLICY_TRACK_DIRECTORY_LAST_WRITE_TIME    CF_INSYNC_POLICY = 0x00000200
	CF_INSYNC_POLICY_TRACK_FILE_ALL                     CF_INSYNC_POLICY = 0x0055550f
	CF_INSYNC_POLICY_TRACK_DIRECTORY_ALL                CF_INSYNC_POLICY = 0x00aaaaf0
	CF_INSYNC_POLICY_TRACK_ALL                          CF_INSYNC_POLICY = 0x00ffffff
	CF_INSYNC_POLICY_PRESERVE_INSYNC_FOR_SYNC_ENGINE    CF_INSYNC_POLICY = 0x80000000
)

type CF_HARDLINK_POLICY uint32

const (
	CF_HARDLINK_POLICY_NONE    CF_HARDLINK_POLICY = 0x00000000
	CF_HARDLINK_POLICY_ALLOWED CF_HARDLINK_POLICY = 0x00000001
)

type CF_PLACEHOLDER_MANAGEMENT_POLICY uint32

const (
	CF_PLACEHOLDER_MANAGEMENT_POLICY_DEFAULT                 CF_PLACEHOLDER_MANAGEMENT_POLICY = 0x00000000
	CF_PLACEHOLDER_MANAGEMENT_POLICY_CREATE_UNRESTRICTED     CF_PLACEHOLDER_MANAGEMENT_POLICY = 0x00000001
	CF_PLACEHOLDER_MANAGEMENT_POLICY_CONVERT_TO_UNRESTRICTED CF_PLACEHOLDER_MANAGEMENT_POLICY = 0x00000002
	CF_PLACEHOLDER_MANAGEMENT_POLICY_UPDATE_UNRESTRICTED     CF_PLACEHOLDER_MANAGEMENT_POLICY = 0x00000004
)

type CF_REGISTER_FLAGS uint32

const (
	CF_REGISTER_FLAG_NONE                                 CF_REGISTER_FLAGS = 0x00000000
	CF_REGISTER_FLAG_UPDATE                               CF_REGISTER_FLAGS = 0x00000001
	CF_REGISTER_FLAG_DISABLE_ON_DEMAND_POPULATION_ON_ROOT CF_REGISTER_FLAGS = 0x00000002
	CF_REGISTER_FLAG_MARK_IN_SYNC_ON_ROOT                 CF_REGISTER_FLAGS = 0x00000004
)

type CF_CONNECT_FLAGS uint32

const (
	CF_CONNECT_FLAG_NONE                          CF_CONNECT_FLAGS = 0x00000000
	CF_CONNECT_FLAG_REQUIRE_PROCESS_INFO          CF_CONNECT_FLAGS = 0x00000002
	CF_CONNECT_FLAG_REQUIRE_FULL_FILE_PATH        CF_CONNECT_FLAGS = 0x00000004
	CF_CONNECT_FLAG_BLOCK_SELF_IMPLICIT_HYDRATION CF_CONNECT_FLAGS = 0x00000008
)

type CF_CONVERT_FLAGS uint32

const (
	CF_CONVERT_FLAG_NONE                        CF_CONVERT_FLAGS = 0x00000000
	CF_CONVERT_FLAG_MARK_IN_SYNC                CF_CONVERT_FLAGS = 0x00000001
	CF_CONVERT_FLAG_DEHYDRATE                   CF_CONVERT_FLAGS = 0x00000002
	CF_CONVERT_FLAG_ENABLE_ON_DEMAND_POPULATION CF_CONVERT_FLAGS = 0x00000004
	CF_CONVERT_FLAG_ALWAYS_FULL                 CF_CONVERT_FLAGS = 0x00000008
	CF_CONVERT_FLAG_FORCE_CONVERT_TO_CLOUD_FILE CF_CONVERT_FLAGS = 0x00000010
)

type USN int64

type CF_PLACEHOLDER_CREATE_INFO struct {
	RelativeFileName   uintptr
	FsMetadata         CF_FS_METADATA
	FileIdentity       uintptr
	FileIdentityLength uint32
	Flags              CF_PLACEHOLDER_CREATE_FLAGS
	Result             uintptr
	CreateUsn          USN
}

type CF_FS_METADATA struct {
	BasicInfo FILE_BASIC_INFO
	FileSize  int64
}

// https://learn.microsoft.com/en-us/windows/win32/api/winbase/ns-winbase-file_basic_info
type FILE_BASIC_INFO struct {
	CreationTime   syscall.Filetime
	LastAccessTime syscall.Filetime
	LastWriteTime  syscall.Filetime
	ChangeTime     syscall.Filetime
	FileAttributes int32
}

type CF_PLACEHOLDER_CREATE_FLAGS uint32

const (
	CF_PLACEHOLDER_CREATE_FLAG_NONE                         CF_PLACEHOLDER_CREATE_FLAGS = 0x00000000
	CF_PLACEHOLDER_CREATE_FLAG_DISABLE_ON_DEMAND_POPULATION CF_PLACEHOLDER_CREATE_FLAGS = 0x00000001
	CF_PLACEHOLDER_CREATE_FLAG_MARK_IN_SYNC                 CF_PLACEHOLDER_CREATE_FLAGS = 0x00000002
	CF_PLACEHOLDER_CREATE_FLAG_SUPERSEDE                    CF_PLACEHOLDER_CREATE_FLAGS = 0x00000004
	CF_PLACEHOLDER_CREATE_FLAG_ALWAYS_FULL                  CF_PLACEHOLDER_CREATE_FLAGS = 0x00000008
)

type CF_CREATE_FLAGS uint32

const (
	CF_CREATE_FLAG_NONE          CF_CREATE_FLAGS = 0x00000000
	CF_CREATE_FLAG_STOP_ON_ERROR CF_CREATE_FLAGS = 0x00000001
)

type CF_OPERATION_INFO struct {
	StructSize        uint32
	Type              CF_OPERATION_TYPE
	ConnectionKey     CF_CONNECTION_KEY
	TransferKey       CF_TRANSFER_KEY
	CorrelationVector *CorrelationVector
	SyncStatus        *CF_SYNC_STATUS
	RequestKey        CF_REQUEST_KEY
}

type CF_OPERATION_TYPE uint32

const (
	CF_OPERATION_TYPE_TRANSFER_DATA         CF_OPERATION_TYPE = 0
	CF_OPERATION_TYPE_RETRIEVE_DATA         CF_OPERATION_TYPE = 1
	CF_OPERATION_TYPE_ACK_DATA              CF_OPERATION_TYPE = 2
	CF_OPERATION_TYPE_RESTART_HYDRATION     CF_OPERATION_TYPE = 3
	CF_OPERATION_TYPE_TRANSFER_PLACEHOLDERS CF_OPERATION_TYPE = 4
	CF_OPERATION_TYPE_ACK_DEHYDRATE         CF_OPERATION_TYPE = 5
	CF_OPERATION_TYPE_ACK_DELETE            CF_OPERATION_TYPE = 6
	CF_OPERATION_TYPE_ACK_RENAME            CF_OPERATION_TYPE = 7
)

type CF_SYNC_STATUS struct {
	StructSize        uint32
	Code              uint32
	DescriptionOffset uint32
	DescriptionLength uint32
	DeviceIdOffset    uint32
	DeviceIdLength    uint32
}

type NTSTATUS uint32

type CF_OPERATION_TRANSFER_DATA_FLAGS uint32

const (
	CF_OPERATION_TRANSFER_DATA_FLAG_NONE CF_OPERATION_TRANSFER_DATA_FLAGS = 0x00000000
)

type CF_OPERATION_PARAMETERS_TransferData struct {
	ParamSize        uint32
	Flags            CF_OPERATION_TRANSFER_DATA_FLAGS
	CompletionStatus NTSTATUS
	Buffer           uintptr
	Offset           int64
	Length           int64
}

type CF_OPERATION_RETRIEVE_DATA_FLAGS uint32

const (
	CF_OPERATION_RETRIEVE_DATA_FLAG_NONE CF_OPERATION_RETRIEVE_DATA_FLAGS = 0x00000000
)

type CF_OPERATION_PARAMETERS_RetrieveData struct {
	ParamSize      uint32
	Flags          CF_OPERATION_RETRIEVE_DATA_FLAGS
	Buffer         uintptr
	Offset         int64
	Length         int64
	ReturnedLength int64
}

type CF_OPERATION_ACK_DATA_FLAGS uint32

const (
	CF_OPERATION_ACK_DATA_FLAG_NONE CF_OPERATION_ACK_DATA_FLAGS = 0x00000000
)

type CF_OPERATION_PARAMETERS_AckData struct {
	ParamSize        uint32
	Flags            CF_OPERATION_ACK_DATA_FLAGS
	CompletionStatus NTSTATUS
	Offset           int64
	Length           int64
}

type CF_OPERATION_RESTART_HYDRATION_FLAGS uint32

const (
	CF_OPERATION_RESTART_HYDRATION_FLAG_NONE         CF_OPERATION_RESTART_HYDRATION_FLAGS = 0x00000000
	CF_OPERATION_RESTART_HYDRATION_FLAG_MARK_IN_SYNC CF_OPERATION_RESTART_HYDRATION_FLAGS = 0x00000001
)

type CF_OPERATION_PARAMETERS_RestartHydration struct {
	ParamSize          uint32
	Flags              CF_OPERATION_RESTART_HYDRATION_FLAGS
	FsMetadata         *CF_FS_METADATA
	FileIdentity       uintptr
	FileIdentityLength uint32
}

type CF_OPERATION_TRANSFER_PLACEHOLDERS_FLAGS uint32

const (
	CF_OPERATION_TRANSFER_PLACEHOLDERS_FLAG_NONE                         CF_OPERATION_TRANSFER_PLACEHOLDERS_FLAGS = 0x00000000
	CF_OPERATION_TRANSFER_PLACEHOLDERS_FLAG_STOP_ON_ERROR                CF_OPERATION_TRANSFER_PLACEHOLDERS_FLAGS = 0x00000001
	CF_OPERATION_TRANSFER_PLACEHOLDERS_FLAG_DISABLE_ON_DEMAND_POPULATION CF_OPERATION_TRANSFER_PLACEHOLDERS_FLAGS = 0x00000002
)

type CF_OPERATION_PARAMETERS_TransferPlaceholders struct {
	ParamSize             uint32
	Flags                 CF_OPERATION_TRANSFER_PLACEHOLDERS_FLAGS
	CompletionStatus      NTSTATUS
	PlaceholderTotalCount int64
	PlaceholderArray      *CF_PLACEHOLDER_CREATE_INFO
	PlaceholderCount      uint32
	EntriesProcessed      uint32
}

type CF_OPERATION_ACK_DEHYDRATE_FLAGS uint32

const (
	CF_OPERATION_ACK_DEHYDRATE_FLAG_NONE CF_OPERATION_ACK_DEHYDRATE_FLAGS = 0x00000000
)

type CF_OPERATION_PARAMETERS_AckDehydrate struct {
	ParamSize          uint32
	Flags              CF_OPERATION_ACK_DEHYDRATE_FLAGS
	CompletionStatus   NTSTATUS
	FileIdentity       uintptr
	FileIdentityLength uint32
}

type CF_OPERATION_ACK_DELETE_FLAGS uint32

const (
	CF_OPERATION_ACK_DELETE_FLAG_NONE CF_OPERATION_ACK_DELETE_FLAGS = 0x00000000
)

type CF_OPERATION_PARAMETERS_AckDelete struct {
	ParamSize        uint32
	Flags            CF_OPERATION_ACK_DELETE_FLAGS
	CompletionStatus NTSTATUS
}

type CF_OPERATION_ACK_RENAME_FLAGS uint32

const (
	CF_OPERATION_ACK_RENAME_FLAG_NONE CF_OPERATION_ACK_RENAME_FLAGS = 0x00000000
)

type CF_OPERATION_PARAMETERS_AckRename struct {
	ParamSize        uint32
	Flags            CF_OPERATION_ACK_RENAME_FLAGS
	CompletionStatus NTSTATUS
}

type CF_PLACEHOLDER_INFO_CLASS uint32

const (
	CF_PLACEHOLDER_INFO_BASIC    CF_PLACEHOLDER_INFO_CLASS = 0
	CF_PLACEHOLDER_INFO_STANDARD CF_PLACEHOLDER_INFO_CLASS = 1
)

type CF_PLACEHOLDER_RANGE_INFO_CLASS uint32

const (
	CF_PLACEHOLDER_RANGE_INFO_ONDISK    CF_PLACEHOLDER_RANGE_INFO_CLASS = 1
	CF_PLACEHOLDER_RANGE_INFO_VALIDATED CF_PLACEHOLDER_RANGE_INFO_CLASS = 2
	CF_PLACEHOLDER_RANGE_INFO_MODIFIED  CF_PLACEHOLDER_RANGE_INFO_CLASS = 3
)

// https://learn.microsoft.com/en-us/windows/win32/api/minwinbase/ne-minwinbase-file_info_by_handle_class
type FILE_INFO_BY_HANDLE_CLASS uint32

const (
	FileBasicInfo                  FILE_INFO_BY_HANDLE_CLASS = 0
	FileStandardInfo               FILE_INFO_BY_HANDLE_CLASS = 1
	FileNameInfo                   FILE_INFO_BY_HANDLE_CLASS = 2
	FileRenameInfo                 FILE_INFO_BY_HANDLE_CLASS = 3
	FileDispositionInfo            FILE_INFO_BY_HANDLE_CLASS = 4
	FileAllocationInfo             FILE_INFO_BY_HANDLE_CLASS = 5
	FileEndOfFileInfo              FILE_INFO_BY_HANDLE_CLASS = 6
	FileStreamInfo                 FILE_INFO_BY_HANDLE_CLASS = 7
	FileCompressionInfo            FILE_INFO_BY_HANDLE_CLASS = 8
	FileAttributeTagInfo           FILE_INFO_BY_HANDLE_CLASS = 9
	FileIdBothDirectoryInfo        FILE_INFO_BY_HANDLE_CLASS = 10
	FileIdBothDirectoryRestartInfo FILE_INFO_BY_HANDLE_CLASS = 11
	FileIoPriorityHintInfo         FILE_INFO_BY_HANDLE_CLASS = 12
	FileRemoteProtocolInfo         FILE_INFO_BY_HANDLE_CLASS = 13
	FileFullDirectoryInfo          FILE_INFO_BY_HANDLE_CLASS = 14
	FileFullDirectoryRestartInfo   FILE_INFO_BY_HANDLE_CLASS = 15
	FileStorageInfo                FILE_INFO_BY_HANDLE_CLASS = 16
	FileAlignmentInfo              FILE_INFO_BY_HANDLE_CLASS = 17
	FileIdInfo                     FILE_INFO_BY_HANDLE_CLASS = 18
	FileIdExtdDirectoryInfo        FILE_INFO_BY_HANDLE_CLASS = 19
	FileIdExtdDirectoryRestartInfo FILE_INFO_BY_HANDLE_CLASS = 20
	FileDispositionInfoEx          FILE_INFO_BY_HANDLE_CLASS = 21
	FileRenameInfoEx               FILE_INFO_BY_HANDLE_CLASS = 22
	FileCaseSensitiveInfo          FILE_INFO_BY_HANDLE_CLASS = 23
	FileNormalizedNameInfo         FILE_INFO_BY_HANDLE_CLASS = 24
	MaximumFileInfoByHandleClass   FILE_INFO_BY_HANDLE_CLASS = 25
)

type CF_PLATFORM_INFO struct {
	BuildNumber       uint32
	RevisionNumber    uint32
	IntegrationNumber uint32
}

type CF_SYNC_ROOT_INFO_CLASS uint32

const (
	CF_SYNC_ROOT_INFO_BASIC    CF_SYNC_ROOT_INFO_CLASS = 0
	CF_SYNC_ROOT_INFO_STANDARD CF_SYNC_ROOT_INFO_CLASS = 1
	CF_SYNC_ROOT_INFO_PROVIDER CF_SYNC_ROOT_INFO_CLASS = 2
)

type CF_HYDRATE_FLAGS uint32

const (
	CF_HYDRATE_FLAG_NONE CF_HYDRATE_FLAGS = 0x00000000
)

type CF_OPEN_FILE_FLAGS uint32

const (
	CF_OPEN_FILE_FLAG_NONE            CF_OPEN_FILE_FLAGS = 0x00000000
	CF_OPEN_FILE_FLAG_EXCLUSIVE       CF_OPEN_FILE_FLAGS = 0x00000001
	CF_OPEN_FILE_FLAG_WRITE_ACCESS    CF_OPEN_FILE_FLAGS = 0x00000002
	CF_OPEN_FILE_FLAG_DELETE_ACCESS   CF_OPEN_FILE_FLAGS = 0x00000004
	CF_OPEN_FILE_FLAG_FOREGROUND_SYNC CF_OPEN_FILE_FLAGS = 0x00000008
)

type CF_SYNC_PROVIDER_STATUS uint32

const (
	CF_PROVIDER_STATUS_DISCONNECTED       CF_SYNC_PROVIDER_STATUS = 0x00000000
	CF_PROVIDER_STATUS_IDLE               CF_SYNC_PROVIDER_STATUS = 0x00000001
	CF_PROVIDER_STATUS_POPULATE_NAMESPACE CF_SYNC_PROVIDER_STATUS = 0x00000002
	CF_PROVIDER_STATUS_POPULATE_METADATA  CF_SYNC_PROVIDER_STATUS = 0x00000004
	CF_PROVIDER_STATUS_POPULATE_CONTENT   CF_SYNC_PROVIDER_STATUS = 0x00000008
	CF_PROVIDER_STATUS_SYNC_INCREMENTAL   CF_SYNC_PROVIDER_STATUS = 0x00000010
	CF_PROVIDER_STATUS_SYNC_FULL          CF_SYNC_PROVIDER_STATUS = 0x00000020
	CF_PROVIDER_STATUS_CONNECTIVITY_LOST  CF_SYNC_PROVIDER_STATUS = 0x00000040
	CF_PROVIDER_STATUS_CLEAR_FLAGS        CF_SYNC_PROVIDER_STATUS = 0x00000080
	CF_PROVIDER_STATUS_TERMINATED         CF_SYNC_PROVIDER_STATUS = 0xC0000001
	CF_PROVIDER_STATUS_ERROR              CF_SYNC_PROVIDER_STATUS = 0xC0000002
)

type CF_REVERT_FLAGS uint32

const (
	CF_REVERT_FLAG_NONE CF_REVERT_FLAGS = 0x00000000
)

type CF_IN_SYNC_STATE uint32

const (
	CF_IN_SYNC_STATE_NOT_IN_SYNC CF_IN_SYNC_STATE = 0x00000000
	CF_IN_SYNC_STATE_IN_SYNC     CF_IN_SYNC_STATE = 0x00000001
)

type CF_SET_IN_SYNC_FLAGS uint32

const (
	CF_SET_IN_SYNC_FLAG_NONE CF_SET_IN_SYNC_FLAGS = 0x00000000
)

type CF_PIN_STATE uint32

const (
	CF_PIN_STATE_UNSPECIFIED CF_PIN_STATE = 0
	CF_PIN_STATE_PINNED      CF_PIN_STATE = 1
	CF_PIN_STATE_UNPINNED    CF_PIN_STATE = 2
	CF_PIN_STATE_EXCLUDED    CF_PIN_STATE = 3
	CF_PIN_STATE_INHERIT     CF_PIN_STATE = 4
)

type CF_SET_PIN_FLAGS uint32

const (
	CF_SET_PIN_FLAG_NONE                  CF_SET_PIN_FLAGS = 0x00000000
	CF_SET_PIN_FLAG_RECURSE               CF_SET_PIN_FLAGS = 0x00000001
	CF_SET_PIN_FLAG_RECURSE_ONLY          CF_SET_PIN_FLAGS = 0x00000002
	CF_SET_PIN_FLAG_RECURSE_STOP_ON_ERROR CF_SET_PIN_FLAGS = 0x00000004
)

type CF_FILE_RANGE struct {
	StartingOffset int64
	Length         int64
}

type CF_UPDATE_FLAGS uint32

const (
	CF_UPDATE_FLAG_NONE                         CF_UPDATE_FLAGS = 0x00000000
	CF_UPDATE_FLAG_VERIFY_IN_SYNC               CF_UPDATE_FLAGS = 0x00000001
	CF_UPDATE_FLAG_MARK_IN_SYNC                 CF_UPDATE_FLAGS = 0x00000002
	CF_UPDATE_FLAG_DEHYDRATE                    CF_UPDATE_FLAGS = 0x00000004
	CF_UPDATE_FLAG_ENABLE_ON_DEMAND_POPULATION  CF_UPDATE_FLAGS = 0x00000008
	CF_UPDATE_FLAG_DISABLE_ON_DEMAND_POPULATION CF_UPDATE_FLAGS = 0x00000010
	CF_UPDATE_FLAG_REMOVE_FILE_IDENTITY         CF_UPDATE_FLAGS = 0x00000020
	CF_UPDATE_FLAG_CLEAR_IN_SYNC                CF_UPDATE_FLAGS = 0x00000040
	CF_UPDATE_FLAG_REMOVE_PROPERTY              CF_UPDATE_FLAGS = 0x00000080
	CF_UPDATE_FLAG_PASSTHROUGH_FS_METADATA      CF_UPDATE_FLAGS = 0x00000100
	CF_UPDATE_FLAG_ALWAYS_FULL                  CF_UPDATE_FLAGS = 0x00000200
	CF_UPDATE_FLAG_ALLOW_PARTIAL                CF_UPDATE_FLAGS = 0x00000400
)

type CF_PLACEHOLDER_BASIC_INFO struct {
	PinState           CF_PIN_STATE
	InSyncState        CF_IN_SYNC_STATE
	FileId             int64
	SyncRootFileId     int64
	FileIdentityLength uint32
	FileIdentity       uintptr
}

type FILE_ATTRIBUTE_TAG_INFO struct {
	FileAttributes uint32
	ReparseTag     uint32
}

type CF_PLACEHOLDER_STATE uint32

const (
	CF_PLACEHOLDER_STATE_NO_STATES              CF_PLACEHOLDER_STATE = 0x00000000
	CF_PLACEHOLDER_STATE_PLACEHOLDER            CF_PLACEHOLDER_STATE = 0x00000001
	CF_PLACEHOLDER_STATE_SYNC_ROOT              CF_PLACEHOLDER_STATE = 0x00000002
	CF_PLACEHOLDER_STATE_ESSENTIAL_PROP_PRESENT CF_PLACEHOLDER_STATE = 0x00000004
	CF_PLACEHOLDER_STATE_IN_SYNC                CF_PLACEHOLDER_STATE = 0x00000008
	CF_PLACEHOLDER_STATE_PARTIAL_FILE_IN_SYNC   CF_PLACEHOLDER_STATE = 0x00000010
	CF_PLACEHOLDER_STATE_PARTIALLY_ON_DISK      CF_PLACEHOLDER_STATE = 0x00000020
	CF_PLACEHOLDER_STATE_INVALID                CF_PLACEHOLDER_STATE = 0xffffffff
)
