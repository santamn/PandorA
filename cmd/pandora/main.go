package main

import (
	"io"
	"log"
	"os"
	"time"

	"pandora/cmd/pandora/icon"

	"github.com/getlantern/systray"
)

func main() {
	// [DEBUG](https://qiita.com/74th/items/441ffcab80a6a28f7ee3)
	logfile, err := os.OpenFile("./test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open test.log:" + err.Error())
	}
	defer logfile.Close()

	// io.MultiWriteで標準出力とファイルの両方を束ねて、logの出力先に設定する
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
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

	download := new(downloadManager)
	window, err := newWindowManager()
	// [DEBUG]
	if err != nil {
		log.Println(err)
		return
	}

	for { // TODO:goroutineの終了を待たなくて良いのか？
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
func menuExit() {}
