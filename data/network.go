package data

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

const (
	// URL for Kyoto University's CAS system
	casURL = "https://cas.ecs.kyoto-u.ac.jp/cas/login?service="
	// Domain and Protocol for PandA
	pandaDomain = "https://panda.ecs.kyoto-u.ac.jp"
	// URL for PandA log in page
	pandaLogin = pandaDomain + "/sakai-login-tool/container"
	// URL for getting all assignments
	pandaAllAssignments = pandaDomain + "/direct/assignment/my.json"
	// URL for Resources
	pandaResources = "/direct/content/site"
)

var (
	// URL for cas login page
	casLogin = casURL + url.QueryEscape(pandaLogin)
)

// 現在受講中の講義のSITEIDを収集する関数
func collectSiteID(loggedInClient *http.Client) (siteIDs []string, err error) {
	// Assignment API からSITEIDを取り出すための構造体
	type (
		assignmentCollection struct {
			Context string `json:"context"`
		}

		myAssignment struct {
			Collection []assignmentCollection `json:"assignment_collection"`
		}
	)

	response, err := loggedInClient.Get(pandaAllAssignments)
	if err != nil {
		return siteIDs, err
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return siteIDs, err
	}

	var assignments myAssignment
	if err := json.Unmarshal(bytes, &assignments); err != nil {
		return siteIDs, err
	}

	// SITEIDの値を重複なくスライスに格納する
	m := make(map[string]struct{})
	siteIDs = make([]string, 0, len(assignments.Collection))
	for _, col := range assignments.Collection {
		if _, ok := m[col.Context]; !ok {
			m[col.Context] = struct{}{}
			siteIDs = append(siteIDs, col.Context)
		}
	}

	return siteIDs, nil
}

// LoggedInClient ログイン済みのクライアントを返す
func LoggedInClient(ecsID, password string) (client *http.Client, err error) {
	// set cookie jar
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client = &http.Client{Jar: jar}

	// pandaURLにGETを行うと、ログインページにリダイレクトされる
	// この際Pandaのドメインに対しJESESSIONIDが紐付けられる
	loginPage, err := client.Get(pandaLogin)
	if err != nil {
		return
	}
	defer loginPage.Body.Close()

	// get LT value
	lt, err := getLT(loginPage.Body)
	if err != nil {
		return
	}

	// login
	client, err = login(client, casLogin, lt, ecsID, password)
	if err != nil {
		return
	}

	return
}

// get Login Ticket(LT) from response body
func getLT(body io.Reader) (lt string, err error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	// get the tag in which lt value is written
	ltTag := doc.Find("input[name=\"lt\"]")
	// get lt value
	lt, exist := ltTag.Attr("value")
	if !exist {
		err = errors.New("LT is not found")
	}

	return
}

// PandAのログインフォームに情報を送信する関数
func login(client *http.Client, loginURL, lt, ecsID, password string) (loggedInClient *http.Client, err error) {
	values := url.Values{}
	// set form data
	values = map[string][]string{
		"_eventId":  {"submit"},
		"execution": {"e1s1"},
		"lt":        {lt},
		"username":  {ecsID},
		"password":  {password},
		"submit":    {"ログイン"},
	}

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(values.Encode()))
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
