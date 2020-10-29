package resource

import (
	"encoding/json"
	"pandora/pkg/dir"
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

	mapFile, err := dir.FetchFile("dmap.dat", "")
	if err != nil {
		return dmap
	}

	json.NewDecoder(mapFile).Decode(&dmap)

	return dmap
}

// writeFile ダウンロードマップをファイルに書き込む
func (dmap downloadMap) writeFile() error {
	mapFile, err := dir.FetchFile("dmap.dat", "")
	if err != nil {
		return err
	}

	e := json.NewEncoder(mapFile)
	e.SetIndent("", "  ")

	return e.Encode(dmap)
}
