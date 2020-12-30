package main

import (
	"sync"
	"time"

	"github.com/getlantern/systray"
)

var (
	downloadClickedCh chan struct{}
	window            *windowManager
)

func init() {
	downloadClickedCh = make(chan struct{})
}

func main() {
	var wg sync.WaitGroup
	// メニューバーを起動
	wg.Add(1)
	go func() {
		defer wg.Done()
		systray.Run(menuReady, menuExit)
	}()

	// 4時間おきにダウンロードを実行
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	var d downloadManager
	for {
		select {
		case <-downloadClickedCh: // ダウンロードボタンを押された場合の実行
			wg.Add(1)
			go func() {
				defer wg.Done()
				d.excute(true)
			}()

		case <-ticker.C: // 定期実行
			wg.Add(1)
			go func() {
				defer wg.Done()
				d.excute(false)
			}()
		}
	}
	wg.Wait()
}

// menuReady メニューを初期化する
func menuReady() {
	// メニューバーにタブを設定
	systray.SetTitle("PandorA")
	download := systray.AddMenuItem("Download", "Download resources in PandA")
	settings := systray.AddMenuItem("Settings", "Settings")
	quit := systray.AddMenuItem("Quit", "Quit PandorA")

	for {
		select {
		case <-download.ClickedCh:
			downloadClickedCh <- struct{}{}

		case <-settings.ClickedCh:
			window.show()

		case <-quit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// menuExit メニューを終了する
func menuExit() {}
