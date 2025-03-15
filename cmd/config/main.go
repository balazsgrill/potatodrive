package main

import (
	"github.com/balazsgrill/potatodrive/bindings"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/rs/zerolog/log"
)

func main() {
	var listBox *walk.ListBox
	var form *walk.Composite
	var configProvider bindings.ConfigWriter

	configProvider = bindings.NewRegistryConfigWriter(log.Logger, "SOFTWARE\\PotatoDrive")

	// Read registry values
	items := configProvider.Keys()

	MainWindow{
		Title:   "PotatoDrive configuration",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		ToolBar: ToolBar{
			Items: []MenuItem{
				Action{
					Text: "Refresh",
					OnTriggered: func() {
						items = configProvider.Keys()
						listBox.SetModel(items)
					},
				},
				Action{
					Text: "Exit",
					OnTriggered: func() {
						walk.App().Exit(0)
					},
				},
			},
		},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					ListBox{
						AssignTo:       &listBox,
						Model:          items,
						MultiSelection: false,
						OnSelectedIndexesChanged: func() {
							if listBox.CurrentIndex() == -1 {
								return
							}
							key := items[listBox.CurrentIndex()]
							config, err := configProvider.ReadConfig(key)
							if err != nil {
								log.Printf("Failed to read config: %v", err)
								return
							}
							form.Children().At(1).(*walk.LineEdit).SetText(config.BaseConfig.LocalPath)
							form.Children().At(3).(*walk.LineEdit).SetText(config.BaseConfig.API)

						},
					},
					Composite{
						AssignTo: &form,
						Layout:   VBox{},
						DataBinder: DataBinder{
							DataSource: bindings.Config{},
						},
						Children: []Widget{
							Label{Text: "LocalPath:"},
							LineEdit{},
							Label{Text: "Form Field 2:"},
							LineEdit{},
						},
					},
				},
			},
		},
	}.Run()
}
