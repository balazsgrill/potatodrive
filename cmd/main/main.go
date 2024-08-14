package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/balazsgrill/potatodrive/bindings"
	"golang.org/x/sys/windows/registry"
)

func initLogger() (zerolog.Logger, io.Closer) {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	logfolder := filepath.Join(cachedir, "PotatoDrive")
	err = os.MkdirAll(logfolder, 0777)
	if err != nil {
		panic(err)
	}

	logfile := "potatodrive.log"
	logf, err := os.OpenFile(filepath.Join(logfolder, logfile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	return log.Output(logf).With().Timestamp().Logger(), logf
}

func startInstance(parentkey registry.Key, keyname string) (io.Closer, error) {
	key, err := registry.OpenKey(parentkey, keyname, registry.QUERY_VALUE)
	if err != nil {
		log.Printf("Open key: %v", err)
		return nil, err
	}

	var basec bindings.BaseConfig
	err = bindings.ReadConfigFromRegistry(key, &basec)
	if err != nil {
		log.Printf("Get base config: %v", err)
		return nil, err
	}
	config := bindings.CreateConfigByType(basec.Type)
	bindings.ReadConfigFromRegistry(key, config)
	err = config.Validate()
	if err != nil {
		log.Printf("Validate config: %v", err)
		return nil, err
	}
	fs, err := config.ToFileSystem()
	if err != nil {
		log.Printf("Create file system: %v", err)
		return nil, err
	}

	log.Printf("Starting %s on %s", keyname, basec.LocalPath)
	c, err := bindings.BindVirtualizationInstance(basec.LocalPath, fs)
	if err != nil {
		log.Print(err)
	}
	log.Printf("%s started", keyname)
	return c, nil
}

func main() {
	var logf io.Closer
	log.Logger, logf = initLogger()
	defer logf.Close()

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
		c, err := startInstance(parentkey, keyname)
		if err != nil {
			log.Printf("Start instance: %v", err)
		} else {
			instances = append(instances, c)
		}
	}

	bindings.CloseOnSigTerm(instances...)

	select {}
}
