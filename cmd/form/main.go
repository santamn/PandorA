package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
)

func main() {
	// アカウント情報を入力するウィンドウを作成
	pandora := app.New()
	pandora.Settings().SetTheme(theme.DarkTheme())
	window := pandora.NewWindow("PandorA")
	object := makeForm(window)
	window.Resize(fyne.NewSize(400, 200))
	window.SetContent(object)
	window.ShowAndRun()
}
