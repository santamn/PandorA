package data

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

const (
	// Domain and Protocol for PandA
	pandaDomain = "https://panda.ecs.kyoto-u.ac.jp"
	// URL for PandA log in page
	pandaLogin = pandaDomain + "/sakai-login-tool/container"
	// URL for getting all assignments
	pandaAllAssignments = pandaDomain + "/direct/assignment/my.json"
	// URL for Resources
	pandaResources = pandaDomain + "/direct/content/site/" // SITEID.json を追記する
	// URL for all sites
	pandaAllSites = pandaDomain + "/direct/site.json"
	// PATH for site folder
	path = "../resource"
)

var (
	// URL for Kyoto University's CAS Login System
	casURL = "https://cas.ecs.kyoto-u.ac.jp/cas/login?service=" + url.QueryEscape(pandaLogin)
)

// Site PandAのサイト情報を取得するための構造体
type Site struct {
	Title       string `json:"title"`
	CreatedDate int64  `json:"createdDate"`
	ID          string `json:"id"`
}

// Resource リソースの情報を表す構造体
type Resource struct {
	Size         int64  `json:"size"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	LessonName   string
	LastModified int64
}

// DownloadMap すでにダウンロードした資料についての情報を表すマップ
// "SiteID1":{
//     "資料名1":"最終修正時刻1",
//     "資料名2":"最終修正時刻2",
//     ...
// },
// "SiteID2:{
//     "資料名1":"最終修正時刻1",
//     "資料名2":"最終修正時刻2",
//     ...
// },
// という構造になっており、最終修正時刻が最後にダンロードした時から変化したものか、ここに登録されていないリソースのみダウンロードする
type DownloadMap map[string]map[string]string

// DownloadPDF 科目のサイトに登録されたPDFをすべてダウンロードする関数 -> 全てのリソースのURLを取得して、それらを並列にダウンロードする関数に変更
func DownloadPDF(loggedInClient *http.Client, siteID string) error {
	type (
		// APIの返すJSONと形を合わせるための構造体
		wrapper struct {
			Collection []Resource `json:"content_collection"`
		}

		// HTTPレスポンスとエラーをどちらも呼び出し側で扱うための構造体
		result struct {
			response *http.Response
			info     Resource
			err      error
		}
	)

	url := pandaResources + siteID + ".json"
	resp, err := loggedInClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var w wrapper
	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return err
	}

	resultChan := make(chan result, len(w.Collection))
	for _, col := range w.Collection {
		if col.Type != "application/pdf" {
			continue
		}

		go func(loggedInClient *http.Client, info resouceInfo) {
			defer close(resultChan)

			r, e := loggedInClient.Get(info.URL)
			resultChan <- result{response: r, info: info, err: e}
		}(loggedInClient, col)
	}

	for r := range resultChan {
		defer r.response.Body.Close()

		// 以下、ファイル作成できるモノについてはファイルを作成し、
		// 失敗したものについてはまとめてエラーを返す

		if r.err != nil {
			return r.err
		}

		file, err := os.Create(r.info.Title)
		if err != nil {
			return err // ここのエラーをほんとに止めるべきかは考える
		}

		written, err := io.Copy(file, r.response.Body)
		if err != nil {
			return err
		}
		if written != r.info.Size {
			return errors.New("Written bytes are not equal to file size")
		}
	}
	close(resultChan)

	return nil
}

// CollectResouceURLs 有効なリソースURLを取得するための関数
func CollectResouceURLs(loggedInClient *http.Client, sites []Site) (resources []Resource, err error) {
	type (
		// APIからリソースの情報を取得するための構造体
		resourceInfo struct {
			ModifiedDate string `json:"modifiedDate"`
			Size         int64  `json:"size"`
			Type         string `json:"type"`
			Title        string `json:"title"`
			URL          string `json:"url"`
		}

		// APIの返すJSONと形を合わせるための構造体
		wrapper struct {
			Collection []resourceInfo `json:"content_collection"`
		}
	)

	resources = make([]Resource, 0, len(sites))

	for _, site := range sites {
		url := pandaResources + site.ID + ".json"
		resp, err := loggedInClient.Get(url)
		if err != nil {
			return resources, err
		}
		defer resp.Body.Close()

		var w wrapper
		if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
			return resources, err
		}

		for _, col := range w.Collection {

		}

	}
}

// CollectSites 現在受講中の講義のSITEIDを収集する関数
func CollectSites(loggedInClient *http.Client) (sites []Site, err error) {
	// サイトの情報を取り出すための構造体
	type wrapper struct {
		Sites []Site `json:"site_collection"`
	}

	sites = make([]Site, 0)

	resp, err := loggedInClient.Get(pandaAllSites)
	if err != nil {
		return sites, err
	}
	defer resp.Body.Close()

	var w wrapper
	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return sites, err
	}

	start, end := getSemesterBound()
	for _, s := range w.Sites {
		// APIから取得した値はUnixミリ秒なので、Unix秒に変換する
		t := s.CreatedDate / 1000
		if start <= t && t <= end {
			// 今学期に作成された科目の情報を取得する
			sites = append(sites, s)
		}
	}

	return sites, nil
}

// NewLoggedInClient ログイン済みのクライアントを返す関数
func NewLoggedInClient(ecsID, password string) (client *http.Client, err error) {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client = &http.Client{Jar: jar}

	// pandaURLにGETを行うと、ログインページにリダイレクトされる
	// この際Pandaのドメインに対しJESESSIONIDが紐付けられる
	loginPage, err := client.Get(pandaLogin)
	if err != nil {
		return client, err
	}
	defer loginPage.Body.Close()

	// ログインページからLT(おそらくログインチケットの略)を取得
	lt, err := getLT(loginPage.Body)
	if err != nil {
		return client, err
	}

	// ログイン
	client, err = login(client, lt, ecsID, password)
	if err != nil {
		return
	}

	return
}

// LTをログインページから取り出す関数
func getLT(body io.Reader) (lt string, err error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	// LTが書かれたタグを取得する
	ltTag := doc.Find("input[name=\"lt\"]")
	// LTを取得
	lt, exist := ltTag.Attr("value")
	if !exist {
		err = errors.New("LT is not found")
	}

	return
}

// 京大のCASシステムにログイン情報をPOSTする関数
func login(client *http.Client, lt, ecsID, password string) (loggedInClient *http.Client, err error) {
	values := url.Values{
		"_eventId":  {"submit"},
		"execution": {"e1s1"},
		"lt":        {lt},
		"username":  {ecsID},
		"password":  {password},
		"submit":    {"ログイン"},
	}

	req, err := http.NewRequest("POST", casURL, strings.NewReader(values.Encode()))
	if err != nil {
		return client, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// ログインフォームに必要なデータを送信すると、PandAのポータルサイトにリダイレクトする
	// この際、クエリパラメータとして発行されるticketを用いて、JSESSIONIDを認証済みにする処理がサーバー側で行われる
	if _, err = client.Do(req); err != nil {
		return client, err
	}

	return client, nil
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

// 並列にPandaからリソースをリクエストする関数
func resourceRequest(loggerInClient *http.Client, url string) {

}

// 並列処理でRequest.Bodyから読み出してファイルに書き込む関数
