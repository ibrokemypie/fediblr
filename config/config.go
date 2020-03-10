package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

func WriteConfig(config map[string]string) {
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

func GetConfig() map[string]string {
	config := make(map[string]string)
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
	}

	if config["tumblrKey"] == "" {
		config["tumblrKey"] = GetTumblrKey()
		WriteConfig(config)
	}

	if config["tumblrUser"] == "" {
		config["tumblrUser"] = GetTumblrUser()
		WriteConfig(config)
	}

	if config["fediInstance"] == "" {
		config["fediInstance"] = GetFediInstance()
		WriteConfig(config)
	}

	if config["fediToken"] == "" {
		config["fediToken"] = GetFediAccessToken(config["fediInstance"])
		WriteConfig(config)
	}

	if config["visibility"] == "" {
		config["visibility"] = "unlisted"
		WriteConfig(config)
	}

	return config
}
