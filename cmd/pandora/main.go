package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"pandora/cmd/pandora/icon"
	"pandora/pkg/dir"

	"github.com/getlantern/systray"
)

var (
	window   *windowManager
	download *downloadManager
)

func init() {
	window = newWindowManager()
	download = &downloadManager{}
}

func main() {
	// ログ出力を設定
	logfile, err := os.OpenFile(
		filepath.Join(dir.WorkingDirecory, "pandoraError.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		panic("cannnot open pandoraError.log:" + err.Error())
	}
	defer logfile.Close()

	log.SetOutput(logfile)
	log.SetFlags(log.Ldate | log.Ltime)

	systray.Run(menuReady, menuExit)
}

// menuReady メニューを初期化する
func menuReady() {
	// メニューバーにタブを設定
	systray.SetTitle("PandorA")
	systray.SetIcon(icon.Data)
	downloadButton := systray.AddMenuItem("Download", "Download resources in PandA")
	settingsButton := systray.AddMenuItem("Settings", "Settings")
	quitButton := systray.AddMenuItem("Quit", "Quit PandorA")

	// 4時間おきにダウンロードを実行
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go download.excute(window, false)

		case <-downloadButton.ClickedCh:
			go download.excute(window, true)

		case <-settingsButton.ClickedCh:
			go window.show()

		case <-quitButton.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// menuExit メニューを終了する
func menuExit() {
	window.quit()
	download.wg.Wait()
}
