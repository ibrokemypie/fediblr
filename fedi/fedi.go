package fedi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Status struct {
	ImageURL      string
	Caption       string
	SourceName    string
	SourceURL     string
	RebloggedName string
	RebloggedURL  string
}

type App struct {
	ClientName   string `json:"name"`
	Website      string `json:"website"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
type Oauth struct {
	AccessToken string `json:"access_token"`
}

func getApp(instanceURL string) App {
	clientName := "fediblr"
	website := "https://github.com/ibrokemypie/fediblr"
	redirectURI := "urn:ietf:wg:oauth:2.0:oob"

	resp, err := http.PostForm(instanceURL+"/api/v1/apps",
		url.Values{
			"client_name":   {clientName},
			"scopes":        {"read write"},
			"website":       {website},
			"redirect_uris": {redirectURI}})
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var result App

	json.NewDecoder(resp.Body).Decode(&result)

	return result
}

func authorizeUser(instanceURL string, app App) string {
	u, err := url.Parse(instanceURL + "/oauth/authorize")
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Set("client_id", app.ClientID)
	q.Set("redirect_uri", app.RedirectURI)
	q.Set("response_type", "code")
	q.Set("force_login", "true")
	q.Set("scope", "read write")
	u.RawQuery = q.Encode()

	fmt.Println("Please open the following URL in your browser.")
	fmt.Println("Once you have authenticated, paste the token from the page into this window.")
	fmt.Println(u)

	var token string
	fmt.Scanln(&token)

	return token
}

func oauthToken(instanceURL string, app App, token string) string {
	resp, err := http.PostForm(instanceURL+"/oauth/token",
		url.Values{
			"client_id":     {app.ClientID},
			"client_secret": {app.ClientSecret},
			"redirect_uri":  {app.RedirectURI},
			"scope":         {"read write"},
			"grant_type":    {"authorization_code"},
			"code":          {token}})
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(body, &result)
	if result["access_token"] == nil {
		fmt.Println(string(body))
		panic(result)
	}

	return result["access_token"].(string)
}

func authorize(instanceURL string) {
	app := getApp(instanceURL)
	token := authorizeUser(instanceURL, app)
	accessToken := oauthToken(instanceURL, app, token)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", instanceURL+"/api/v1/apps/verify_credentials", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, _ := client.Do(req)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

func PostStatus(status Status, instanceURL string) {
	authorize(instanceURL)

	fmt.Println(status.ImageURL)
	fmt.Println(status.Caption)
	fmt.Println(status.SourceName + " " + status.SourceURL)
	fmt.Println(status.RebloggedName + " " + status.RebloggedURL)
}

// img
// text
// source
// reblogged from
