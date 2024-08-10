//go:build windows

package cfapi

import (
	"syscall"
	"unsafe"
)

var (
	cldapilib                             = syscall.NewLazyDLL("cldapi.dll")
	cfCloseHandle                         = cldapilib.NewProc("CfCloseHandle")
	cfConnectSyncRoot                     = cldapilib.NewProc("CfConnectSyncRoot")
	cfConvertToPlaceholder                = cldapilib.NewProc("CfConvertToPlaceholder")
	cfCreatePlaceholders                  = cldapilib.NewProc("CfCreatePlaceholders")
	cfDisconnectSyncRoot                  = cldapilib.NewProc("CfDisconnectSyncRoot")
	cfExecute                             = cldapilib.NewProc("CfExecute")
	cfGetCorrelationVector                = cldapilib.NewProc("CfGetCorrelationVector")
	cfGetPlaceholderInfo                  = cldapilib.NewProc("CfGetPlaceholderInfo")
	cfGetPlaceholderRangeInfo             = cldapilib.NewProc("CfGetPlaceholderRangeInfo")
	cfGetPlaceholderRangeInfoForHydration = cldapilib.NewProc("CfGetPlaceholderRangeInfoForHydration")
	cfGetPlaceholderStateFromAttributeTag = cldapilib.NewProc("CfGetPlaceholderStateFromAttributeTag")
	cfGetPlaceholderStateFromFileInfo     = cldapilib.NewProc("CfGetPlaceholderStateFromFileInfo")
	cfGetPlaceholderStateFromFindData     = cldapilib.NewProc("CfGetPlaceholderStateFromFindData")
	cfGetPlatformInfo                     = cldapilib.NewProc("CfGetPlatformInfo")
	cfGetSyncRootInfoByHandle             = cldapilib.NewProc("CfGetSyncRootInfoByHandle")
	cfGetSyncRootInfoByPath               = cldapilib.NewProc("CfGetSyncRootInfoByPath")
	cfGetTransferKey                      = cldapilib.NewProc("CfGetTransferKey")
	cfGetWin32HandleFromProtectedHandle   = cldapilib.NewProc("CfGetWin32HandleFromProtectedHandle")
	cfHydratePlaceholder                  = cldapilib.NewProc("CfHydratePlaceholder")
	cfOpenFileWithOplock                  = cldapilib.NewProc("CfOpenFileWithOplock")
	cfQuerySyncProviderStatus             = cldapilib.NewProc("CfQuerySyncProviderStatus")
	cfReferenceProtectedHandle            = cldapilib.NewProc("CfReferenceProtectedHandle")
	cfRegisterSyncRoot                    = cldapilib.NewProc("CfRegisterSyncRoot")
	cfReleaseProtectedHandle              = cldapilib.NewProc("CfReleaseProtectedHandle")
	cfReleaseTransferKey                  = cldapilib.NewProc("CfReleaseTransferKey")
	cfReportProviderProgress              = cldapilib.NewProc("CfReportProviderProgress")
	cfReportProviderProgress2             = cldapilib.NewProc("CfReportProviderProgress2")
	cfReportSyncStatus                    = cldapilib.NewProc("CfReportSyncStatus")
	cfRevertPlaceholder                   = cldapilib.NewProc("CfRevertPlaceholder")
	cfSetCorrelationVector                = cldapilib.NewProc("CfSetCorrelationVector")
	cfSetInSyncState                      = cldapilib.NewProc("CfSetInSyncState")
	cfSetPinState                         = cldapilib.NewProc("CfSetPinState")
	cfUnregisterSyncRoot                  = cldapilib.NewProc("CfUnregisterSyncRoot")
	cfUpdatePlaceholder                   = cldapilib.NewProc("CfUpdatePlaceholder")
	cfUpdateSyncProviderStatus            = cldapilib.NewProc("CfUpdateSyncProviderStatus")
)

func CfCloseHandle(Handle syscall.Handle) uintptr {
	ret, _, _ := cfCloseHandle.Call(uintptr(Handle))
	return ret
}
func CfConnectSyncRoot(SyncRootPath uintptr, CallbackTable []CF_CALLBACK_REGISTRATION, CallbackContext uintptr, ConnectFlags CF_CONNECT_FLAGS, ConnectionKey *CF_CONNECTION_KEY) uintptr {
	ret, _, _ := cfConnectSyncRoot.Call(SyncRootPath, uintptr(unsafe.Pointer(&CallbackTable[0])), CallbackContext, uintptr(ConnectFlags), uintptr(unsafe.Pointer(ConnectionKey)))
	return ret
}

func CfRegisterSyncRoot(SyncRootPath uintptr, Registration *CF_SYNC_REGISTRATION, Policies *CF_SYNC_POLICIES, RegisterFlags CF_REGISTER_FLAGS) uintptr {
	ret, _, _ := cfRegisterSyncRoot.Call(SyncRootPath, uintptr(unsafe.Pointer(Registration)), uintptr(unsafe.Pointer(Policies)), uintptr(RegisterFlags))
	return ret
}

func CfConvertToPlaceholder(FileHandle syscall.Handle, FileIdentity uintptr, FileIdentityLength uint32, ConvertFlags CF_CONVERT_FLAGS, ConvertUsn uintptr, Overlapped uintptr) uintptr {
	ret, _, _ := cfConvertToPlaceholder.Call(uintptr(FileHandle), FileIdentity, uintptr(FileIdentityLength), uintptr(ConvertFlags), ConvertUsn, Overlapped)
	return ret
}

func CfCreatePlaceholders(BaseDirectoryPath uintptr, PlaceholderArray *CF_PLACEHOLDER_CREATE_INFO, PlaceholderCount uint32, CreateFlags CF_CREATE_FLAGS, EntriesProcessed *uint32) uintptr {
	ret, _, _ := cfCreatePlaceholders.Call(BaseDirectoryPath, uintptr(unsafe.Pointer(PlaceholderArray)), uintptr(PlaceholderCount), uintptr(CreateFlags), uintptr(unsafe.Pointer(EntriesProcessed)))
	return ret
}

func CfDisconnectSyncRoot(ConnectionKey CF_CONNECTION_KEY) uintptr {
	ret, _, _ := cfDisconnectSyncRoot.Call(uintptr(ConnectionKey))
	return ret
}

func CfExecute(OpInfo *CF_OPERATION_INFO, OpParams uintptr) uintptr {
	ret, _, _ := cfExecute.Call(uintptr(unsafe.Pointer(OpInfo)), OpParams)
	return ret
}

func CfGetCorrelationVector(FileHandle syscall.Handle, CorrelationVector *CorrelationVector) uintptr {
	ret, _, _ := cfGetCorrelationVector.Call(uintptr(unsafe.Pointer(CorrelationVector)))
	return ret
}

func CfGetPlaceholderInfo(FileHandle syscall.Handle, InfoClass CF_PLACEHOLDER_INFO_CLASS, InfoBuffer uintptr, InfoBufferLength uint32, ReturnedLength *uint32) uintptr {
	ret, _, _ := cfGetPlaceholderInfo.Call(uintptr(FileHandle), uintptr(InfoClass), InfoBuffer, uintptr(InfoBufferLength), uintptr(unsafe.Pointer(ReturnedLength)))
	return ret
}

func CfGetPlaceholderRangeInfo(FileHandle syscall.Handle, InfoClass CF_PLACEHOLDER_RANGE_INFO_CLASS, StartingOffset int64, Length int64, InfoBuffer uintptr, InfoBufferLength uint32, ReturnedLength *uint32) uintptr {
	ret, _, _ := cfGetPlaceholderRangeInfo.Call(uintptr(FileHandle), uintptr(InfoClass), uintptr(StartingOffset), uintptr(Length), InfoBuffer, uintptr(InfoBufferLength), uintptr(unsafe.Pointer(ReturnedLength)))
	return ret
}

func CfGetPlaceholderRangeInfoForHydration(ConnectionKey CF_CONNECTION_KEY, TransferKey CF_TRANSFER_KEY, FileId int64, InfoClass CF_PLACEHOLDER_RANGE_INFO_CLASS, StartingOffset int64, RangeLength int64, InfoBuffer uintptr, InfoBufferSize uint32, InfoBufferWritten *uint32) uintptr {
	ret, _, _ := cfGetPlaceholderRangeInfoForHydration.Call(uintptr(ConnectionKey), uintptr(TransferKey), uintptr(FileId), uintptr(InfoClass), uintptr(StartingOffset), uintptr(RangeLength), InfoBuffer, uintptr(InfoBufferSize), uintptr(unsafe.Pointer(InfoBufferWritten)))
	return ret
}

func CfGetPlaceholderStateFromAttributeTag(FileAttributes uint32, ReparseTag uint32) uintptr {
	ret, _, _ := cfGetPlaceholderStateFromAttributeTag.Call(uintptr(FileAttributes), uintptr(ReparseTag))
	return ret
}

func CfGetPlaceholderStateFromFileInfo(InfoBuffer uintptr, InfoClass FILE_INFO_BY_HANDLE_CLASS) uintptr {
	ret, _, _ := cfGetPlaceholderStateFromFileInfo.Call(InfoBuffer, uintptr(InfoClass))
	return ret
}

func CfGetPlaceholderStateFromFindData(FindData uintptr) uintptr {
	ret, _, _ := cfGetPlaceholderStateFromFindData.Call(FindData)
	return ret
}

func CfGetPlatformInfo(PlatformVersion *CF_PLATFORM_INFO) uintptr {
	ret, _, _ := cfGetPlatformInfo.Call(uintptr(unsafe.Pointer(PlatformVersion)))
	return ret
}

func CfGetSyncRootInfoByHandle(FileHandle syscall.Handle, InfoClass CF_SYNC_ROOT_INFO_CLASS, InfoBuffer uintptr, InfoBufferLength uint32, ReturnedLength *uint32) uintptr {
	ret, _, _ := cfGetSyncRootInfoByHandle.Call(uintptr(FileHandle), uintptr(InfoClass), InfoBuffer, uintptr(InfoBufferLength), uintptr(unsafe.Pointer(ReturnedLength)))
	return ret
}

func CfGetSyncRootInfoByPath(SyncRootPath uintptr, InfoClass CF_SYNC_ROOT_INFO_CLASS, InfoBuffer uintptr, InfoBufferLength uint32, ReturnedLength *uint32) uintptr {
	ret, _, _ := cfGetSyncRootInfoByPath.Call(SyncRootPath, uintptr(InfoClass), InfoBuffer, uintptr(InfoBufferLength), uintptr(unsafe.Pointer(ReturnedLength)))
	return ret
}

func CfGetTransferKey(FileHandle syscall.Handle, TransferKey *CF_TRANSFER_KEY) uintptr {
	ret, _, _ := cfGetTransferKey.Call(uintptr(FileHandle), uintptr(unsafe.Pointer(TransferKey)))
	return ret
}
func CfGetWin32HandleFromProtectedHandle(ProtectedHandle uintptr) uintptr {
	ret, _, _ := cfGetWin32HandleFromProtectedHandle.Call(ProtectedHandle)
	return ret
}

func CfHydratePlaceholder(FileHandle syscall.Handle, StartingOffset int64, Length int64, HydrateFlags CF_HYDRATE_FLAGS, Overlapped uintptr) uintptr {
	ret, _, _ := cfHydratePlaceholder.Call(uintptr(FileHandle), uintptr(StartingOffset), uintptr(Length), uintptr(HydrateFlags))
	return ret
}

func CfOpenFileWithOplock(FilePath uintptr, Flags CF_OPEN_FILE_FLAGS, ProtectedHandle *syscall.Handle) uintptr {
	ret, _, _ := cfOpenFileWithOplock.Call(FilePath, uintptr(Flags), uintptr(unsafe.Pointer(ProtectedHandle)))
	return ret
}

func CfQuerySyncProviderStatus(ConnectionKey CF_CONNECTION_KEY, SyncProviderStatus *CF_SYNC_PROVIDER_STATUS) uintptr {
	ret, _, _ := cfQuerySyncProviderStatus.Call(uintptr(ConnectionKey), uintptr(unsafe.Pointer(SyncProviderStatus)))
	return ret
}

func CfReferenceProtectedHandle(ProtectedHandle syscall.Handle) uintptr {
	ret, _, _ := cfReferenceProtectedHandle.Call(uintptr(ProtectedHandle))
	return ret
}

func CfReleaseProtectedHandle(ProtectedHandle syscall.Handle) uintptr {
	ret, _, _ := cfReleaseProtectedHandle.Call(uintptr(ProtectedHandle))
	return ret
}

func CfReleaseTransferKey(FileHandle syscall.Handle, TransferKey CF_TRANSFER_KEY) uintptr {
	ret, _, _ := cfReleaseTransferKey.Call(uintptr(FileHandle), uintptr(TransferKey))
	return ret
}

func CfReportProviderProgress(ConnectionKey CF_CONNECTION_KEY, TransferKey CF_TRANSFER_KEY, ProviderProgressTotal int64, ProviderProgressCompleted int64) uintptr {
	ret, _, _ := cfReportProviderProgress.Call(uintptr(ConnectionKey), uintptr(TransferKey), uintptr(ProviderProgressTotal), uintptr(ProviderProgressCompleted))
	return ret
}

func CfReportProviderProgress2(ConnectionKey CF_CONNECTION_KEY, TransferKey CF_TRANSFER_KEY, RequestKey CF_REQUEST_KEY, ProviderProgressTotal int64, ProviderProgressCompleted int64, TargetSessionId uint32) uintptr {
	ret, _, _ := cfReportProviderProgress2.Call(uintptr(ConnectionKey), uintptr(TransferKey), uintptr(RequestKey), uintptr(ProviderProgressTotal), uintptr(ProviderProgressCompleted), uintptr(TargetSessionId))
	return ret
}

func CfReportSyncStatus(SyncRootPath uintptr, SyncStatus *CF_SYNC_STATUS) uintptr {
	ret, _, _ := cfReportSyncStatus.Call(SyncRootPath, uintptr(unsafe.Pointer(SyncStatus)))
	return ret
}

func CfRevertPlaceholder(FileHandle syscall.Handle, RevertFlags CF_REVERT_FLAGS, Overlapped uintptr) uintptr {
	ret, _, _ := cfRevertPlaceholder.Call(uintptr(FileHandle), uintptr(RevertFlags))
	return ret
}

func CfSetCorrelationVector(FileHandle syscall.Handle, CorrelationVector *CorrelationVector) uintptr {
	ret, _, _ := cfSetCorrelationVector.Call(uintptr(FileHandle), uintptr(unsafe.Pointer(CorrelationVector)))
	return ret
}

func CfSetInSyncState(FileHandle syscall.Handle, InSyncState CF_IN_SYNC_STATE, InSyncFlags CF_SET_IN_SYNC_FLAGS, InSyncUsn *USN) uintptr {
	ret, _, _ := cfSetInSyncState.Call(uintptr(FileHandle), uintptr(InSyncState), uintptr(InSyncFlags), uintptr(unsafe.Pointer(InSyncUsn)))
	return ret
}
func CfSetPinState(FileHandle syscall.Handle, PinState CF_PIN_STATE, PinFlags CF_SET_PIN_FLAGS, Overlapped uintptr) uintptr {
	ret, _, _ := cfSetPinState.Call(uintptr(FileHandle), uintptr(PinState), uintptr(PinFlags), Overlapped)
	return ret
}

func CfUnregisterSyncRoot(SyncRootPath uintptr) uintptr {
	ret, _, _ := cfUnregisterSyncRoot.Call(SyncRootPath)
	return ret
}

func CfUpdatePlaceholder(FileHandle syscall.Handle, FsMetadata *CF_FS_METADATA, FileIdentity uintptr, FileIdentityLength uint32, DehydrateRangeArray *CF_FILE_RANGE, DehydrateRangeCount uint32, UpdateFlags CF_UPDATE_FLAGS, UpdateUsn *USN, Overlapped uintptr) uintptr {
	ret, _, _ := cfUpdatePlaceholder.Call(uintptr(FileHandle), uintptr(unsafe.Pointer(FsMetadata)), FileIdentity, uintptr(FileIdentityLength), uintptr(unsafe.Pointer(DehydrateRangeArray)), uintptr(DehydrateRangeCount), uintptr(UpdateFlags), uintptr(unsafe.Pointer(UpdateUsn)), Overlapped)
	return ret
}

func CfUpdateSyncProviderStatus(ConnectionKey CF_CONNECTION_KEY, SyncProviderStatus *CF_SYNC_PROVIDER_STATUS) uintptr {
	ret, _, _ := cfUpdateSyncProviderStatus.Call(uintptr(ConnectionKey), uintptr(unsafe.Pointer(SyncProviderStatus)))
	return ret
}
