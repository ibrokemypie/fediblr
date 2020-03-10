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
	"path"
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

func uploadImage(imageURL string, instanceUrl string, authToken string) string {
	resp, err := http.Get(imageURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	req, err := newfileUploadRequest(instanceUrl+"/api/v1/media", nil, "file", resp.Body, path.Base(imageURL))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "bearer "+authToken)

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

func PostStatus(status Status, instanceURL string, authToken string) {
	mediaID := uploadImage(status.ImageURL, instanceURL, authToken)

	fmt.Println(mediaID)
	// fmt.Println(status.ImageURL)
	// fmt.Println(status.Caption)
	// fmt.Println(status.SourceName + " " + status.SourceURL)
	// fmt.Println(status.RebloggedName + " " + status.RebloggedURL)
}

// img
// text
// source
// reblogged from
