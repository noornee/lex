package logic

import (
	"crypto/tls"
	"fmt"
	"log"
	"main/logic/types"
	"net/http"

	"github.com/goccy/go-json"
)

var (
	Client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}
)

func GetSubredditData(subreddit string) types.Subreddit {
	url := fmt.Sprintf("https://reddit.com/r/%v/about.json?raw_json=1", subreddit)

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

	var Sub types.Subreddit

	err = json.NewDecoder(resp.Body).Decode(&Sub)
	if err != nil {
		log.Println(err)
	}

	return Sub
}

func GetPosts(after, sort, subreddit string) types.Posts {
	url := fmt.Sprintf("https://reddit.com/r/%v.json?raw_json=1", subreddit)
	if len(sort) != 0 {
		url = fmt.Sprintf("https://reddit.com/r/%v/top.json?raw_json=1&t=%v", subreddit, sort)
	}
	if len(after) != 0 {
		url += fmt.Sprintf("&after=%v", after)
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

	var Posts types.Posts

	err = json.NewDecoder(resp.Body).Decode(&Posts)
	if err != nil {
		log.Println(err)
	}

	return Posts
}
