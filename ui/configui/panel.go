package configui

import (
	"log"

	"github.com/balazsgrill/potatodrive/bindings"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func ConfigPanel(databinder **walk.DataBinder, configProvider bindings.ConfigWriter) Widget {
	return Composite{
		Layout: VBox{},
		DataBinder: DataBinder{
			AssignTo: databinder,
			DataSource: &ConfigValues{
				HasValue:    false,
				NotHasValue: true,
			},
			OnSubmitted: func() {
				configvalues := (*databinder).DataSource().(*ConfigValues)
				value := WriteTo(configvalues)
				if value != nil {
					configProvider.WriteConfig(*value)
				}
			},
			AutoSubmit: true,
		},
		Children: []Widget{
			Label{
				Visible: Bind("NotHasValue"),
				Text:    "No configuration selected",
				MinSize: Size{500, 20},
			},
			Composite{
				Visible: Bind("HasValue"),
				Layout:  Grid{Columns: 3},
				Children: []Widget{
					Label{Text: "id:"},
					Label{Text: Bind("ID"), ColumnSpan: 2},
					Label{Text: "LocalPath:"},
					LineEdit{Text: Bind("Base.LocalPath")},
					PushButton{
						Text: "Browse",
						OnClicked: func() {
							dlg := new(walk.FileDialog)
							dlg.Title = "Select folder"
							dlg.Filter = "Folders|*"
							if ok, err := dlg.ShowBrowseFolder(nil); err != nil {
								log.Printf("Failed to show dialog: %v", err)
								return
							} else if !ok {
								return
							}
							(*databinder).DataSource().(*ConfigValues).Base.LocalPath = dlg.FilePath
							(*databinder).Reset()
						},
					},
					Label{Text: "Filesystem integration API:"},
					ComboBox{
						Value:      Bind("Base.API"),
						ColumnSpan: 2,
						Model:      []string{bindings.APIType_CFAPI, bindings.APIType_PRJFS, bindings.APIType_CFAPI_Simplfied},
					},
				},
			},
			Composite{
				Visible: Bind("HasS3"),
				Layout:  Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "Endpoint:"},
					LineEdit{Text: Bind("S3Config.Endpoint")},
					Label{Text: "Region:"},
					LineEdit{Text: Bind("S3Config.Region")},
					Label{Text: "Bucket:"},
					LineEdit{Text: Bind("S3Config.Bucket")},
					Label{Text: "KeyId:"},
					LineEdit{Text: Bind("S3Config.KeyId")},
					Label{Text: "KeySecret:"},
					LineEdit{Text: Bind("S3Config.KeySecret")},
					Label{Text: "UseSSL:"},
					CheckBox{Checked: Bind("S3Config.UseSSL")},
				},
			},
			Composite{
				Visible: Bind("HasSFTP"),
				Layout:  Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "Host:"},
					LineEdit{Text: Bind("SFTPConfig.Host")},
					Label{Text: "Base path:"},
					LineEdit{Text: Bind("SFTPConfig.Basepath")},
					Label{Text: "User:"},
					LineEdit{Text: Bind("SFTPConfig.User")},
					Label{Text: "Password:"},
					LineEdit{Text: Bind("SFTPConfig.Password")},
					Label{Text: "PrivateKey:"},
					LineEdit{Text: Bind("SFTPConfig.PrivateKey")},
				},
			},
			Composite{
				Visible: Bind("HasGPhotos"),
				Layout:  Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "ClientId:"},
					LineEdit{Text: Bind("GPhotosConfig.ClientId")},
					Label{Text: "ClientSecret:"},
					LineEdit{Text: Bind("GPhotosConfig.ClientSecret")},
					Label{Text: "RedirectURL:"},
					LineEdit{Text: Bind("GPhotosConfig.RedirectURL")},
				},
			},
		},
	}
}
