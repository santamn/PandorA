package main

import (
	"os/exec"
	"pandora/pkg/account"
	"pandora/pkg/resource"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

func main() {
	// メニューバーを起動
	systray.Run(menuReady, menuExit)
}

// menuReady メニューを初期化する
func menuReady() {
	// メニューバーにタブを設定
	systray.SetTitle("PandorA")
	download := systray.AddMenuItem("Download", "Download resources in PandA")
	settings := systray.AddMenuItem("Settings", "Settings")
	quit := systray.AddMenuItem("Quit", "Quit PandorA")

	windowExist := false

	for {
		select {
		case <-download.ClickedCh:
			ecsID, password, rejectable, err := account.ReadAccountInfo()
			if err != nil {
				beeep.Alert("PandorA Error", "Failed to read account info", "")
			}
			resource.Download(ecsID, password, rejectable)

		case <-settings.ClickedCh:
			if !windowExist {
				// 画面が二つ以上表示されないようにする
				windowExist = true
				if err := exec.Command("../form/form").Run(); err != nil {
					beeep.Alert("PandorA Error", err.Error(), "")
				}
				windowExist = false
			}

		case <-quit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// menuExit メニューを終了する
func menuExit() {}
