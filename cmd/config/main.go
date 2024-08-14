// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"image"

	"github.com/rs/zerolog/log"

	_ "image/png"

	"github.com/balazsgrill/potatodrive/assets"
	"github.com/lxn/walk"
)

//go:generate rsrc -manifest .\main.exe.manifest -o rsrc.syso

func main() {
	// We need either a walk.MainWindow or a walk.Dialog for their message loop.
	// We will not make it visible in this example, though.
	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	image, _, err := image.Decode(bytes.NewReader(assets.PotatoPng))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	iconbt, err := walk.NewBitmapFromImageForDPI(image, 96)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// We load our icon from a file.
	icon, err := walk.NewIconFromBitmap(iconbt)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// Create the notify icon and make sure we clean it up on exit.
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer ni.Dispose()

	// Set the icon and a tool tip text.
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := ni.SetToolTip("Click for info or use the context menu to exit."); err != nil {
		log.Fatal().Err(err).Send()
	}

	// When the left mouse button is pressed, bring up our balloon.
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		if err := ni.ShowCustom(
			"Walk NotifyIcon Example",
			"There are multiple ShowX methods sporting different icons.",
			icon); err != nil {

			log.Fatal().Err(err).Send()
		}
	})

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal().Err(err).Send()
	}
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal().Err(err).Send()
	}

	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		log.Fatal().Err(err).Send()
	}

	// Now that the icon is visible, we can bring up an info balloon.
	if err := ni.ShowInfo("Walk NotifyIcon Example", "Click the icon to show again."); err != nil {
		log.Fatal().Err(err).Send()
	}

	// Run the message loop.
	mw.Run()
}
