package main

import (
	"os"

	"github.com/ibrokemypie/fediblr/fedi"
	"github.com/ibrokemypie/fediblr/tumblr"
)

func main() {
	if len(os.Args) != 4 {
		panic("Must have 3 arguments (api key, blog url and fedi isntance)")
	}
	apiKey := os.Args[1]
	blogURL := os.Args[2]
	instanceURL := os.Args[3]

	status := tumblr.GetPost(apiKey, blogURL)

	fedi.PostStatus(status, instanceURL)
}
