package main

import (
	"fmt"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/integrii/flaggy"
	"github.com/saltosystems/winrt-go/windows/storage"
	"github.com/saltosystems/winrt-go/windows/storage/provider"
)

func main() {
	err := ole.RoInitialize(1)
	if err != nil {
		panic(err)
	}
	listcmd := flaggy.NewSubcommand("list")
	flaggy.AttachSubcommand(listcmd, 1)
	flaggy.Parse()

	if listcmd.Used {
		syncrootsvector, err := provider.StorageProviderSyncRootManagerGetCurrentSyncRoots()
		if err != nil {
			panic(err)
		}
		count, err := syncrootsvector.GetSize()
		if err != nil {
			panic(err)
		}
		for i := uint32(0); i < count; i++ {
			syncrootptr, err := syncrootsvector.GetAt(i)
			if err != nil {
				panic(err)
			}
			syncroot := (*provider.StorageProviderSyncRootInfo)(syncrootptr)
			id, err := syncroot.GetId()
			if err != nil {
				panic(err)
			}
			providerid, err := syncroot.GetProviderId()
			if err != nil {
				panic(err)
			}
			path, err := syncroot.GetPath()
			if err != nil {
				panic(err)
			}
			iuk, err := path.QueryInterface(ole.NewGUID(storage.SignatureIStorageItem))
			if err != nil {
				panic(err)
			}
			defer iuk.Release()
			storageitem := (*storage.IStorageItem)(unsafe.Pointer(iuk))
			ospath, err := storageitem.GetPath()
			if err != nil {
				panic(err)
			}
			fmt.Printf("ID: %s, ProviderID: %v, Path: %s\n", id, providerid, ospath)
		}

	}
}
