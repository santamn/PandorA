package main

import (
	"log"
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
			exec.Command("./form").Run()
			log.Print("画面出力")

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
