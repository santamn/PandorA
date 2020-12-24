package main

import (
	"log"
	"pandora/pkg/account"
	"pandora/pkg/resource"
	"pandora/pkg/view"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

var window fyne.Window

func main() {
	// アカウント情報を記入するフォームを作成
	pandora := app.New()
	pandora.Settings().SetTheme(theme.DarkTheme())
	window = pandora.NewWindow("PandorA")
	object := view.MakeForm(window)
	window.Resize(fyne.NewSize(250, 100))
	window.SetContent(object)

	go systray.Run(menuReady, menuExit)

	pandora.Run()

	log.Println("終了")
}

// menuReady メニューを初期化する
func menuReady() {
	// メニューバーにタブを設定
	systray.SetTitle("PandorA")
	download := systray.AddMenuItem("Download", "Download resources in PandA")
	settings := systray.AddMenuItem("Settings", "Settings")
	logButton := systray.AddMenuItem("Log", "Print log")
	quit := systray.AddMenuItem("Quit", "Quit PandorA")

	for {
		select {
		case <-download.ClickedCh:
			ecsID, password, rejectable, err := account.ReadAccountInfo()
			if err != nil {
				beeep.Alert("PandorA Error", "Failed to read account info", "")
			}
			resource.Download(ecsID, password, rejectable)

		case <-settings.ClickedCh:
			window.Show()

		case <-logButton.ClickedCh:
			log.Println("ログ出力")

		case <-quit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// menuExit メニューを終了する
func menuExit() {}
