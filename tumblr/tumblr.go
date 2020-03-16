package tumblr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/ibrokemypie/fediblr/config"
	"github.com/ibrokemypie/fediblr/fedi"
)

type Top struct {
	Meta     Meta     `json:"meta"`
	Response Response `json:"response"`
}

type Meta struct {
	ResponseStatus int    `json:"status"`
	Message        string `json:"msg"`
}

type Response struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	Type          string  `json:"type"`
	ID            int     `json:"id"`
	Photos        []Photo `json:"photos"`
	SourceName    string  `json:"source_title"`
	SourceLink    string  `json:"source_url"`
	RebloggedName string  `json:"reblogged_from_name"`
	RebloggedLink string  `json:"reblogged_from_url"`
	Summary       string  `json:"summary"`
	Caption       string  `json:"caption"`
	Body          string  `json:"body"`
}

type Photo struct {
	Caption  string `json:"Caption"`
	Original Size   `json:"original_size"`
}

type Size struct {
	Link string `json:"URL"`
}

func GetPost(configuration config.Config) fedi.Status {
	baseURL := "https://api.tumblr.com/v2/blog/" + configuration.TumblrUser + "/posts" +
		"?api_key=" + configuration.TumblrKey + "&limit=10&reblog_info=true"

	resp, err := http.Get(baseURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result Top

	json.NewDecoder(resp.Body).Decode(&result)

	if result.Meta.ResponseStatus != 200 {
		fmt.Println("http failure, code: " + strconv.Itoa(result.Meta.ResponseStatus) + " message: " + result.Meta.Message)
		os.Exit(1)
	}

	if len(result.Response.Posts) <= 0 {
		fmt.Println("no posts in the api response")
		os.Exit(1)
	}

	var post Post

	for i := len(result.Response.Posts)/2 - 1; i >= 0; i-- {
		opp := len(result.Response.Posts) - 1 - i
		result.Response.Posts[i], result.Response.Posts[opp] = result.Response.Posts[opp], result.Response.Posts[i]
	}

posts:
	for _, p := range result.Response.Posts {
		if len(configuration.LastId) > 0 {
			for _, v := range configuration.LastId {
				if p.ID == v {
					continue posts
				}
			}
		}

		if len(configuration.LastId) >= 10 {
			configuration.LastId = configuration.LastId[1:]
		}
		configuration.LastId = append(configuration.LastId, p.ID)

		config.WriteConfig(configuration)

		if p.Type == "photo" {
			if len(p.Photos) <= 0 {
				continue
			}
		}

		post = p
		break
	}

	if post.ID <= 0 {
		fmt.Println("No new posts")
		os.Exit(1)
	}

	post, caption := sanitizeCaption(post)

	images := []string{}
	for _, s := range post.Photos {
		images = append(images, s.Original.Link)
	}

	status := fedi.Status{
		Images:        images,
		Caption:       strip.StripTags(caption),
		SourceName:    post.SourceName,
		SourceURL:     removeRedirect(post.SourceLink),
		RebloggedName: post.RebloggedName,
		RebloggedURL:  removeRedirect(post.RebloggedLink)}

	return status
}

func removeRedirect(urlString string) string {
	if strings.HasPrefix(urlString, "https://t.umblr.com/redirect") {
		urlString = strings.TrimPrefix(urlString, "https://t.umblr.com/redirect?z=")
		urlString = strings.Split(urlString, "&t=")[0]
		var err error
		urlString, err = url.QueryUnescape(urlString)
		if err != nil {
			panic(err)
		}
	}

	return urlString
}

func sanitizeCaption(post Post) (Post, string) {
	var caption string
	if post.Body != "" {
		caption = post.Body
	} else if post.Caption != "" {
		caption = post.Caption
	} else {
		return post, post.Summary
	}

	caption = strings.Replace(caption, "<br>", "\n", -1)
	caption = strings.Replace(caption, "<p>", "\n", -1)

	if strings.Contains(caption, "img src=\"") {
		images := strings.Split(caption, "img src=\"")
		for _, v := range images {
			if strings.HasPrefix(v, "http") {
				newSize := Size{Link: strings.Split(v, "\"")[0]}
				newPhoto := Photo{Original: newSize}
				post.Photos = append(post.Photos, newPhoto)
			}
		}
	}

	return post, strip.StripTags(caption)
}
