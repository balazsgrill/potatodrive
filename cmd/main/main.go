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

	uicontext := &ui.UIContext{
		Logger:  mgr.Logger,
		LogFile: mgr.logfilepath,
		Version: Version,
	}
	statuslist := ui.NewStatusList()
	ui.CreateStatusWindow(uicontext, statuslist)

	icon := ui.CreateNotifyIcon(uicontext)
	icon.Logger.Info().Str("version", Version).Msg("Starting PotatoDrive")
	defer icon.Close()

	err = core.CheckAlreadyRunning()
	if err != nil {
		icon.NotificationInfo("Can't start PotatoDrive", "Already running")
		return
	}

	icon.AddAction("Show statuses", func() {
		uicontext.MainWindow.Show()
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
				FileStateCallback: func(fss core.FileSyncState) {
					go uicontext.MainWindow.Synchronize(func() {
						statuslist.AddState(fss)
					})
				},
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
