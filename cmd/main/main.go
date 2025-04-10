package main

import (
	"github.com/rs/zerolog/log"

	"github.com/balazsgrill/potatodrive/bindings"
	"github.com/balazsgrill/potatodrive/core"
	"github.com/balazsgrill/potatodrive/ui"
)

var Version string = "0.0.0-dev"

func main() {

	uicontext := ui.NewUIContext(Version)

	mgr, err := New(uicontext)
	if err != nil {
		log.Print(err)
		return
	}

	statuslist := ui.NewTaskListModel()
	ui.CreateTaskListWindow(uicontext, statuslist)

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
						icon.Logger.Err(state.LastSyncError).Msgf("%s is offline %v", keyname, state.LastSyncError)
					}
				},
				FileStateCallback: core.AsCallbacks(statuslist.TaskStateListener),
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
