package potatodrive

//go:generate winrt-go-gen -class Windows.Storage.IStorageFolder
//go:generate winrt-go-gen -class Windows.Storage.Provider.StorageProviderHydrationPolicy
//go:generate winrt-go-gen -class Windows.Storage.Provider.StorageProviderHydrationPolicyModifier
//go:generate winrt-go-gen -class Windows.Storage.Provider.StorageProviderPopulationPolicy
//go:generate winrt-go-gen -class Windows.Storage.Provider.StorageProviderInSyncPolicy
//go:generate winrt-go-gen -class Windows.Storage.Provider.StorageProviderHardlinkPolicy
//go:generate winrt-go-gen -class Windows.Storage.Provider.StorageProviderProtectionMode
//go:generate winrt-go-gen -class Windows.Storage.Provider.StorageProviderSyncRootManager
//go:generate winrt-go-gen -class Windows.Storage.Provider.StorageProviderSyncRootInfo -method-filter !put_RecycleBinUri -method-filter !get_RecycleBinUri
//go:generate winrt-go-gen -class Windows.Storage.CreationCollisionOption
