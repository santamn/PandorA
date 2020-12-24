package view

import (
	"pandora/pkg/account"
	"pandora/pkg/resource"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

// MenuReady メニューを初期化する
func MenuReady() {
	// アカウント情報を記入するフォームを作成
	pandora := app.New()
	pandora.Settings().SetTheme(theme.DarkTheme())
	window := pandora.NewWindow("PandorA")
	object := MakeForm(window)
	window.Resize(fyne.NewSize(250, 100))
	window.SetContent(object)

	// メニューバーにタブを設定
	systray.SetTitle("PandorA")
	download := systray.AddMenuItem("Download", "Download resources in PandA")
	quit := systray.AddMenuItem("Quit", "Quit PandorA")
	settings := systray.AddMenuItem("Settings", "Settings")

	for {
		select {
		case <-download.ClickedCh:
			ecsID, password, rejectable, err := account.ReadAccountInfo()
			if err != nil {
				beeep.Alert("PandorA Error", "Failed to read account info", "")
			}
			resource.Download(ecsID, password, rejectable)

		case <-settings.ClickedCh:
			window.ShowAndRun()

		case <-quit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// MenuExit メニューを終了する
func MenuExit() {}
