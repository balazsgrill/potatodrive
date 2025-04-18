package ui

import (
	"fmt"

	"github.com/balazsgrill/potatodrive/core"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/rs/zerolog/log"
)

func openlink(link *walk.LinkLabelLink) {
	err := core.OpenURL(link.URL())
	if err != nil {
		log.Error().Err(err).Msg("Failed to open url")
	}
}

func aboutDialog(_ walk.Form, context *UIContext) (int, error) {
	return MainWindow{
		Title:   context.Get("About PotatoDrive"),
		Icon:    "#2\\0409",
		MinSize: Size{300, 200},
		Size:    Size{300, 200},
		MaxSize: Size{300, 200},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Alignment: AlignHCenterVCenter,
				Text:      fmt.Sprintf("PotatoDrive %s", context.Version),
			},
			LinkLabel{
				Alignment:       AlignHCenterVCenter,
				Text:            context.Get("License"),
				OnLinkActivated: openlink,
			},
			LinkLabel{
				Alignment:       AlignHCenterVCenter,
				Text:            context.Get("Homepage"),
				OnLinkActivated: openlink,
			},
		},
	}.Run()
}
