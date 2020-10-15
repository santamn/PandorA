package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

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
	pandaResources = pandaDomain + "/direct/content/site/"
)

var (
	// URL for Kyoto University's CAS Login System
	casURL = "https://cas.ecs.kyoto-u.ac.jp/cas/login?service=" + url.QueryEscape(pandaLogin)
)

// DownloadPDF 科目のサイトに登録されたPDFをすべてダウンロードする関数
func DownloadPDF(loggedInClient *http.Client, siteID string) error {
	// リソース情報を取得するための構造体
	type (
		contentCollection struct {
			Size  int64  `json:"size"`
			Type  string `json:"type"`
			Title string `json:"title"`
			URL   string `json:"url"`
		}

		resources struct {
			Collection []contentCollection `json:"content_collection"`
		}
	)

	url := pandaResources + siteID + ".json"
	resp, err := loggedInClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r resources
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	for _, col := range r.Collection {
		if col.Type != "application/pdf" {
			continue
		}

		resp, err := loggedInClient.Get(col.URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Println("Title", col.Title)
		file, err := os.Create(col.Title)
		if err != nil {
			return err
		}

		written, err := io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
		if written != col.Size {
			return errors.New("Invalid bytes")
		}
	}

	return nil
}

// CollectSiteID 現在受講中の講義のSITEIDを収集する関数
func CollectSiteID(loggedInClient *http.Client) (siteIDs []string, err error) {
	// Assignment API からSITEIDを取り出すための構造体
	type (
		assignmentCollection struct {
			Context string `json:"context"`
		}

		myAssignments struct {
			Collection []assignmentCollection `json:"assignment_collection"`
		}
	)

	resp, err := loggedInClient.Get(pandaAllAssignments)
	if err != nil {
		return siteIDs, err
	}
	defer resp.Body.Close()

	var a myAssignments
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return siteIDs, err
	}

	// SITEIDの値を重複なくスライスに格納する
	m := make(map[string]struct{})
	siteIDs = make([]string, 0, len(a.Collection))
	for _, col := range a.Collection {
		if _, ok := m[col.Context]; !ok {
			m[col.Context] = struct{}{}
			siteIDs = append(siteIDs, col.Context)
		}
	}

	return siteIDs, nil
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
	values := url.Values{}
	values = map[string][]string{
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
