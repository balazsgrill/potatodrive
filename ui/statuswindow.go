package ui

import (
	"log"

	"github.com/balazsgrill/potatodrive/win"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type StatusList struct {
	walk.ReflectListModelBase
	statuses     []win.FileSyncState
	pathToStatus map[string]int
}

func NewStatusList() *StatusList {
	return &StatusList{
		pathToStatus: make(map[string]int),
	}
}

func (sl *StatusList) AddState(state win.FileSyncState) {
	currentindex, exists := sl.pathToStatus[state.Path]
	if exists {
		sl.statuses[currentindex] = state
	} else {
		sl.statuses = append(sl.statuses, state)
		sl.pathToStatus[state.Path] = len(sl.statuses) - 1
	}
}

func StatusWindow() {
	var mw *walk.MainWindow
	var lb *walk.ListBox

	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "PotatoDrive status",
		MinSize:  Size{200, 200},
		Size:     Size{800, 600},
		Font:     Font{Family: "Segoe UI", PointSize: 9},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				DoubleBuffering: true,
				Layout:          VBox{},
				Children: []Widget{
					ListBox{
						AssignTo:       &lb,
						MultiSelection: true,
						Model:          model,
						ItemStyler:     styler,
					},
				},
			},
		},
	}).Create(); err != nil {
		log.Fatal(err)
	}
}
