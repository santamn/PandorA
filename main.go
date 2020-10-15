package main

import (
	"fmt"
	"pandora/data"
)

func main() {
	// pandora := app.New()
	// pandora.Settings().SetTheme(theme.DarkTheme())

	// window := pandora.NewWindow("PandorA")
	// object := screens.SemesterScreen(window)

	// window.Resize(fyne.NewSize(700, 500))
	// window.SetContent(object)
	// window.ShowAndRun()

	client, err := data.NewLoggedInClient("a0180935", "SCP-8900-ex")
	if err != nil {
		fmt.Println(err)
	}

	if err := data.DownloadPDF(client, "2020-888-H730-002"); err != nil {
		fmt.Println(err)
	}
}
