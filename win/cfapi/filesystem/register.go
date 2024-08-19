package filesystem

import (
	"runtime"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/saltosystems/winrt-go"
	"github.com/saltosystems/winrt-go/windows/foundation"
	"github.com/saltosystems/winrt-go/windows/storage"
	"github.com/saltosystems/winrt-go/windows/storage/provider"
)

func getFolder(folder string) (*storage.IStorageFolder, error) {
	op, err := storage.StorageFolderGetFolderFromPathAsync(folder)
	if err != nil {
		return nil, err
	}
	semaphore := make(chan bool)
	iid := winrt.ParameterizedInstanceGUID(foundation.GUIDAsyncOperationCompletedHandler, storage.SignatureStorageFolder)
	handler := foundation.NewAsyncOperationCompletedHandler(ole.NewGUID(iid), func(instance *foundation.AsyncOperationCompletedHandler, asyncInfo *foundation.IAsyncOperation, asyncStatus foundation.AsyncStatus) {
		semaphore <- true
	})
	err = op.SetCompleted(handler)
	if err != nil {
		return nil, err
	}
	<-semaphore
	ptr, err := op.GetResults()
	return (*storage.IStorageFolder)(unsafe.Pointer(ptr)), err
}

func RegisterRootPath(id string, rootPath string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	err := ole.RoInitialize(1)
	if err != nil {
		return err
	}
	info, err := provider.NewStorageProviderSyncRootInfo()
	if err != nil {
		return err
	}
	folder, err := getFolder(rootPath)
	if err != nil {
		return err
	}

	existing, err := provider.StorageProviderSyncRootManagerGetSyncRootInformationForFolder(folder)
	if err == nil && existing != nil {
		// Already registered
		eid, _ := existing.GetId()
		/*if eid == id {
			// No need to register again
			return nil
		} else {*/
		// unregister first
		err = provider.StorageProviderSyncRootManagerUnregister(eid)
		if err != nil {
			return err
		}
		//}
	}

	info.SetAllowPinning(true)
	info.SetId(id)
	info.SetPath(folder)
	info.SetDisplayNameResource("PotatoDrive " + rootPath)
	info.SetIconResource("C:\\git\\potatodrive\\potato.ico")
	info.SetVersion("1")
	info.SetHydrationPolicy(provider.StorageProviderHydrationPolicyFull)
	info.SetHydrationPolicyModifier(provider.StorageProviderHydrationPolicyModifierAutoDehydrationAllowed)
	info.SetPopulationPolicy(provider.StorageProviderPopulationPolicyAlwaysFull)
	info.SetInSyncPolicy(provider.StorageProviderInSyncPolicyDirectoryLastWriteTime | provider.StorageProviderInSyncPolicyFileLastWriteTime)
	info.SetHardlinkPolicy(provider.StorageProviderHardlinkPolicyNone)
	return provider.StorageProviderSyncRootManagerRegister(info)
}
