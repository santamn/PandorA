package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

const (
	// Domain and Protocol for PandA
	pandaDomain = "https://panda.ecs.kyoto-u.ac.jp"
	// URL for PandA log in page
	pandaLogin = pandaDomain + "/sakai-login-tool/container"
	// URL for all sites
	pandaAllSites = pandaDomain + "/direct/site.json"
	// URL for Resources
	pandaResources = pandaDomain + "/direct/content/site/" // {SITEID}.json を追記する
	// URL for Resource Acception
	pandaAcception = pandaDomain + "/access/accept?"
	// PATH for resource folder
	resourcePath = "~/Desktop/PandorA"
)

var (
	// URL for Kyoto University's CAS Login System
	casURL = "https://cas.ecs.kyoto-u.ac.jp/cas/login?service=" + url.QueryEscape(pandaLogin)
	// Infomation of downloaded resource
	downloaded = make(downloadMap, 0)
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
	lic, err := newLoggedInClient(ecsID, password)
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
func paraDownloadPDF(loggedInClient *http.Client, resources []resource) (errors []error) {
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
		if res.Type != "application/pdf" {
			continue
		}

		wg.Add(1)
		go func(lic *http.Client, info resource) {
			defer wg.Done()
			// 自動リダイレクトをOFFにする
			lic.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			// リソースをダウンロード
			resp, err := lic.Get(info.URL)
			if resp.Status == http.StatusText(302) {
				// 資料に著作権制限がついているために確認ページヘと飛ばされてしまう場合
				// 自動リダイレクトをONに戻す
				lic.CheckRedirect = nil

				// 資料のダウンロードの許可をくれるパスへクエリを投げ、そのままリダイレクト先で資料を取得
				path := "/content/group/" + info.lessonSite.ID + "/" + info.Title
				query := "ref=" + path + "&" + "url=" + path

				url := pandaAcception + url.QueryEscape(query)
				resp, err = lic.Get(url)
				resultChan <- result{response: resp, info: info, err: err}
			}
		}(loggedInClient, res)
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

		file, err := fetchFile(result.info.Title, result.info.lessonSite.Title)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if _, err := io.Copy(file, result.response.Body); err != nil { // ここを並列化するかどうかを考える
			errors = append(errors, err)
		}
	}
	return
}

// collectUnacquiredResouceInfo 未取得のリソースの情報を取得
func collectUnacquiredResouceInfo(loggedInClient *http.Client, sites []site) (resources []resource, err error) {
	// APIの返すJSONと形を合わせるための構造体
	type wrapper struct {
		Collection []resource `json:"content_collection"`
	}

	dmap := readDownloadMap()

	resources = make([]resource, 0, len(sites))
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

// collectSites 現在受講中の講義のSITEIDを収集する関数
func collectSites(loggedInClient *http.Client) (sites []site, err error) {
	// サイトの情報を取り出すための構造体
	type wrapper struct {
		Sites []site `json:"site_collection"`
	}

	sites = make([]site, 0)

	resp, err := loggedInClient.Get(pandaAllSites)
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

// newLoggedInClient ログイン済みのクライアントを返す関数
func newLoggedInClient(ecsID, password string) (client *http.Client, err error) {
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
