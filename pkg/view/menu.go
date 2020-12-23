package view

import (
	"pandora/pkg/account"
	"pandora/pkg/resource"

	"github.com/gen2brain/beeep"
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
			ecsID, password, rejectable, err := account.ReadAccountInfo()
			if err != nil {
				beeep.Alert("PandorA Error", "Failed to read account info", "")
			}
			resource.Download(ecsID, password, rejectable)

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
