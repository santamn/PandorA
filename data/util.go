package data

import (
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

// 今学期の開始時刻と終了時刻をUnixタイムを返す関数
func getSemesterBound() (start, end int64) {
	year, month, _ := time.Now().Date()
	loc, _ := time.LoadLocation("Asia/Tokyo")

	switch {
	case 3 <= month && month <= 8:
		// 前期
		start = time.Date(year, 3, 1, 0, 0, 0, 0, loc).Unix()
		end = time.Date(year, 8, 31, 24, 0, 0, 0, loc).Unix()
	case (9 <= month && month <= 12) || (1 <= month && month <= 2):
		// 後期
		start = time.Date(year, 9, 1, 0, 0, 0, 0, loc).Unix()
		end = time.Date(year+1, 3, 0, 0, 0, 0, 0, loc).Unix()
	}

	return
}

// デスクトップにファイルを作成する関数
func fetchFile(filename, foldername string) (file *os.File, err error) {
	// あらかじめ戻り先を絶対パスに展開しておく
	prev, err := filepath.Abs(".")
	if err != nil {
		return file, err
	}
	defer os.Chdir(prev)

	if err := os.Chdir(os.Getenv("HOME") + "/Desktop"); err != nil {
		return file, err
	}

	if info, err := os.Stat("PandorA"); os.IsNotExist(err) || !info.IsDir() {
		// PandorAフォルダが存在しない場合は作成する
		if err := os.Mkdir("PandorA", 0766); err != nil {
			return file, err
		}
	}

	if err := os.Chdir("PandorA"); err != nil {
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
