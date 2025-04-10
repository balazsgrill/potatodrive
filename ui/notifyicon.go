package ui

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/balazsgrill/potatodrive/assets"
	"github.com/balazsgrill/potatodrive/core"
	"github.com/lxn/walk"
	"github.com/rs/zerolog"
	"golang.org/x/sys/windows"
)

type NotifyIcon struct {
	*UIContext

	*walk.MainWindow
	zerolog.Logger
	ni   *walk.NotifyIcon
	icon *walk.Icon
}

func (ui *NotifyIcon) NotificationInfo(title string, msg string) {
	if err := ui.ni.ShowCustom(title, msg, ui.icon); err != nil {
		ui.Logger.Fatal().Err(err).Send()
	}
}

func (ui *NotifyIcon) AddAction(title string, action func()) {
	anAction := walk.NewAction()
	if err := anAction.SetText(title); err != nil {
		ui.Logger.Fatal().Err(err).Send()
	}
	anAction.Triggered().Attach(action)
	if err := ui.ni.ContextMenu().Actions().Add(anAction); err != nil {
		ui.Logger.Fatal().Err(err).Send()
	}
}

func (ui *NotifyIcon) AddActions() {
	ui.AddAction("&Open Log", func() {
		core.OpenFile(windows.Handle(ui.MainWindow.Handle()), ui.LogFile)
	})

	ui.AddAction("E&xit", func() {
		walk.App().Exit(0)
	})

	ui.AddAction("&About", func() {
		_, err := aboutDialog(ui.MainWindow, ui.UIContext)
		if err != nil {
			ui.Logger.Error().Err(err).Send()
		}
	})

}

func (ui *NotifyIcon) Close() {
	ui.ni.Dispose()
}

func CreateNotifyIcon(context *UIContext) *NotifyIcon {
	ui := &NotifyIcon{
		UIContext: context,
	}
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

	ui.AddActions()

	// The notify icon is hidden initially, so we have to make it visible.
	if err := ui.ni.SetVisible(true); err != nil {
		ui.Logger.Fatal().Err(err).Send()
	}

	ui.ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		_, err := aboutDialog(ui.MainWindow, ui.UIContext)
		if err != nil {
			logger.Error().Err(err).Send()
		}
	})

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
