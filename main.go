package main

import "github.com/getlantern/systray"

func main() {
	// pandora := app.New()
	// pandora.Settings().SetTheme(theme.DarkTheme())

	// window := pandora.NewWindow("PandorA")
	// object := screens.SemesterScreen(window)

	// window.Resize(fyne.NewSize(700, 500))
	// window.SetContent(object)
	// window.ShowAndRun()

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome")
	systray.AddMenuItem("Quit", "Quit the whole app")
}

func onExit() {}
