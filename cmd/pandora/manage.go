package main

import (
	"fmt"
	"pandora/pkg/account"
	pandaapi "pandora/pkg/pandaAPI"
	"pandora/pkg/resource"
	"sync"
	"time"
)

type downloadManager struct {
	isRunning        bool
	lastExecutedTime time.Time
	mu               sync.Mutex
}

func (download *downloadManager) excute(clicked bool) {
	download.mu.Lock()
	if !download.isRunning {
		download.isRunning = true
		download.mu.Unlock()

		if min := time.Now().Sub(download.lastExecutedTime).Minutes(); min < 10 && clicked {
			// 前のダウンロードからの経過時間が10分以内にユーザーによる再度の実行の要求があれば警告を出して終了する
			alert(fmt.Sprintf("PandorA needs cool time. Please try after %d min later at least", uint(10-min)))
			download.isRunning = false
			return
		}

		ecsID, password, rejectable, err := account.ReadAccountInfo()
		if err != nil {
			// アカウント情報を入力させる
			showWindow() // TODO:既にウィンドウが開かれている状態だとここでアカウント情報が修正されるのを待つことができない
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
				alert("Network Error: something wrong with connecting the Internet") //TODO:ネッt障害のエラーをどこかにログとして出力して置くかどうか考える
			case *pandaapi.DeadPandAError:
				alert(err.Error())
			case *pandaapi.FailedLoginError:
				alert(err.Error())
				showWindow()
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
