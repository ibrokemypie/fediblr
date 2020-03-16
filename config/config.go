package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	FediInstance string
	FediToken    string
	LastId       []int
	TumblrKey    string
	TumblrUser   string
	Visibility   string
}

func WriteConfig(config Config) {
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

func GetConfig() Config {
	config := Config{}
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
	}

	if config.TumblrKey == "" {
		config.TumblrKey = GetTumblrKey()
		WriteConfig(config)
	}

	if config.TumblrUser == "" {
		config.TumblrUser = GetTumblrUser()
		WriteConfig(config)
	}

	if config.FediInstance == "" {
		config.FediInstance = GetFediInstance()
		WriteConfig(config)
	}

	if config.FediToken == "" {
		config.FediToken = GetFediAccessToken(config.FediInstance)
		WriteConfig(config)
	}

	if config.Visibility == "" {
		config.Visibility = "unlisted"
		WriteConfig(config)
	}

	if config.LastId == nil {
		config.LastId = make([]int, 10, 10)
		WriteConfig(config)
	}

	return config
}
