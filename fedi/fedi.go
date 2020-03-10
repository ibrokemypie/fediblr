package fedi

import "fmt"

type Status struct {
	ImageURL      string
	Caption       string
	SourceName    string
	SourceURL     string
	RebloggedName string
	RebloggedURL  string
}

func PostStatus(status Status) {
	fmt.Println(status.ImageURL)
	fmt.Println(status.Caption)
	fmt.Println(status.SourceName + " " + status.SourceURL)
	fmt.Println(status.RebloggedName + " " + status.RebloggedURL)
}

// img
// text
// source
// reblogged from
