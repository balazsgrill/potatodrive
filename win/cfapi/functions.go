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

func CfCreatePlaceholders(BaseDirectoryPath uintptr, PlaceholderArray []CF_PLACEHOLDER_CREATE_INFO, PlaceholderCount uint32, CreateFlags CF_CREATE_FLAGS, EntriesProcessed *uint32) uintptr {
	ret, _, _ := cfCreatePlaceholders.Call(BaseDirectoryPath, uintptr(unsafe.Pointer(&PlaceholderArray[0])), uintptr(PlaceholderCount), uintptr(CreateFlags), uintptr(unsafe.Pointer(EntriesProcessed)))
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
