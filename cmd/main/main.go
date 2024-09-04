package main

import (
	"github.com/rs/zerolog/log"

	"github.com/balazsgrill/potatodrive/bindings"
	"github.com/balazsgrill/potatodrive/win"
)

var Version string = "0.0.0-dev"

func main() {

	mgr, err := New()
	if err != nil {
		log.Print(err)
		return
	}
	mgr.Logger.Info().Str("version", Version).Msg("Starting PotatoDrive")
	ui := createUI(UIContext{
		Logger:  mgr.Logger,
		LogFile: mgr.logfilepath,
	})
	defer ui.ni.Dispose()

	err = win.CheckAlreadyRunning()
	if err != nil {
		ui.NotificationInfo("Can't start PotatoDrive", "Already running")
		return
	}

	keys, _ := mgr.InstanceList()
	for _, keyname := range keys {
		go func(keyname string) {
			err := mgr.StartInstance(keyname, ui.Logger, func(state win.ConnectionState) {
				if state.LastSyncError != nil {
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
}
