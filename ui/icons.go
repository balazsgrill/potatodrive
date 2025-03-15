package ui

import (
	"image"

	"github.com/balazsgrill/potatodrive/assets"
	"github.com/lxn/walk"
)

func (context *UIContext) GetImageForAsset(name string) walk.Image {
	file, err := assets.Icons.Open(name)
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		context.Logger.Fatal().Err(err).Send()
	}
	bt, err := walk.NewBitmapFromImageForDPI(image, 96)
	if err != nil {
		context.Logger.Fatal().Err(err).Send()
	}
	return bt
}
