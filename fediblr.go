package main

import (
	"github.com/ibrokemypie/fediblr/config"
	"github.com/ibrokemypie/fediblr/fedi"
	"github.com/ibrokemypie/fediblr/tumblr"
)

func main() {
	configuration := config.GetConfig()

	status := tumblr.GetPost(configuration["tumblrKey"], configuration["tumblrUser"])

	fedi.PostStatus(status, configuration)
}
