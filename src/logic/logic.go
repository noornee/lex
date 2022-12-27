package logic

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/logic/types"
	"net/http"
	"strings"
)

var (
	Client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}
)

func GetSubredditData(subreddit string) types.Subreddit {
	url := fmt.Sprintf("https://reddit.com/r/%v/about.json", subreddit)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("User-Agent", "go:lex:cmd777")

	resp, err := Client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var Sub types.Subreddit

	errum := json.Unmarshal(body, &Sub)
	if errum != nil {
		log.Println(errum)
	}

	return Sub
}

func GetPosts(after, sort, subreddit string) types.Posts {
	url := fmt.Sprintf("https://reddit.com/r/%v.json", subreddit)
	if len(sort) != 0 {
		url = fmt.Sprintf("https://reddit.com/r/%v/top.json?t=%v", subreddit, sort)
	}
	if len(after) != 0 {
		if strings.Contains(url, "?") {
			url += fmt.Sprintf("&after=%v", after)
		} else {
			url += fmt.Sprintf("?after=%v", after)
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("User-Agent", "go:lex:cmd777")

	resp, err := Client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var Posts types.Posts

	errum := json.Unmarshal(body, &Posts)
	if errum != nil {
		log.Println(errum)
	}

	return Posts
}
