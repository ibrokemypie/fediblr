package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetFediInstance() string {
	fmt.Println("Paste the fedi instance url (eg: https://mastodon.social) and press enter")

	var instance string
	fmt.Scanln(&instance)

	return instance
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

func GetFediAccessToken(instanceURL string) string {
	app := getApp(instanceURL)
	token := authorizeUser(instanceURL, app)
	accessToken := oauthToken(instanceURL, app, token)

	return accessToken
}
