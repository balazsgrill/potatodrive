package assets

import (
	"embed"
	_ "embed"
)

//go:embed potato.png
var PotatoPng []byte

//go:embed success.html
var AuthSuccessHtml []byte

//go:embed locales/*
var Locales embed.FS

//go:embed icons/*
var Icons embed.FS

const IconBucket = "icons/bucket.png"
const IconSFTP = "icons/sftp.png"
const IconGPhotos = "icons/gphotos.png"
const IconDelete = "icons/delete.png"
const IconRefresh = "icons/refresh.png"
