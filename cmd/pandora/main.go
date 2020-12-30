package main

import (
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(menuReady, menuExit)
}

// menuReady メニューを初期化する
func menuReady() {
	// メニューバーにタブを設定
	systray.SetTitle("PandorA")
	downloadButton := systray.AddMenuItem("Download", "Download resources in PandA")
	settingsButton := systray.AddMenuItem("Settings", "Settings")
	quitButton := systray.AddMenuItem("Quit", "Quit PandorA")

	// 4時間おきにダウンロードを実行
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	download := new(downloadManager)
	window := new(windowManager)

	for {
		select {
		case <-ticker.C:
			download.excute(window, false)

		case <-downloadButton.ClickedCh:
			download.excute(window, true)

		case <-settingsButton.ClickedCh:
			window.show()

		case <-quitButton.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// menuExit メニューを終了する
func menuExit() {}
