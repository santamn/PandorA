package dir

import (
	"os"
	"path/filepath"
)

// PandorAフォルダへ移動する
func cdPandorA() error {
	if err := os.Chdir(os.Getenv("HOME") + "/Desktop"); err != nil {
		return err
	}

	if info, err := os.Stat("PandorA"); os.IsNotExist(err) || !info.IsDir() {
		// PandorAフォルダが存在しない場合は作成する
		if err := os.Mkdir("PandorA", 0766); err != nil {
			return err
		}
	}

	return os.Chdir("PandorA")
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
			if err := os.Mkdir(foldername, 0766); err != nil {
				return file, err
			}
		}

		if err := os.Chdir(foldername); err != nil {
			return file, err
		}
	}

	// ファイルがなければ作成し、読み書き両用でファイルを開く
	file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0766)
	return
}
