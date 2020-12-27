package main

import (
	"log"
	"os/exec"
	"pandora/pkg/account"
	pandaapi "pandora/pkg/pandaAPI"
	"pandora/pkg/resource"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

const (
	pathToForm = "../form/form"
)

var (
	showCh chan struct{}
)

func init() {
	showCh = make(chan struct{})
}

func main() {
	// メニューバーを起動
	go systray.Run(menuReady, menuExit)

	// ユーザー情報を設定するウィンドウの起動を監視
	go windowManager("../form/form", showCh)

	// 4時間おきにダウンロードを実行
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ecsID, password, rejectable, err := account.ReadAccountInfo()
			if err != nil {
				// アカウント情報を入力させる
				showCh <- struct{}{}
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
				case *pandaapi.DeadPandAError:
				case *pandaapi.FailedLoginError:
				default:
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
			showCh <- struct{}{}

		case <-quit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// menuExit メニューを終了する
func menuExit() {}

// windowManager ユーザー情報を入力するウィンドウの起動を監視
func windowManager(path string, showCh <-chan struct{}) {
	cmd := exec.Command(path)
	windowExist := false

	for range showCh {
		log.Println("はじめ")
		if !windowExist {
			log.Println("画面")
			windowExist = true
			// UIを起動
			if err := cmd.Run(); err != nil {
				beeep.Alert("PandorA Error", err.Error(), "")
			}
			windowExist = false
		}
		log.Println("終わり")
	}
}
