package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ibrokemypie/fediblr/fedi"
	"github.com/ibrokemypie/fediblr/tumblr"
)

func main() {
	var config map[string]string
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
	}

	if config["tumblrKey"] == "" {
		config["tumblrKey"] = tumblr.GetTumblrKey()
		writeConfig(config)
	}

	if config["tumblrUser"] == "" {
		config["tumblrUser"] = tumblr.GetTumblrUser()
		writeConfig(config)
	}

	if config["fediInstance"] == "" {
		config["fediInstance"] = fedi.GetFediInstance()
		writeConfig(config)
	}

	if config["fediToken"] == "" {
		config["fediToken"] = fedi.GetFediAccessToken(config["fediInstance"])
		writeConfig(config)
	}

	status := tumblr.GetPost(config["tumblrKey"], config["tumblrUser"])

	fedi.PostStatus(status, config["fediInstance"], config["fediToken"])
}

func writeConfig(config map[string]string) {
	f, err := os.Create("config.toml")
	if err != nil {
		// failed to create/open the file
		log.Fatal(err)
	}
	if err := toml.NewEncoder(f).Encode(config); err != nil {
		// failed to encode
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		// failed to close the file
		log.Fatal(err)

	}
}
