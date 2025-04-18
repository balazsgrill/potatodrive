package configui

import (
	"log"

	"github.com/balazsgrill/potatodrive/assets"
	"github.com/balazsgrill/potatodrive/bindings"
	"github.com/balazsgrill/potatodrive/ui"
	"github.com/google/uuid"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type configUI struct {
	mw *walk.MainWindow
}

func Create(uicontext *ui.UIContext, configProvider bindings.ConfigWriter) MainWindow {
	var listBox *walk.ListBox
	var db *walk.DataBinder
	cui := &configUI{}
	// Read registry values
	items := configProvider.Keys()
	refresh := func() {
		items = configProvider.Keys()
		listBox.SetModel(items)
	}
	return MainWindow{
		AssignTo: &cui.mw,
		Title:    "PotatoDrive configuration",
		MinSize:  Size{600, 400},
		Layout:   VBox{},
		ToolBar: ToolBar{
			ButtonStyle: ToolBarButtonImageAboveText,
			Items: []MenuItem{
				Action{
					Text:        "Refresh",
					Image:       uicontext.GetImageForAsset(assets.IconRefresh),
					OnTriggered: refresh,
				},
				Action{
					Text:  "Mount S3",
					Image: uicontext.GetImageForAsset(assets.IconBucket),
					OnTriggered: func() {
						db.SetDataSource(&ConfigValues{
							ID: uuid.NewString(),
							Base: bindings.BaseConfig{
								Type: bindings.TYPE_S3,
								API:  bindings.APIType_CFAPI,
							},
							HasValue: true,
							HasS3:    true,
						})
						db.Reset()
						refresh()
					},
				},
				Action{
					Text:  "Mount SFTP",
					Image: uicontext.GetImageForAsset(assets.IconSFTP),
					OnTriggered: func() {
						db.SetDataSource(&ConfigValues{
							ID: uuid.NewString(),
							Base: bindings.BaseConfig{
								Type: bindings.TYPE_SFTP,
								API:  bindings.APIType_CFAPI,
							},
							HasValue: true,
							HasSFTP:  true,
						})
						db.Reset()
						refresh()
					},
				},
				Action{
					Text:  "Mount GPhotos",
					Image: uicontext.GetImageForAsset(assets.IconGPhotos),
					OnTriggered: func() {
						db.SetDataSource(&ConfigValues{
							ID: uuid.NewString(),
							Base: bindings.BaseConfig{
								Type: bindings.TYPE_GPHOTOS,
								API:  bindings.APIType_CFAPI,
							},
							HasValue:    true,
							HasGPhotos:  true,
							NotHasValue: false,
						})
						db.Reset()
						refresh()
					},
				},
				Action{
					Text:  "Delete",
					Image: uicontext.GetImageForAsset(assets.IconDelete),
					OnTriggered: func() {
						if listBox.CurrentIndex() == -1 {
							return
						}
						key := items[listBox.CurrentIndex()]
						if err := configProvider.DeleteConfig(key); err != nil {
							log.Printf("Failed to delete config: %v", err)
							return
						}
						refresh()
					},
				},
			},
		},
		Children: []Widget{
			HSplitter{
				Persistent: true,
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
								// log the error but do not return early
								// this way the user can still edit the config
							}
							// Set the data source even if there is an error
							db.SetDataSource(ReadFrom(&config))
							db.Reset()

							listBox.Parent().RequestLayout()
						},
					},
					cui.configPanel(&db, configProvider),
				},
			},
		},
	}
}
