package pandaapi

import (
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
	// Domain and Protocol for PandA
	pandaDomain = "https://panda.ecs.kyoto-u.ac.jp"
	// URL for PandA log in page
	pandaLogin = pandaDomain + "/sakai-login-tool/container"
	// URL for all sites
	pandaAllSites = pandaDomain + "/direct/site.json"
	// URL for Resources Infomation
	pandaResourcesInfo = pandaDomain + "/direct/content/site/" // {SITEID}.json を追記する
	// URL for getting resource
	pandaResource = pandaDomain + "/access" // {SITEID}/{フォルダ名(あれば)}/{資料名} を追記する
	// URL for Resource Acception
	pandaAcception = pandaDomain + "/access/accept?"
	// Error message appears when failed to log in
	loginErrorMessage = "あなたが入力した認証情報は，認証可能なものであることが確認できませんでした．"
)

var (
	// URL for Kyoto University's CAS Login System
	casURL = "https://cas.ecs.kyoto-u.ac.jp/cas/login?service=" + url.QueryEscape(pandaLogin)
)

// LoggedInClient PandAにログイン済みのクライアントを表す
type LoggedInClient struct {
	c *http.Client
}

// CheckPandaStatus PandAサーバが生きているかどうかを判定する
func CheckPandaStatus() error {
	// リダイレクトを無効にする
	c := http.DefaultClient
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := c.Head(pandaDomain)
	if err != nil {
		return &NetworkError{err: err}
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return &DeadPandAError{code: resp.StatusCode, err: nil, url: pandaDomain}
	}

	return nil
}

// FetchAllSites 全ての授業サイトの情報を取得するAPI レスポンスボディをクローズする必要がある
func (lic *LoggedInClient) FetchAllSites() (resp *http.Response, err error) {
	resp, err = lic.c.Get(pandaAllSites)
	// 200以外のレスポンスが帰ってくる場合はサーバーが死んでいるとみなす
	if resp.StatusCode != 200 {
		return resp, &DeadPandAError{code: resp.StatusCode, err: err, url: pandaAllSites}
	}

	return
}

// FetchSiteResources 授業サイトに登録されているリソースの情報を取得するAPI レスポンスボディをクローズする必要がある
func (lic *LoggedInClient) FetchSiteResources(siteID string) (resp *http.Response, err error) {
	siteURL := pandaResourcesInfo + siteID + ".json"

	resp, err = lic.c.Get(siteURL)
	// 200以外のレスポンスが帰ってくる場合はサーバーが死んでいるとみなす
	if resp.StatusCode != 200 {
		return resp, &DeadPandAError{code: resp.StatusCode, err: err, url: siteURL}
	}

	return
}

// FetchResource リソースを取得するAPI レスポンスボディをクローズする必要がある
func (lic *LoggedInClient) FetchResource(uri string) (resp *http.Response, err error) {
	resp, err = lic.c.Get(uri)
	if err != nil {
		return resp, &NetworkError{err: err}
	}

	// 通常のダウンロードに成功した場合
	if resp.StatusCode == 200 {
		return
	}

	// 著作権制限付きダウンロード警告がでる場合
	if resp.StatusCode == 302 {
		// /{SITEID}/{フォルダパス}/{資料名}の部分を取得
		path := strings.Replace(uri, pandaResource, "", 1)
		// 資料のダウンロードの許可をくれるパスへクエリを投げる
		query := "ref=" + path + "&url=" + path
		r, e := lic.c.Get(pandaAcception + query)
		if e != nil {
			return resp, &NetworkError{err: e}
		}
		defer func() {
			io.Copy(ioutil.Discard, r.Body)
			r.Body.Close()
		}()

		resp, err = lic.c.Get(uri)
		if err != nil {
			return resp, &NetworkError{err: err}
		}
		if resp.StatusCode != 200 {
			return resp, &DeadPandAError{code: resp.StatusCode, err: err, url: uri}
		}
		return
	}

	// 200と302以外のレスポンスを返す場合はサーバーが死んでいるとみなす
	return resp, &DeadPandAError{code: resp.StatusCode, err: err, url: uri}
}

// NewLoggedInClient ログイン済みのクライアントを返す関数
func NewLoggedInClient(ecsID, password string) (lic *LoggedInClient, err error) {
	// Cookieを保存する
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client := &http.Client{Jar: jar}

	// まずPandAの生存確認を行う
	// この関数内ではここで生存が確認された場合にはログイン中はPandAが死んでいないものと推定する
	if err := CheckPandaStatus(); err != nil {
		return &LoggedInClient{c: client}, err
	}

	// pandaURLにGETを行うと、ログインページにリダイレクトされる
	// この際Pandaのドメインに対しJESESSIONIDが紐付けられる
	loginPage, err := client.Get(pandaLogin)
	if err != nil {
		return &LoggedInClient{c: client}, &NetworkError{err: err}
	}
	defer loginPage.Body.Close()

	// ログインページからLT(おそらくログインチケットの略)を取得
	lt, err := getLT(loginPage)
	if err != nil {
		return &LoggedInClient{c: client}, err
	}

	// ログイン
	client, err = login(client, lt, ecsID, password)
	if err != nil {
		return &LoggedInClient{c: client}, err
	}

	//　リダイレクトを無効にする
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &LoggedInClient{c: client}, nil
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
	resp, err := client.Do(req)
	if err != nil {
		return client, &NetworkError{err: err}
	}
	defer resp.Body.Close()

	// ログインに成功したかどうかを確認する
	auth, err := isAuthorized(resp)
	if err != nil {
		return client, err
	}
	if !auth {
		return client, &FailedLoginError{EscID: ecsID, Password: password}
	}

	return client, nil
}

// LTをログインページから取り出す関数
func getLT(resp *http.Response) (lt string, err error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
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

// ログインに成功していかどうかを判定する関数
func isAuthorized(resp *http.Response) (bool, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return false, err
	}

	// ログインに失敗したときのメッセージが書かれたタグを取得する
	msg := doc.Find("#msg")
	if msg.Length() == 0 {
		// ログインに成功した場合
		return true, nil
	}

	if msg.Text() == loginErrorMessage {
		// ログインに失敗した場合
		return false, nil
	}

	// 想定外の動作
	return false, errors.New("There's something wrong with the login system")
}
