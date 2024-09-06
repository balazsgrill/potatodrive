package cfapi

import "syscall"
import "C"

type Callback_FetchData func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_FetchData) uintptr
type Callback_ValidateData func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_ValidateData) uintptr
type Callback_CancelFetchData func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_Cancel) uintptr
type Callback_FetchPlaceholders func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_FetchPlaceholders) uintptr
type Callback_CancelFetchPlaceholders func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_Cancel) uintptr
type Callback_OpenCompletion func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_OpenCompletion) uintptr
type Callback_CloseCompletion func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_CloseCompletion) uintptr
type Callback_Dehydrate func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_Dehydrate) uintptr
type Callback_DehydrateCompletion func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_DehydrateCompletion) uintptr
type Callback_Delete func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_Delete) uintptr
type Callback_DeleteCompletion func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_DeleteCompletion) uintptr
type Callback_Rename func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_Rename) uintptr
type Callback_RenameCompletion func(*CF_CALLBACK_INFO, *CF_CALLBACK_PARAMETERS_RenameCompletion) uintptr

type Callbacks struct {
	FetchData               Callback_FetchData
	ValidateData            Callback_ValidateData
	CancelFetchData         Callback_CancelFetchData
	FetchPlaceholders       Callback_FetchPlaceholders
	CancelFetchPlaceholders Callback_CancelFetchPlaceholders
	OpenCompletion          Callback_OpenCompletion
	CloseCompletion         Callback_CloseCompletion
	Dehydrate               Callback_Dehydrate
	DehydrateCompletion     Callback_DehydrateCompletion
	Delete                  Callback_Delete
	DeleteCompletion        Callback_DeleteCompletion
	Rename                  Callback_Rename
	RenameCompletion        Callback_RenameCompletion
}

func (cb *Callbacks) CreateCallbackTable() []CF_CALLBACK_REGISTRATION {
	result := make([]CF_CALLBACK_REGISTRATION, 14)
	count := 0
	if cb.FetchData != nil {
		result[count].Callback = syscall.NewCallback(cb.FetchData)
		result[count].Type = CF_CALLBACK_TYPE_FETCH_DATA
		count++
	}
	if cb.ValidateData != nil {
		result[count].Callback = syscall.NewCallback(cb.ValidateData)
		result[count].Type = CF_CALLBACK_TYPE_VALIDATE_DATA
		count++
	}
	if cb.CancelFetchData != nil {
		result[count].Callback = syscall.NewCallback(cb.CancelFetchData)
		result[count].Type = CF_CALLBACK_TYPE_CANCEL_FETCH_DATA
		count++
	}
	if cb.FetchPlaceholders != nil {
		result[count].Callback = syscall.NewCallback(cb.FetchPlaceholders)
		result[count].Type = CF_CALLBACK_TYPE_FETCH_PLACEHOLDERS
		count++
	}
	if cb.CancelFetchPlaceholders != nil {
		result[count].Callback = syscall.NewCallback(cb.CancelFetchPlaceholders)
		result[count].Type = CF_CALLBACK_TYPE_CANCEL_FETCH_PLACEHOLDERS
		count++
	}
	if cb.OpenCompletion != nil {
		result[count].Callback = syscall.NewCallback(cb.OpenCompletion)
		result[count].Type = CF_CALLBACK_TYPE_NOTIFY_FILE_OPEN_COMPLETION
		count++
	}
	if cb.CloseCompletion != nil {
		result[count].Callback = syscall.NewCallback(cb.CloseCompletion)
		result[count].Type = CF_CALLBACK_TYPE_NOTIFY_FILE_CLOSE_COMPLETION
		count++
	}
	if cb.Dehydrate != nil {
		result[count].Callback = syscall.NewCallback(cb.Dehydrate)
		result[count].Type = CF_CALLBACK_TYPE_NOTIFY_DEHYDRATE
		count++
	}
	if cb.DehydrateCompletion != nil {
		result[count].Callback = syscall.NewCallback(cb.DehydrateCompletion)
		result[count].Type = CF_CALLBACK_TYPE_NOTIFY_DEHYDRATE_COMPLETION
		count++
	}
	if cb.Delete != nil {
		result[count].Callback = syscall.NewCallback(cb.Delete)
		result[count].Type = CF_CALLBACK_TYPE_NOTIFY_DELETE
		count++
	}
	if cb.DeleteCompletion != nil {
		result[count].Callback = syscall.NewCallback(cb.DeleteCompletion)
		result[count].Type = CF_CALLBACK_TYPE_NOTIFY_DELETE_COMPLETION
		count++
	}
	if cb.Rename != nil {
		result[count].Callback = syscall.NewCallback(cb.Rename)
		result[count].Type = CF_CALLBACK_TYPE_NOTIFY_RENAME
		count++
	}
	if cb.RenameCompletion != nil {
		result[count].Callback = syscall.NewCallback(cb.RenameCompletion)
		result[count].Type = CF_CALLBACK_TYPE_NOTIFY_RENAME_COMPLETION
		count++
	}
	result[count].Type = CF_CALLBACK_TYPE_NONE
	return result
}
