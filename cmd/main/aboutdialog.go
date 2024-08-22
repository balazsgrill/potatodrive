package main

import (
	"fmt"
	"os/exec"

	"github.com/balazsgrill/potatodrive/win"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/rs/zerolog/log"
)

func openlink(link *walk.LinkLabelLink) {
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", link.URL()).Start()
	if err != nil {
		log.Error().Err(err).Msg("Failed to open url")
	}
}

func aboutDialog(_ walk.Form) (int, error) {
	return MainWindow{
		Title:   "About PotatoDrive",
		Icon:    win.InstalledFile(win.POTATOICO),
		MinSize: Size{300, 200},
		Size:    Size{300, 200},
		MaxSize: Size{300, 200},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Alignment: AlignHCenterVCenter,
				Text:      fmt.Sprintf("PotatoDrive %s", Version),
			},
			LinkLabel{
				Alignment:       AlignHCenterVCenter,
				Text:            `License: <a id="this" href="https://github.com/balazsgrill/potatodrive?tab=MIT-1-ov-file#readme">MIT</a>`,
				OnLinkActivated: openlink,
			},
			LinkLabel{
				Alignment:       AlignHCenterVCenter,
				Text:            `<a id="this" href="https://github.com/balazsgrill/potatodrive">https://github.com/balazsgrill/potatodrive</a>`,
				OnLinkActivated: openlink,
			},
		},
	}.Run()
}
