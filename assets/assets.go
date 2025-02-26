package assets

import (
	"embed"
	_ "embed"
)

//go:embed potato.png
var PotatoPng []byte

//go:embed locales/*
var Locales embed.FS
