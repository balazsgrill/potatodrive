package main

import (
	"io"
	"log"

	"github.com/balazsgrill/potatodrive/bindings"
	"golang.org/x/sys/windows/registry"
)

func main() {
	parentkey, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\PotatoDrive", registry.QUERY_VALUE|registry.READ)
	if err != nil {
		panic(err)
	}

	keys, err := parentkey.ReadSubKeyNames(0)
	if err != nil {
		panic(err)
	}

	var instances []io.Closer

	for _, keyname := range keys {

		key, err := registry.OpenKey(parentkey, keyname, registry.QUERY_VALUE)
		if err != nil {
			log.Printf("Open key: %v", err)
			continue
		}

		var basec bindings.BaseConfig
		err = bindings.ReadConfigFromRegistry(key, &basec)
		if err != nil {
			log.Printf("Get base config: %v", err)
			continue
		}
		config := bindings.CreateConfigByType(basec.Type)
		bindings.ReadConfigFromRegistry(key, config)
		err = config.Validate()
		if err != nil {
			log.Printf("Validate config: %v", err)
			continue
		}
		fs, err := config.ToFileSystem()
		if err != nil {
			log.Printf("Create file system: %v", err)
			continue
		}

		log.Printf("Starting %s on %s", keyname, basec.LocalPath)
		c, err := bindings.BindVirtualizationInstance(basec.LocalPath, fs)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%s started", keyname)
		instances = append(instances, c)

	}

	bindings.CloseOnSigTerm(instances...)

	select {}
}
