package main

import (
	"pandora/pkg/view"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
)

func main() {
	pandora := app.New()
	pandora.Settings().SetTheme(theme.DarkTheme())

	window := pandora.NewWindow("PandorA")
	object := view.MakeForm(window)

	window.Resize(fyne.NewSize(250, 100))
	window.SetContent(object)
	pandora.Run()
}
