package main

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/balazsgrill/potatodrive/assets"
	"github.com/balazsgrill/potatodrive/win"
	"github.com/lxn/walk"
	"github.com/rs/zerolog"
	"golang.org/x/sys/windows"
)

//go:generate rsrc -manifest .\main.exe.manifest -o rsrc.syso -ico ../../potato.ico

type UI struct {
	*walk.MainWindow
	zerolog.Logger
	ni   *walk.NotifyIcon
	icon *walk.Icon
}

type UIContext struct {
	Logger  zerolog.Logger
	LogFile string
}

func (ui *UI) NotificationInfo(title string, msg string) {
	if err := ui.ni.ShowCustom(title, msg, ui.icon); err != nil {
		ui.Logger.Fatal().Err(err).Send()
	}
}

func createUI(context UIContext) *UI {
	ui := &UI{}
	logger := context.Logger

	// We need either a walk.MainWindow or a walk.Dialog for their message loop.
	// We will not make it visible in this example, though.
	var err error
	ui.MainWindow, err = walk.NewMainWindow()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	image, _, err := image.Decode(bytes.NewReader(assets.PotatoPng))
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	iconbt, err := walk.NewBitmapFromImageForDPI(image, 96)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	// We load our icon from a file.
	ui.icon, err = walk.NewIconFromBitmap(iconbt)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	// Create the notify icon and make sure we clean it up on exit.
	ui.ni, err = walk.NewNotifyIcon(ui.MainWindow)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	// Set the icon and a tool tip text.
	if err := ui.ni.SetIcon(ui.icon); err != nil {
		logger.Fatal().Err(err).Send()
	}
	if err := ui.ni.SetToolTip("Click for info or use the context menu to exit."); err != nil {
		logger.Fatal().Err(err).Send()
	}

	// When the left mouse button is pressed, bring up our balloon.
	ui.ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		_, err := aboutDialog(ui.MainWindow)
		if err != nil {
			logger.Error().Err(err).Send()
		}
	})

	openLogAction := walk.NewAction()
	if err := openLogAction.SetText("&Open Log"); err != nil {
		logger.Fatal().Err(err).Send()
	}
	openLogAction.Triggered().Attach(func() {
		win.OpenFile(windows.Handle(ui.MainWindow.Handle()), context.LogFile)
	})
	if err := ui.ni.ContextMenu().Actions().Add(openLogAction); err != nil {
		logger.Fatal().Err(err).Send()
	}

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		logger.Fatal().Err(err).Send()
	}
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ui.ni.ContextMenu().Actions().Add(exitAction); err != nil {
		logger.Fatal().Err(err).Send()
	}

	// aboutDialog
	aboudDialogAction := walk.NewAction()
	if err := aboudDialogAction.SetText("About PotatoDrive"); err != nil {
		logger.Fatal().Err(err).Send()
	}
	aboudDialogAction.Triggered().Attach(func() {
		_, err := aboutDialog(ui.MainWindow)
		if err != nil {
			logger.Error().Err(err).Send()
		}
	})
	if err := ui.ni.ContextMenu().Actions().Add(aboudDialogAction); err != nil {
		logger.Fatal().Err(err).Send()
	}

	// The notify icon is hidden initially, so we have to make it visible.
	if err := ui.ni.SetVisible(true); err != nil {
		logger.Fatal().Err(err).Send()
	}

	ui.Logger = logger.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		switch level {
		case zerolog.ErrorLevel:
			ui.ni.ShowError("Error", msg)
		case zerolog.WarnLevel, zerolog.FatalLevel:
			ui.ni.ShowWarning("Warning", msg)
		case zerolog.InfoLevel:
			ui.ni.ShowInfo("Info", msg)
		}
	}))
	return ui
}
