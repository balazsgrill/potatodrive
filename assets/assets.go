package assets

import (
	"embed"
	_ "embed"
)

//go:embed potato.png
var PotatoPng []byte

//go:embed locales/*
var Locales embed.FS

//go:embed icons/*
var Icons embed.FS

const IconBucket = "icons/bucket.png"
const IconSFTP = "icons/sftp.png"
