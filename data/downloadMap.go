package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// downloadMap すでにダウンロードした資料についての情報を表すマップ
//
// 	"SiteID1":{
// 	　"資料名1":"最終修正時刻1",
//   　　"資料名2":"最終修正時刻2",
//    ...
// 	},
//
// 	"SiteID2":{
// 		"資料名3":"最終修正時刻3",
// 		"資料名4":"最終修正時刻4",
// 		...
//	},
//
// という構造になっており、最終修正時刻が最後にダンロードした時から変化したものか、ここに登録されていないリソースのみダウンロードする
type downloadMap map[string]map[string]string

// readDownloadMap ダウンロードマップをファイルから読み出す
func readDownloadMap() downloadMap {
	dmap := make(downloadMap)
	// あらかじめ戻り先を絶対パスに展開しておく
	prev, err := filepath.Abs(".")
	if err != nil {
		return dmap
	}
	defer os.Chdir(prev)

	if err := cdPandorA(); err != nil {
		return dmap
	}

	mapFile, err := os.Open("dmap")
	if err != nil {
		// ファイルが開けない場合は空のマップを返す
		return dmap
	}

	json.NewDecoder(mapFile).Decode(&dmap)

	return dmap
}

// writeFile ダウンロードマップをファイルに書き込む
func (dmap downloadMap) writeFile() error {
	// あらかじめ戻り先を絶対パスに展開しておく
	prev, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	defer os.Chdir(prev)

	if err := cdPandorA(); err != nil {
		return err
	}

	mapFile, err := os.OpenFile("dmap", os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}

	e := json.NewEncoder(mapFile)
	e.SetIndent("", "  ")

	return e.Encode(dmap)
}
