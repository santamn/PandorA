package data

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ASCIIコードで33(!)-126(~)をrotする
func rot47(target []byte) (result []byte) {
	result = make([]byte, len(target))

	for i, t := range target {
		result[i] = (t-33+47)%94 + 33
	}

	return
}

// 科目名に含まれる "2020前期" の部分を作成する
func makeSemesterDescription() (text string) {
	year, month, _ := time.Now().Date()

	switch {
	case 3 <= month && month <= 8:
		// 前期
		text = fmt.Sprint(year) + "前期"
	default:
		// 後期
		text = fmt.Sprint(year) + "後期"
	}

	return
}

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

// PandorAフォルダにファイルを作成する関数
func fetchFile(filename, foldername string) (file *os.File, err error) {
	// あらかじめ戻り先を絶対パスに展開しておく
	prev, err := filepath.Abs(".")
	if err != nil {
		return file, err
	}
	defer os.Chdir(prev)

	if err := cdPandorA(); err != nil {
		return file, err
	}

	if info, err := os.Stat(foldername); os.IsNotExist(err) || !info.IsDir() {
		// 授業用のフォルダが存在しない場合は作成する
		if err := os.Mkdir(foldername, 0766); err != nil {
			return file, err
		}
	}

	if err := os.Chdir(foldername); err != nil {
		return file, err
	}

	file, err = os.Create(filename)
	return
}
