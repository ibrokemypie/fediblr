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
	"os"
	"path"
	"strings"

	"github.com/ibrokemypie/fediblr/config"
)

type Status struct {
	ImageURL      string
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

func uploadImage(configuration map[string]string, imageURL string) string {
	if configuration["lastImage"] != path.Base(imageURL) {
		configuration["lastImage"] = path.Base(imageURL)
		config.WriteConfig(configuration)
	} else {
		fmt.Println("Already posted this image before, skipping.")
		os.Exit(1)
	}

	resp, err := http.Get(imageURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	req, err := newfileUploadRequest(configuration["fediInstance"]+"/api/v1/media", nil, "file", resp.Body, path.Base(imageURL))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "bearer "+configuration["fediToken"])

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

func createStatus(configuration map[string]string, mediaID string, status Status) {
	statusText := status.Caption +
		"\n\nSource: " + status.SourceName + " " + status.SourceURL +
		"\nReblogged From: " + status.RebloggedName + " " + status.RebloggedURL
	params := url.Values{}
	params.Add("media_ids[]", mediaID)
	params.Add("status", statusText)
	params.Add("visibility", configuration["visibility"])
	postData := strings.NewReader(params.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", configuration["fediInstance"]+"/api/v1/statuses", postData)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "bearer "+configuration["fediToken"])
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func PostStatus(status Status, configuration map[string]string) {
	mediaID := uploadImage(configuration, status.ImageURL)
	createStatus(configuration, mediaID, status)
}
