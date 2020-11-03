package view

import (
	"github.com/getlantern/systray"
)

// MenuReady メニューを初期化する
func MenuReady() {
	systray.SetTitle("PandorA")
	download := systray.AddMenuItem("Download", "Download resources in PandA")
	quit := systray.AddMenuItem("Quit", "Quit PandorA")
	settings := systray.AddMenuItem("Settings", "Settings")

	for {
		select {
		case <-download.ClickedCh:
		case <-settings.ClickedCh:
		case <-quit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// MenuExit メニューを終了する
func MenuExit() {

}
