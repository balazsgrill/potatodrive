package main

import (
	"flag"
	"log"

	"github.com/balazsgrill/potatodrive/bindings"
	cs3 "github.com/balazsgrill/potatodrive/bindings/s3"
	"golang.org/x/sys/windows/registry"
)

func main() {
	regkey := flag.String("regkey", "", "Registry key that holds configuration. If set, all other arguments are ignored")
	config := &cs3.Config{}
	bindings.ConfigToFlags(config)
	flag.Parse()

	if *regkey != "" {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, *regkey, registry.QUERY_VALUE)
		if err != nil {
			log.Panic(err)
		}
		bindings.ReadConfigFromRegistry(key, config)
	}

	err := config.Validate()
	if err != nil {
		panic(err)
	}

	fs, err := config.ToFileSystem()
	if err != nil {
		log.Panic(err)
	}

	closer, err := bindings.BindVirtualizationInstance(config.LocalPath, fs)
	if err != nil {
		log.Panic(err)
	}
	bindings.CloseOnSigTerm(closer)
}
