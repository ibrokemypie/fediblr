package main

import (
	"github.com/ibrokemypie/fediblr/config"
	"github.com/ibrokemypie/fediblr/fedi"
	"github.com/ibrokemypie/fediblr/tumblr"
)

func main() {
	configuration := config.GetConfig()

	status := tumblr.GetPost(configuration)

	fedi.PostStatus(status, configuration)
}
