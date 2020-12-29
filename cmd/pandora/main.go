package main

import (
	"os/exec"
	"pandora/pkg/account"
	pandaapi "pandora/pkg/pandaAPI"
	"pandora/pkg/resource"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

const (
	pathToForm = "../form/form"
)

var (
	windowExist bool
	windowMu    sync.Mutex
)

func main() {
	// メニューバーを起動
	go systray.Run(menuReady, menuExit)

	// 4時間おきにダウンロードを実行
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ecsID, password, rejectable, err := account.ReadAccountInfo()
			if err != nil {
				// アカウント情報を入力させる
				showWindow()
			}
			ecsID, password, rejectable, err = account.ReadAccountInfo()
			if err != nil {
				// 2回目にエラーが出た場合はエラーを表示して終了する
				beeep.Alert("PandorA Error", err.Error(), "")
				return
			}

			if err := resource.Download(ecsID, password, rejectable); err != nil {
				switch err.(type) {
				case *pandaapi.NetworkError:
					beeep.Alert("PandorA Error", "Network Error: something wrong with connecting the Internet", "")
				case *pandaapi.DeadPandAError:
					beeep.Alert("PandorA Error", err.Error(), "")
				case *pandaapi.FailedLoginError:
					beeep.Alert("PandorA Error", err.Error(), "")
					showWindow()
				default:
					beeep.Alert("PandorA Error", "System Error: "+err.Error(), "")
				}
			}
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
			ecsID, password, rejectable, err := account.ReadAccountInfo()
			if err != nil {
				beeep.Alert("PandorA Error", "Failed to read account info", "")
			}
			resource.Download(ecsID, password, rejectable)

		case <-settings.ClickedCh:
			showWindow()

		case <-quit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// menuExit メニューを終了する
func menuExit() {
}

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
