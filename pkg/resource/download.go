package resource

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	dir "pandora/pkg/dir"
	pandaapi "pandora/pkg/pandaAPI"
	"strings"
	"sync"
	"time"
)

// site PandAのサイト情報を取得するための構造体
type site struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}

// resource リソースの情報を表す構造体
type resource struct {
	Size         int64  `json:"size"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	LastModified string `json:"modifiedDate"`
	lessonSite   site
}

// Download 資料をダウンロード
func Download(ecsID, password string) {
	lic, err := pandaapi.NewLoggedInClient(ecsID, password)
	if err != nil {
		fmt.Println(err)
	}

	sites, err := collectSites(lic)
	if err != nil {
		fmt.Println(err)
	}

	resources, err := collectUnacquiredResouceInfo(lic, sites)
	if err != nil {
		fmt.Println(err)
	}

	if errors := paraDownloadPDF(lic, resources); len(errors) > 0 {
		fmt.Println(errors)
	}
}

// paraDownloadPDF 未取得のリソースを並列にダウンロードする関数
func paraDownloadPDF(lic *pandaapi.LoggedInClient, resources []resource) (errors []error) {
	// HTTPレスポンスとエラーをどちらも呼び出し側で扱うための構造体
	type result struct {
		response *http.Response
		info     resource
		err      error
	}

	errors = make([]error, 0)

	var wg sync.WaitGroup
	resultChan := make(chan result, len(resources))

	for _, res := range resources {
		if res.Type != "application/pdf" { // TODO:音声・動画ファイルの形式を調べてそれのみを弾くようにする
			continue
		}

		wg.Add(1)
		go func(lic *pandaapi.LoggedInClient, info resource) {
			defer wg.Done()

			// リソースをダウンロード
			resp, err := lic.FetchResource(info.lessonSite.ID, info.Title)
			resultChan <- result{response: resp, info: info, err: err}

		}(lic, res)
	}

	// 送信するものがなくなったらチャネルをクローズする
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		defer result.response.Body.Close()

		if result.err != nil {
			errors = append(errors, result.err)
			continue
		}

		file, err := dir.FetchFile(result.info.Title, result.info.lessonSite.Title)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if _, err := io.Copy(file, result.response.Body); err != nil { // TODO:ここを並列化するかどうかを考える
			errors = append(errors, err)
		}
	}
	return
}

// collectUnacquiredResouceInfo 未取得のリソースの情報を取得
func collectUnacquiredResouceInfo(lic *pandaapi.LoggedInClient, sites []site) (resources []resource, err error) {
	// APIの返すJSONと形を合わせるための構造体
	type wrapper struct {
		Collection []resource `json:"content_collection"`
	}

	dmap := readDownloadMap()

	resources = make([]resource, 0, len(sites))
	for _, site := range sites {
		// サイトのリソース情報を取得
		resp, err := lic.FetchSiteResources(site.ID)
		if err != nil {
			return resources, err
		}
		defer resp.Body.Close()

		var w wrapper
		if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
			return resources, err
		}

		for _, res := range w.Collection {
			// リソース情報に講義名を追加
			res.lessonSite = site

			// ダウンロードしていない資料もしくは最終編集時刻が変更されているもののみダウンロード候補へ追加する
			resourceMap, ok := dmap[site.ID]
			if !ok { // サイトIDがダウンロードマップに存在しない場合(= その講義にはじめて資料が追加された)
				resources = append(resources, res)

				// ダウンロードマップを更新
				dmap[site.ID] = map[string]string{res.Title: res.LastModified}
				continue
			}

			if lastModified, ok := resourceMap[res.Title]; !ok || lastModified != res.LastModified {
				// 資料名がダウンロードマップに登録されていない場合(= いままでにダウンロードされたことがない)
				// もしくは最終編集時刻が過去のものと異なっている場合

				resources = append(resources, res)
				// ダウンロードマップを更新
				dmap[site.ID][res.Title] = res.LastModified
				continue
			}
		}
	}

	if err := dmap.writeFile(); err != nil {
		return resources, err
	}

	return
}

// collectSites 現在受講中の講義の授業サイトに関する情報を収集
func collectSites(lic *pandaapi.LoggedInClient) (sites []site, err error) {
	// サイトの情報を取り出すための構造体
	type wrapper struct {
		Sites []site `json:"site_collection"`
	}

	sites = make([]site, 0)

	resp, err := lic.FetchAllSites()
	if err != nil {
		return sites, err
	}
	defer resp.Body.Close()

	var w wrapper
	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return sites, err
	}

	semesterText := makeSemesterDescription()

	for _, s := range w.Sites {
		if strings.Contains(s.Title, semesterText) {
			// 科目名に含まれる"2020前期"の部分で科目が現在受講中かどうかを判定する
			sites = append(sites, s)
		}
	}

	return sites, nil
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
