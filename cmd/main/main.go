package main

import (
	"github.com/rs/zerolog/log"

	"github.com/balazsgrill/potatodrive/bindings"
	"github.com/balazsgrill/potatodrive/core"
	"github.com/balazsgrill/potatodrive/ui"
)

var Version string = "0.0.0-dev"

func main() {

	mgr, err := New()
	if err != nil {
		log.Print(err)
		return
	}

	icon := ui.CreateNotifyIcon(ui.UIContext{
		Logger:  mgr.Logger,
		LogFile: mgr.logfilepath,
		Version: Version,
	})
	icon.Logger.Info().Str("version", Version).Msg("Starting PotatoDrive")
	defer icon.Close()

	err = core.CheckAlreadyRunning()
	if err != nil {
		icon.NotificationInfo("Can't start PotatoDrive", "Already running")
		return
	}

	statuslist := ui.NewStatusList()
	icon.AddAction("Show statuses", func() {
		go ui.StatusWindow(statuslist)
	})

	keys, _ := mgr.InstanceList()
	for _, keyname := range keys {
		go func(keyname string) {
			context := bindings.InstanceContext{
				Logger: icon.Logger,
				StateCallback: func(state core.ConnectionState) {
					if state.LastSyncError != nil {
						icon.Logger.Err(err).Msgf("%s is offline %v", keyname, err)
					}
				},
				FileStateCallback: statuslist.AddState,
			}
			err := mgr.StartInstance(keyname, context)
			if err != nil {
				icon.Logger.Err(err).Msgf("Failed to start %s", keyname)
			}
		}(keyname)
	}

	go bindings.CloseOnSigTerm(mgr)
	icon.Run()
}
