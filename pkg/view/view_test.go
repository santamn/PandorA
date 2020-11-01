package view_test

import (
	"pandora/pkg/view"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func TestView(t *testing.T) {
	pandora := test.NewApp()
	pandora.Settings().SetTheme(theme.DarkTheme())

	window := pandora.NewWindow("PandorA")
	object := view.MakeUserForm(window)

	window.Resize(fyne.NewSize(700, 500))
	window.SetContent(object)
	window.ShowAndRun()
}
