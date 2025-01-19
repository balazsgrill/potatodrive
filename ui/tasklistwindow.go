package ui

import (
	"log"
	"syscall"

	"github.com/balazsgrill/potatodrive/core"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func CreateTaskListWindow(context *UIContext, model *TaskListModel) {
	var lb *walk.ListBox

	styler := &TaskStyler{
		lb:                  &lb,
		model:               model,
		dpi2StampSize:       make(map[int]walk.Size),
		widthDPI2WsPerLine:  make(map[widthDPI]int),
		textWidthDPI2Height: make(map[textWidthDPI]int),
		stateicons:          make(map[core.FileSyncStateEnum]*walk.Icon),
	}

	//styler.loadIcons()

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
