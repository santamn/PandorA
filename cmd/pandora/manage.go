package main

import (
	"fmt"
	"log"
	"os/exec"
	"pandora/pkg/account"
	"pandora/pkg/dir"
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
	wg               sync.WaitGroup
}

func (d *downloadManager) excute(window *windowManager, clicked bool) {
	defer func() {
		d.isRunning = false
	}()

	d.mu.Lock()
	if !d.isRunning {
		d.isRunning = true
		d.wg.Add(1)
		d.mu.Unlock()

		if min := time.Now().Sub(d.lastExecutedTime).Minutes(); min < 10 && clicked {
			// 前のダウンロードからの経過時間が10分以内にユーザーによる再度の実行の要求があれば警告を出して終了する
			alert(fmt.Sprintf("PandorA needs cool time. Please try after %d minute(s) later at least.", uint(10-min)))
			return
		}

		ecsID, password, rejectable, err := account.ReadAccountInfo()
		if err != nil {
			log.Println("read account error 1:", err)
			// アカウント情報を入力させる
			window.show()
			window.wg.Wait() // アカウント情報の入力を待つ
		}
		ecsID, password, rejectable, err = account.ReadAccountInfo()
		if err != nil {
			log.Println("read account error 2:", err)
			// 2回目にエラーが出た場合はエラーを表示して終了する
			alert(err.Error())
			return
		}

		notify("NOW DOWNLOADING")

		d.lastExecutedTime = time.Now()
		if errors := resource.Download(ecsID, password, rejectable); len(errors) > 0 {
			for _, err := range errors {
				log.Println("Download error:", err)

				switch err.(type) {
				case *pandaapi.NetworkError:
					alert("Network Error: something wrong with connecting the Internet")
				case *pandaapi.DeadPandAError:
					alert(err.Error())
				case *pandaapi.FailedLoginError:
					alert(err.Error())
					go window.show()
				default:
					alert("System Error: " + err.Error())
				}
			}
		} else {
			notify("Download succeeded!")
		}
		// wg.Waitを使えばここでダウンロードが終了することを待つことができる
		d.wg.Done()
	} else {
		d.mu.Unlock()
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

func newWindowManager() *windowManager {
	return &windowManager{
		path: filepath.Join(dir.WorkingDirecory, "form"),
	}
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

// ウィンドウを終了する
func (w *windowManager) quit() {
	if w.cmd != nil {
		w.cmd.Process.Kill()
	}
}
