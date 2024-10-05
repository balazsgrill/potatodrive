package ui

import (
	"log"
	"syscall"

	"github.com/balazsgrill/potatodrive/core"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

type StatusList struct {
	walk.ReflectListModelBase
	statuses     []core.FileSyncState
	pathToStatus map[string]int
}

var _ walk.ReflectListModel = (*StatusList)(nil)

func NewStatusList() *StatusList {
	return &StatusList{
		pathToStatus: make(map[string]int),
	}
}

func (sl *StatusList) Items() interface{} {
	return sl.statuses
}

func (sl *StatusList) AddState(state core.FileSyncState) {
	currentindex, exists := sl.pathToStatus[state.Path]
	if exists {
		sl.statuses[currentindex] = state
	} else {
		sl.statuses = append(sl.statuses, state)
		sl.pathToStatus[state.Path] = len(sl.statuses) - 1
	}
	// TODO find a better way to update the list
	sl.PublishItemsReset()
}

func CreateStatusWindow(context *UIContext, model *StatusList) {
	var lb *walk.ListBox

	styler := &Styler{
		lb:                  &lb,
		model:               model,
		dpi2StampSize:       make(map[int]walk.Size),
		widthDPI2WsPerLine:  make(map[widthDPI]int),
		textWidthDPI2Height: make(map[textWidthDPI]int),
		stateicons:          make(map[core.FileSyncStateEnum]*walk.Icon),
	}

	styler.loadIcons()

	if err := (MainWindow{
		AssignTo: &context.MainWindow,
		Title:    "PotatoDrive status",
		MinSize:  Size{200, 200},
		Size:     Size{800, 600},
		Font:     Font{Family: "Segoe UI", PointSize: 9},
		Layout:   VBox{},
		Visible:  false,
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

	// https://github.com/lxn/walk/issues/326#issuecomment-461074992
	var prevWndProcPtr uintptr
	prevWndProcPtr = win.SetWindowLongPtr(context.MainWindow.Handle(), win.GWL_WNDPROC,
		syscall.NewCallback(func(hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
			if msg == win.WM_CLOSE {
				win.ShowWindow(hWnd, win.SW_HIDE)
				return 0
			}
			return win.CallWindowProc(prevWndProcPtr, hWnd, msg, wParam, lParam)
		}))
}
