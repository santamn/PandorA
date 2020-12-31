package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"pandora/pkg/account"
	pandaapi "pandora/pkg/pandaAPI"
	"pandora/pkg/resource"
	"path/filepath"
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
	defer func() {
		download.isRunning = false
	}()

	download.mu.Lock()
	if !download.isRunning {
		download.isRunning = true
		download.mu.Unlock()

		if min := time.Now().Sub(download.lastExecutedTime).Minutes(); min < 10 && clicked {
			// 前のダウンロードからの経過時間が10分以内にユーザーによる再度の実行の要求があれば警告を出して終了する
			alert(fmt.Sprintf("PandorA needs cool time. Please try after %d minute(s) later at least.", uint(10-min)))
			return
		}

		ecsID, password, rejectable, err := account.ReadAccountInfo()
		if err != nil {
			// [DEBUG]
			log.Println("read account error 1", err)
			// アカウント情報を入力させる
			window.show()
			window.wg.Wait() // アカウント情報の入力を待つ
		}
		ecsID, password, rejectable, err = account.ReadAccountInfo()
		if err != nil {
			// [DEBUG]
			log.Println("read account error 2", err)
			// 2回目にエラーが出た場合はエラーを表示して終了する
			alert(err.Error())
			return
		}

		notify("NOW DOWNLOADING")

		download.lastExecutedTime = time.Now()
		if err := resource.Download(ecsID, password, rejectable); err != nil {
			// [DEBUG]
			log.Println("Download error", err)

			switch err.(type) {
			case *pandaapi.NetworkError:
				alert("Network Error: something wrong with connecting the Internet") //TODO:ネット障害のエラーをどこかにログとして出力して置くかどうか考える
			case *pandaapi.DeadPandAError:
				alert(err.Error())
			case *pandaapi.FailedLoginError:
				alert(err.Error())
				go window.show()
			default:
				alert("System Error: " + err.Error())
			}
		} else {
			notify("Download succeeded!")
		}
	} else {
		download.mu.Unlock()
	}
}

// windowManager ウィンドウが画面に一つだけ表示されるよう管理する
type windowManager struct {
	cmd       *exec.Cmd
	isShowing bool
	mu        sync.Mutex
	path      string
	wg        sync.WaitGroup
}

func newWindowManager() (*windowManager, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(filepath.Dir(exe), "form")

	return &windowManager{
		path: path,
	}, nil
}

func (w *windowManager) show() {
	w.mu.Lock()
	if !w.isShowing {
		w.isShowing = true
		w.wg.Add(1)
		w.mu.Unlock()

		// UIを起動
		w.cmd = exec.Command(w.path)
		if err := w.cmd.Run(); err != nil {
			// [DEBUG]
			log.Println("show error:", err)
			beeep.Alert("PandorA Error", err.Error(), "")
		}
		w.cmd = nil
		w.isShowing = false
		// wg.Waitを使えばここで画面が終了することを待つことができる
		w.wg.Done()
	} else {
		w.mu.Unlock()
	}
}

// ウィンドウを表示しているプロセスをKillする
func (w *windowManager) quit() {
	if w.cmd != nil {
		if err := w.cmd.Process.Kill(); err != nil {
			log.Println(w.cmd.Process.Kill()) // TODO:プロセスをKillするのは強引か?
		}
	}
}
