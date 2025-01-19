package main

import "github.com/saltosystems/winrt-go/windows/storage/provider"

func unreg(id string) {
	err := provider.StorageProviderSyncRootManagerUnregister(id)
	if err != nil {
		panic(err)
	}
}
