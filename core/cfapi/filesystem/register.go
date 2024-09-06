package filesystem

import (
	"runtime"
	"syscall"
	"unsafe"

	"github.com/balazsgrill/potatodrive/core"
	"github.com/balazsgrill/potatodrive/core/cfapi"
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

func RegisterRootPathSimple(id syscall.GUID, rootPath string) error {
	var registration cfapi.CF_SYNC_REGISTRATION
	registration.ProviderId = id
	registration.ProviderName = core.GetPointer("PotatoDrive")
	registration.ProviderVersion = core.GetPointer("0.1")
	registration.StructSize = uint32(unsafe.Sizeof(registration))
	var policies cfapi.CF_SYNC_POLICIES
	policies.StructSize = uint32(unsafe.Sizeof(policies))
	policies.Hydration.Primary = cfapi.CF_HYDRATION_POLICY_FULL
	policies.Hydration.Modifier = cfapi.CF_HYDRATION_POLICY_MODIFIER_AUTO_DEHYDRATION_ALLOWED
	policies.Population.Primary = cfapi.CF_POPULATION_POLICY_ALWAYS_FULL
	policies.InSync = cfapi.CF_INSYNC_POLICY_TRACK_ALL
	policies.HardLink = cfapi.CF_HARDLINK_POLICY_NONE
	policies.PlaceholderManagement = cfapi.CF_PLACEHOLDER_MANAGEMENT_POLICY_DEFAULT
	hr := cfapi.CfRegisterSyncRoot(core.GetPointer(rootPath), &registration, &policies, cfapi.CF_REGISTER_FLAG_NONE)
	if hr != 0 {
		return core.ErrorByCode(hr)
	}
	return nil
}

func UnregisterRootPathSimple(rootPath string) error {
	hr := cfapi.CfUnregisterSyncRoot(core.GetPointer(rootPath))
	return core.ErrorByCode(hr)
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
		// TODO should not register again, but it fails if I do so
		// unregister first
		err = provider.StorageProviderSyncRootManagerUnregister(eid)
		if err != nil {
			return err
		}

	}

	info.SetAllowPinning(true)
	info.SetId(id)
	info.SetPath(folder)
	info.SetDisplayNameResource("PotatoDrive " + rootPath)
	info.SetIconResource(core.InstalledFile(core.POTATOICO, true))
	info.SetVersion("1")
	info.SetHydrationPolicy(provider.StorageProviderHydrationPolicyFull)
	info.SetHydrationPolicyModifier(provider.StorageProviderHydrationPolicyModifierAutoDehydrationAllowed)
	info.SetPopulationPolicy(provider.StorageProviderPopulationPolicyAlwaysFull)
	info.SetInSyncPolicy(provider.StorageProviderInSyncPolicyDirectoryLastWriteTime | provider.StorageProviderInSyncPolicyFileLastWriteTime)
	info.SetHardlinkPolicy(provider.StorageProviderHardlinkPolicyNone)
	return provider.StorageProviderSyncRootManagerRegister(info)
}
