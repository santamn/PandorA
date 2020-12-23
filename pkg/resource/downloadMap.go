package resource

import (
	"encoding/json"
	"pandora/pkg/dir"
)

// downloadMap すでにダウンロードした資料についての情報を表すマップ
//
// 	"SiteID1":{
//		"resource1":"last-modified1",
//		"resource2":"last-modified2",
//		...
// 	},
// 	"SiteID2":{
// 		"resource3":"last-modified3",
// 		"resource4":"last-modified4",
// 		...
//	},
//
// という構造になっており、最終修正時刻が最後にダンロードした時から変化したものか、ここに登録されていないリソースのみダウンロードする
//
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

// writeToFile ダウンロードマップをファイルに書き込む
func (dmap downloadMap) writeToFile() error {
	mapFile, err := dir.FetchFile("dmap.dat", "")
	if err != nil {
		return err
	}

	e := json.NewEncoder(mapFile)
	e.SetIndent("", "  ")

	return e.Encode(dmap)
}
