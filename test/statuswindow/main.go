package main

import (
	"github.com/balazsgrill/potatodrive/core"
	"github.com/balazsgrill/potatodrive/ui"
)

func main() {
	list := ui.NewStatusList()
	list.AddState(core.FileSyncState{Path: "C:\\test1.txt", State: core.FileSyncStateDone})
	ui.StatusWindow(list)
}
