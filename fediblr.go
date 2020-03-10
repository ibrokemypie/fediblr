package main

import (
	"os"

	"github.com/ibrokemypie/fediblr/fedi"
	"github.com/ibrokemypie/fediblr/tumblr"
)

func main() {
	if len(os.Args) != 3 {
		panic("Must have 2 arguments (api key and blog url)")
	}
	apiKey := os.Args[1]
	blogURL := os.Args[2]

	status := tumblr.GetPost(apiKey, blogURL)

	fedi.PostStatus(status)
}
