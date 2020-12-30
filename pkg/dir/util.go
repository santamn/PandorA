package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

// デスクトップまでのパスを取得する
// 取得できない場合は$HOMEを返す
func getPathToDesktop() (path string) {
	var home string
	switch runtime.GOOS {
	case "linux", "darwin":
		home = os.Getenv("HOME")
	case "windows":
		home = os.Getenv("USERPROFILE")
	}

	if err := os.Chdir(home); err != nil {
		return home
	}

	if info, err := os.Stat(home + "/Desktop"); !os.IsNotExist(err) && info.IsDir() {
		// $HOME/Desktopが存在する場合
		return home + "/Desktop"
	}

	if info, err := os.Stat(home + "/デスクトップ"); !os.IsNotExist(err) && info.IsDir() {
		// $HOME/デスクトップが存在する場合
		return home + "/デスクトップ"
	}

	return home
}

// PandorAフォルダへ移動する
func cdPandorA() error {
	pathToDesktop := getPathToDesktop()
	folderName := "PandorA Box"

	if err := os.Chdir(pathToDesktop); err != nil {
		return err
	}

	if info, err := os.Stat(folderName); os.IsNotExist(err) || !info.IsDir() {
		// PandorAフォルダが存在しない場合は作成する
		if err := os.Mkdir(folderName, 0766); err != nil {
			return err
		}
	}

	return os.Chdir(folderName)
}

// FetchFile PandorAフォルダ内のファイルを取得する関数 フォルダ名が空の場合はPandorAフォルダに直でファイルを作成・取得する
func FetchFile(filename, foldername string) (file *os.File, err error) {
	// あらかじめ戻り先を絶対パスに展開しておく
	prev, err := filepath.Abs(".")
	if err != nil {
		return file, err
	}
	defer os.Chdir(prev)

	if err := cdPandorA(); err != nil {
		return file, err
	}

	if foldername != "" {
		if info, err := os.Stat(foldername); os.IsNotExist(err) || !info.IsDir() {
			// 授業用のフォルダが存在しない場合は作成する
			if err := os.Mkdir(foldername, 0644); err != nil {
				return file, err
			}
		}

		if err := os.Chdir(foldername); err != nil {
			return file, err
		}
	}

	// ファイルがなければ作成し、読み書き両用でファイルを開く
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		// 同名のファイルが既に存在している場合にはファイル名に(n)をつけたファイルを作成
		text := strings.Split(filename, ".")

		changed := false
		for i := 1; i < 10; i++ {
			var newName string

			if len(text) == 2 {
				newName = text[0] + fmt.Sprintf("(%d)", i) + "." + text[1]
			} else {
				newName = fmt.Sprintf("(%d)", i) + filename
			}

			if _, err := os.Stat(newName); os.IsNotExist(err) {
				changed = true
				filename = newName
				break
			}
		}

		if !changed {
			// 10個以上同名のファイルが存在する場合にはuuidをファイル名の頭につける
			u, err := uuid.NewRandom()
			if err != nil {
				return file, err
			}
			filename = u.String() + filename
		}
	}

	file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	return
}

// FetchSettingsFile 設定ファイルを実行ファイルと同じディレクトリに生成する
func FetchSettingsFile(filename string) (file *os.File, err error) {
	// ファイルがなければ作成し、存在する場合は既に存在するファイルをオープンする
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
}
