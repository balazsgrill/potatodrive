//go:build windows

package projfs

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	projectedfslib                   = syscall.NewLazyDLL("ProjectedFSLib.dll")
	prjAllocateAlignedBuffer         = projectedfslib.NewProc("PrjAllocateAlignedBuffer")
	prjClearNegativePathCache        = projectedfslib.NewProc("PrjClearNegativePathCache")
	prjCompleteCommand               = projectedfslib.NewProc("PrjCompleteCommand")
	prjDeleteFile                    = projectedfslib.NewProc("PrjDeleteFile")
	prjDoesNameContainWildCards      = projectedfslib.NewProc("PrjDoesNameContainWildCards")
	prjFileNameCompare               = projectedfslib.NewProc("PrjFileNameCompare")
	prjFileNameMatch                 = projectedfslib.NewProc("PrjFileNameMatch")
	prjFillDirEntryBuffer            = projectedfslib.NewProc("PrjFillDirEntryBuffer")
	prjFillDirEntryBuffer2           = projectedfslib.NewProc("PrjFillDirEntryBuffer2")
	prjFreeAlignedBuffer             = projectedfslib.NewProc("PrjFreeAlignedBuffer")
	prjGetOnDiskFileState            = projectedfslib.NewProc("PrjGetOnDiskFileState")
	prjGetVirtualizationInstanceInfo = projectedfslib.NewProc("PrjGetVirtualizationInstanceInfo")
	prjMarkDirectoryAsPlaceholder    = projectedfslib.NewProc("PrjMarkDirectoryAsPlaceholder")
	prjStartVirtualizing             = projectedfslib.NewProc("PrjStartVirtualizing")
	prjStopVirtualizing              = projectedfslib.NewProc("PrjStopVirtualizing")
	prjUpdateFileIfNeeded            = projectedfslib.NewProc("PrjUpdateFileIfNeeded")
	prjWriteFileData                 = projectedfslib.NewProc("PrjWriteFileData")
	prjWritePlaceholderInfo          = projectedfslib.NewProc("PrjWritePlaceholderInfo")
	prjWritePlaceholderInfo2         = projectedfslib.NewProc("PrjWritePlaceholderInfo2")
)

func ErrorByCode(result uintptr) error {
	if result == 0 {
		return nil
	} else {
		return fmt.Errorf("error result: %x", result)
	}
}

func PrjAllocateAlignedBuffer(namespaceVirtualizationContext PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT, size uint32) uintptr {
	res, _, _ := prjAllocateAlignedBuffer.Call(uintptr(namespaceVirtualizationContext), uintptr(size))
	return res
}

func PrjClearNegativePathCache(namespaceVirtualizationContext PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT, totalEntryNumber *uint32) uintptr {
	res, _, _ := prjClearNegativePathCache.Call(uintptr(namespaceVirtualizationContext), uintptr(unsafe.Pointer(totalEntryNumber)))
	return res
}

func PrjGetVirtualizationInstanceInfo(namespaceVirtualizationContext PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT, virtualizationInstanceInfo *PRJ_VIRTUALIZATION_INSTANCE_INFO) uintptr {
	res, _, _ := prjGetVirtualizationInstanceInfo.Call(uintptr(namespaceVirtualizationContext), uintptr(unsafe.Pointer(virtualizationInstanceInfo)))
	return res
}

func PrjCompleteCommand(namespaceVirtualizationContext PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT, commandId uint32, completionresult int32, extendedParameters *PRJ_COMPLETE_COMMAND_EXTENDED_PARAMETERS) uintptr {
	res, _, _ := prjCompleteCommand.Call(uintptr(namespaceVirtualizationContext), uintptr(commandId), uintptr(completionresult), uintptr(unsafe.Pointer(extendedParameters)))
	return res
}

func PrjDeleteFile(namespaceVirtualizationContext PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT, destinationFileName string, updateFlags uint32, failureReason *PRJ_UPDATE_FAILURE_CAUSES) uintptr {
	sf := GetPointer(destinationFileName)
	res, _, _ := prjDeleteFile.Call(uintptr(namespaceVirtualizationContext), sf, uintptr(updateFlags), uintptr(unsafe.Pointer(failureReason)))
	return res
}

func PrjDoesNameContainWildCards(searchExpression uintptr) bool {
	b, _, _ := prjDoesNameContainWildCards.Call(searchExpression)
	return b != 0
}

func PrjFileNameCompare(f1 string, f2 string) int32 {
	sf1 := GetPointer(f1)
	sf2 := GetPointer(f2)
	i1, _, _ := prjFileNameCompare.Call(sf1, sf2)
	return int32(i1)
}

func PrjFileNameMatch(name string, pattern uintptr) bool {
	sf1 := GetPointer(name)
	i1, _, _ := prjFileNameMatch.Call(sf1, pattern)
	return i1 != 0
}

func PrjFillDirEntryBuffer(filename string, fileBasicInfo *PRJ_FILE_BASIC_INFO, dirEntryBufferHandle PRJ_DIR_ENTRY_BUFFER_HANDLE) uintptr {
	sf1 := GetPointer(filename)
	res, _, _ := prjFillDirEntryBuffer.Call(sf1, uintptr(unsafe.Pointer(fileBasicInfo)), uintptr(dirEntryBufferHandle))
	return res
}

func PrjFillDirEntryBuffer2(dirEntryBufferHandle PRJ_DIR_ENTRY_BUFFER_HANDLE, filename string, fileBasicInfo *PRJ_FILE_BASIC_INFO, extendedInfo *PRJ_EXTENDED_INFO) uintptr {
	sf1 := GetPointer(filename)
	res, _, _ := prjFillDirEntryBuffer2.Call(uintptr(dirEntryBufferHandle), sf1, uintptr(unsafe.Pointer(fileBasicInfo)), uintptr(unsafe.Pointer(extendedInfo)))
	return res
}

func PrjFreeAlignedBuffer(buffer *any) uintptr {
	res, _, _ := prjFreeAlignedBuffer.Call(uintptr(unsafe.Pointer(buffer)))
	return res
}

func PrjGetOnDiskFileState(filename string, fileState *PRJ_FILE_STATE) uintptr {
	sf1 := GetPointer(filename)
	res, _, _ := prjGetOnDiskFileState.Call(sf1, uintptr(unsafe.Pointer(fileState)))
	return res
}

func PrjMarkDirectoryAsPlaceholder(rootPathName string, targetPathName string, versionInfo *PRJ_PLACEHOLDER_VERSION_INFO, virtualizationInstanceID *syscall.GUID) uintptr {
	sf1 := GetPointer(rootPathName)
	var sf2 uintptr
	if targetPathName != "" {
		sf2 = GetPointer(targetPathName)
	}
	res, _, _ := prjMarkDirectoryAsPlaceholder.Call(sf1, sf2, uintptr(unsafe.Pointer(versionInfo)), uintptr(unsafe.Pointer(virtualizationInstanceID)))
	return res
}

func PrjStartVirtualizing(virtualizationRootPath string, callbacks *PRJ_CALLBACKS, instanceContext any, options *PRJ_STARTVIRTUALIZING_OPTIONS, namespaceVirtualizationContext *PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT) uintptr {
	sf1 := GetPointer(virtualizationRootPath)
	res, _, _ := prjStartVirtualizing.Call(sf1, uintptr(unsafe.Pointer(callbacks.to_raw())), uintptr(unsafe.Pointer(&instanceContext)), uintptr(unsafe.Pointer(options)), uintptr(unsafe.Pointer(namespaceVirtualizationContext)))
	return res
}

func PrjStopVirtualizing(namespaceVirtualizationContext PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT) {
	prjStopVirtualizing.Call(uintptr(namespaceVirtualizationContext))
}

func PrjUpdateFileIfNeeded(namespaceVirtualizationContext PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT, destinationFileName string, placeholderInfo *PRJ_PLACEHOLDER_INFO, placeholderInfoSize uint32, updateFlags PRJ_UPDATE_TYPES, failureReason *PRJ_UPDATE_FAILURE_CAUSES) uintptr {
	sf1 := GetPointer(destinationFileName)
	res, _, _ := prjUpdateFileIfNeeded.Call(uintptr(namespaceVirtualizationContext), sf1, uintptr(unsafe.Pointer(placeholderInfo)), uintptr(placeholderInfoSize), uintptr(updateFlags), uintptr(unsafe.Pointer(failureReason)))
	return res
}

func PrjWriteFileData(namespaceVirtualizationContext PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT, dataStreamId *syscall.GUID, buffer *byte, byteoffset uint64, length uint32) uintptr {
	res, _, _ := prjWriteFileData.Call(uintptr(namespaceVirtualizationContext), uintptr(unsafe.Pointer(dataStreamId)), uintptr(unsafe.Pointer(buffer)), uintptr(byteoffset), uintptr(length))
	return res
}

func PrjWritePlaceholderInfo(namespaceVirtualizationContext PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT, destinationFileName string, placeholderInfo *PRJ_PLACEHOLDER_INFO, placeholderInfoSize uint32) uintptr {
	sf1 := GetPointer(destinationFileName)
	res, _, _ := prjWritePlaceholderInfo.Call(uintptr(namespaceVirtualizationContext), sf1, uintptr(unsafe.Pointer(placeholderInfo)), uintptr(placeholderInfoSize))
	return res
}

func PrjWritePlaceholderInfo2(namespaceVirtualizationContext PRJ_NAMESPACE_VIRTUALIZATION_CONTEXT, destinationFileName string, placeholderInfo *PRJ_PLACEHOLDER_INFO, placeholderInfoSize uint32, ExtendedInfo *PRJ_EXTENDED_INFO) uintptr {
	sf1 := GetPointer(destinationFileName)
	res, _, _ := prjWritePlaceholderInfo2.Call(uintptr(namespaceVirtualizationContext), sf1, uintptr(unsafe.Pointer(placeholderInfo)), uintptr(placeholderInfoSize), uintptr(unsafe.Pointer(ExtendedInfo)))
	return res
}
