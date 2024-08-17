package filesystem

import (
	"runtime"
	"unsafe"

	"github.com/saltosystems/winrt-go/windows/storage"
	"github.com/saltosystems/winrt-go/windows/storage/provider"
)

func getFolder(folder string) (*storage.IStorageFolder, error) {
	op, err := storage.StorageFolderGetFolderFromPathAsync(folder)
	if err != nil {
		return nil, err
	}
	ptr, err := op.GetResults()
	return (*storage.IStorageFolder)(unsafe.Pointer(ptr)), err
}

func RegisterRootPath(id string, rootPath string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var info provider.StorageProviderSyncRootInfo
	folder, err := getFolder(rootPath)
	if err != nil {
		return err
	}

	existing, err := provider.StorageProviderSyncRootManagerGetSyncRootInformationForFolder(folder)
	if err == nil && existing != nil {
		// Already registered
		eid, _ := existing.GetId()
		if eid == id {
			// No need to register again
			return nil
		} else {
			// unregister first
			err = provider.StorageProviderSyncRootManagerUnregister(eid)
			if err != nil {
				return err
			}
		}
	}

	info.SetAllowPinning(true)
	info.SetId(id)
	info.SetPath(folder)
	info.SetDisplayNameResource("PotatoDrive " + rootPath)
	info.SetIconResource("potato.ico")
	info.SetVersion("1")
	info.SetHydrationPolicy(provider.StorageProviderHydrationPolicyFull)
	info.SetHydrationPolicyModifier(provider.StorageProviderHydrationPolicyModifierAutoDehydrationAllowed)
	info.SetPopulationPolicy(provider.StorageProviderPopulationPolicyAlwaysFull)
	info.SetInSyncPolicy(provider.StorageProviderInSyncPolicyDirectoryLastWriteTime | provider.StorageProviderInSyncPolicyFileLastWriteTime)
	info.SetHardlinkPolicy(provider.StorageProviderHardlinkPolicyNone)
	return provider.StorageProviderSyncRootManagerRegister(&info)
}
