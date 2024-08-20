package main

import (
	"github.com/rs/zerolog/log"

	"github.com/balazsgrill/potatodrive/bindings"
)

func main() {
	mgr, err := New()
	if err != nil {
		log.Print(err)
		return
	}
	ui := createUI(UIContext{
		Logger:  mgr.Logger,
		LogFile: mgr.logfilepath,
	})
	defer ui.ni.Dispose()

	keys, _ := mgr.InstanceList()
	for _, keyname := range keys {
		go func(keyname string) {
			err := mgr.StartInstance(keyname, ui.Logger, func(err error) {
				if err != nil {
					ui.Logger.Err(err).Msgf("%s is offline %v", keyname, err)
				}
			})
			if err != nil {
				ui.Logger.Err(err).Msgf("Failed to start %s", keyname)
			}
		}(keyname)
	}

	go bindings.CloseOnSigTerm(mgr)
	ui.Run()
	//select {}
}
