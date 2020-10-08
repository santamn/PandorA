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
	"time"

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
)

var (
	// URL for cas login page
	casLogin = casURL + url.QueryEscape(pandaLogin)
)

// Assignment 課題
type Assignment struct {
	AssignmentID   string
	AssignmentName string
	CloseTime      time.Time
	DueTime        time.Time
	Instructions   string
	LessonName     string
	Status         uint8
}

// Keitai APIから取得したJSONを平坦化した構造体
// ここに科目名と提出状況が欲しい
// あとtime型へ変形したい
type flattenAssignment struct {
	assignmentID   string
	assignmentName string
	lessonID       string
	instructions   string
	dueTime        int64
	closeTime      int64
}

// 課題一覧取得APIから科目ID・課題ID・課題名・課題内容・締め切りを取得する
// 科目IDについては別途一覧を返す
func fetchAssingmentInfo(client *http.Client) (flatten []flattenAssignment, lessonIDs []string, err error) {
	// JSONをパースする用の構造体
	type (
		Assignment struct {
			Close struct {
				Time int64 `json:"time"`
			} `json:"closeTime"`
			Due struct {
				Time int64 `json:"time"`
			} `json:"dueTime"`
			AssignmentID   string `json:"id"`
			AssignmentName string `json:"title"`
			Instructions   string `json:"instructions"`
			LessonID       string `json:"context"`
		}

		Reciever struct {
			Coll []Assignment `json:"assignment_collection"`
		}
	)

	resp, err := client.Get(pandaAllAssignments)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	bytesBody, _ := ioutil.ReadAll(resp.Body)
	var reciever Reciever
	if err := json.Unmarshal(bytesBody, &reciever); err != nil {
		return nil, nil, err
	}

	flatten = make([]flattenAssignment, len(reciever.Coll))
	for i, item := range reciever.Coll {
		flatten[i].assignmentID = item.AssignmentID
		flatten[i].assignmentName = item.AssignmentName
		flatten[i].closeTime = item.Close.Time
		flatten[i].dueTime = item.Due.Time
		flatten[i].instructions = item.Instructions
		flatten[i].lessonID = item.LessonID
	}

	// 重複なく科目IDをスライスにまとめる
	lessonIDs = make([]string, 0, len(flatten))
	m := make(map[string]bool)
	for _, f := range flatten {
		if !m[f.lessonID] {
			m[f.lessonID] = true
			lessonIDs = append(lessonIDs, f.lessonID)
		}
	}

	return
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

// send data to login form
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
