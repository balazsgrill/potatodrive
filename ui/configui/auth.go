package configui

import (
	"context"
	"log"

	"github.com/balazsgrill/potatodrive/bindings/gphotos"
	"github.com/balazsgrill/potatodrive/core"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func auth(ownerWindow walk.Form, config *gphotos.Config) error {
	authurlchan := make(chan string)
	var fdlg *walk.Dialog
	var cancelButton *walk.PushButton

	if config == nil {
		return nil
	}
	if config.ClientID == "" || config.ClientSecret == "" || config.RedirectURL == "" {
		return nil
	}
	authcontext, cancelfunc := context.WithCancel(context.Background())
	go func() {
		err := config.Authenticate(authcontext, func(url string) {
			authurlchan <- url
			err := core.OpenURL(url)
			if err != nil {
				log.Printf("Error opening URL: %v", err)
			}
		})
		if err != nil {
			close(authurlchan)
			log.Printf("Error during authentication: %v", err)
		} else {
			fdlg.Close(0)
		}
	}()
	authurl, ok := <-authurlchan
	close(authurlchan)
	if !ok {
		log.Println("Authentication failed")
		cancelfunc()
		return nil
	}

	dlg := declarative.Dialog{
		AssignTo:     &fdlg,
		Title:        "Authentication",
		MinSize:      declarative.Size{Width: 300, Height: 100},
		Layout:       declarative.VBox{Alignment: declarative.AlignHCenterVCenter},
		CancelButton: &cancelButton,
		Children: []declarative.Widget{
			declarative.Label{
				Text: "Google authentication is in progress",
			},
			declarative.LinkLabel{
				Text: authurl,
			},
			declarative.PushButton{
				AssignTo: &cancelButton,
				Text:     "Cancel",
				OnClicked: func() {
					cancelfunc()
					fdlg.Cancel()
				},
			},
		},
	}

	dlg.Run(ownerWindow)
	return nil
}
