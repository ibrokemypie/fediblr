package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Top struct {
	Meta     Meta     `json:"meta"`
	Response Response `json:"response"`
}

type Meta struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
}

type Response struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	Photos        []Photo `json:"photos"`
	SourceName    string  `json:"source_title"`
	SourceLink    string  `json:"source_url"`
	RebloggedName string  `json:"reblogged_from_name"`
	RebloggedLink string  `json:"reblogged_from_url"`
}

type Photo struct {
	Caption  string `json:"Caption"`
	Original Size   `json:"original_size"`
}

type Size struct {
	Link string `json:"URL"`
}

func main() {
	if len(os.Args) != 3 {
		panic("Must have 2 arguments (api key and blog url)")
	}
	apiKey := os.Args[1]
	blogURL := os.Args[2]

	baseURL := "https://api.tumblr.com/v2/blog/" + blogURL + "/posts/photo" +
		"?api_key=" + apiKey + "&limit=10&reblog_info=true"

	resp, err := http.Get(baseURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result Top

	json.NewDecoder(resp.Body).Decode(&result)

	if result.Meta.Status != 200 {
		panic(result.Meta.Message)
	}

	if len(result.Response.Posts) <= 0 {
		panic(result.Response.Posts)
	}

	var post Post
	for i := 0; i <= len(result.Response.Posts); i++ {
		if len(result.Response.Posts[i].Photos) > 0 {
			post = result.Response.Posts[i]
			break
		}
	}

	fmt.Println(post.Photos[0].Original.Link)
	fmt.Println(post.Photos[0].Caption)
	fmt.Println(post.SourceName + " " + post.SourceLink)
	fmt.Println(post.RebloggedName + " " + post.RebloggedLink)

}

// img
// text
// source
// reblogged from
