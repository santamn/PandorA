package main

import (
	"os/exec"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

const (
	pathToForm = "../form/form"
)

var (
	windowExist       bool
	windowMu          sync.Mutex
	downloadClickedCh chan struct{}
)

func init() {
	downloadClickedCh = make(chan struct{})
}

func main() {
	// メニューバーを起動
	go systray.Run(menuReady, menuExit)

	// 4時間おきにダウンロードを実行
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	var d downloadManager

	for {
		select {
		case <-downloadClickedCh: // ボタンを押された場合の実行
			go d.excute(true)
		case <-ticker.C: // 定期実行
			go d.excute(false)
		}
	}
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
			showWindow()

		case <-quit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// menuExit メニューを終了する
func menuExit() {}

// showWindow ユーザー情報を入力するウィンドウを起動
// どのスレッドから呼ばれても画面が一つしか表示されないようにする
func showWindow() {
	path := "../form/form"

	windowMu.Lock()
	if !windowExist {
		windowExist = true
		windowMu.Unlock()
		// UIを起動
		if err := exec.Command(path).Run(); err != nil {
			beeep.Alert("PandorA Error", err.Error(), "")
		}
		windowExist = false
	} else {
		windowMu.Unlock()
	}
}
