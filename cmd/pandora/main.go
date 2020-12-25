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

var (
	pandora       fyne.App
	window        fyne.Window
	menuRestartCh chan struct{}
	quitCh        chan struct{}
)

func init() {
	menuRestartCh = make(chan struct{})
	quitCh = make(chan struct{})
}

func main() {
	// アカウント情報を入力するウィンドウを作成
	pandora = app.New()
	pandora.Settings().SetTheme(theme.DarkTheme())
	window = pandora.NewWindow("PandorA")
	object := view.MakeForm(window)
	window.Resize(fyne.NewSize(250, 100))
	window.SetContent(object)
	window.SetOnClosed(window.Hide)

	// メニューバーを起動
	go systray.Run(menuReady, menuExit)

	for {
		select {
		case <-quitCh:
			return
		default:
			log.Println("画面起動")
			pandora.Run()
		}
	}

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
			log.Print("画面出力")

		case <-logButton.ClickedCh:
			log.Println("ログ出力")

		case <-quit.ClickedCh:
			pandora.Quit()
			systray.Quit()
			quitCh <- struct{}{}
			return
		}
	}
}

// menuExit メニューを終了する
func menuExit() {}
