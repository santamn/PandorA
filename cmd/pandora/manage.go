package main

import (
	"fmt"
	"os/exec"
	"pandora/pkg/account"
	pandaapi "pandora/pkg/pandaAPI"
	"pandora/pkg/resource"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
)

// downloadManager ダウンロード実行中に並列して実行されたり、短いタイムスパンでダウンロードが実行されないように制御する
type downloadManager struct {
	isRunning        bool
	lastExecutedTime time.Time
	mu               sync.Mutex
}

func (download *downloadManager) excute(window *windowManager, clicked bool) {
	download.mu.Lock()
	if !download.isRunning {
		download.isRunning = true
		download.mu.Unlock()

		if min := time.Now().Sub(download.lastExecutedTime).Minutes(); min < 10 && clicked {
			// 前のダウンロードからの経過時間が10分以内にユーザーによる再度の実行の要求があれば警告を出して終了する
			alert(fmt.Sprintf("PandorA needs cool time. Please try after %d minute(s) later at least.", uint(10-min)))
			download.isRunning = false
			return
		}

		ecsID, password, rejectable, err := account.ReadAccountInfo()
		if err != nil {
			// アカウント情報を入力させる
			window.show()
			window.wg.Wait() // アカウント情報の入力を待つ
		}
		ecsID, password, rejectable, err = account.ReadAccountInfo()
		if err != nil {
			// 2回目にエラーが出た場合はエラーを表示して終了する
			alert(err.Error())
			return
		}

		notify("NOW DOWNLOADING")

		download.lastExecutedTime = time.Now()
		if err := resource.Download(ecsID, password, rejectable); err != nil {
			download.isRunning = false

			switch err.(type) {
			case *pandaapi.NetworkError:
				alert("Network Error: something wrong with connecting the Internet") //TODO:ネット障害のエラーをどこかにログとして出力して置くかどうか考える
			case *pandaapi.DeadPandAError:
				alert(err.Error())
			case *pandaapi.FailedLoginError:
				alert(err.Error())
				window.show()
			default:
				alert("System Error: " + err.Error())
			}
		} else {
			download.isRunning = false
			notify("Download succeeded!")
		}
	} else {
		download.mu.Unlock()
	}
}

// windowManager ウィンドウが画面に一つだけ表示されるよう管理する
type windowManager struct {
	isShowing bool
	mu        sync.Mutex
	wg        sync.WaitGroup
}

func (w *windowManager) show() {
	path := "../form/form"

	w.mu.Lock()
	if !w.isShowing {
		w.isShowing = true
		w.wg.Add(1)
		w.mu.Unlock()

		// UIを起動
		if err := exec.Command(path).Run(); err != nil {
			beeep.Alert("PandorA Error", err.Error(), "")
		}

		w.isShowing = false
		// wg.Waitを使えばここで画面が終了することを待つことができる
		w.wg.Done()
	} else {
		w.mu.Unlock()
	}
}
