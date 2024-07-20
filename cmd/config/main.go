package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	binding.NewUntypedList()

	//list := widget.NewListWithData()
	//left := container.NewVBox(list)
	//editor := container.New(layout.NewFormLayout())

	w.SetContent(container.NewBorder(nil, nil, nil, nil))
	w.ShowAndRun()
}
