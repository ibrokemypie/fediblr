package fedi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/ibrokemypie/fediblr/config"
)

type Status struct {
	Images        []string
	Caption       string
	SourceName    string
	SourceURL     string
	RebloggedName string
	RebloggedURL  string
}

func newfileUploadRequest(uri string, params map[string]string, paramName string, file io.ReadCloser, fileName string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func uploadImage(configuration config.Config, imageURL string) string {
	resp, err := http.Get(imageURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	req, err := newfileUploadRequest(configuration.FediInstance+"/api/v1/media", nil, "file", resp.Body, path.Base(imageURL))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "bearer "+configuration.FediToken)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(resp)
	}

	var result map[string]interface{}
	readBody, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(readBody, &result)
	if result["type"] == "unknown" {
		fmt.Println(string(readBody))
		panic(result)
	}

	return result["id"].(string)
}

func createStatus(configuration config.Config, mediaIDs []string, status Status) {
	statusText := status.Caption +
		"\n\nSource: " + status.SourceName + " " + status.SourceURL +
		"\nReblogged From: " + status.RebloggedName + " " + status.RebloggedURL
	params := url.Values{}
	for _, id := range mediaIDs {
		params.Add("media_ids[]", id)
	}
	params.Add("status", statusText)
	params.Add("visibility", configuration.Visibility)
	postData := strings.NewReader(params.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", configuration.FediInstance+"/api/v1/statuses", postData)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "bearer "+configuration.FediToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func PostStatus(status Status, configuration config.Config) {
	mediaIDs := []string{}
	for _, image := range status.Images {
		mediaIDs = append(mediaIDs, uploadImage(configuration, image))
	}
	createStatus(configuration, mediaIDs, status)
}
