package logic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cmd777/lex/src/logic/types"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/utils"
)

func GetSubredditData(subreddit string) types.Subreddit {
	url := fmt.Sprintf("https://www.reddit.com/r/%v/about.json?raw_json=1", subreddit)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		log.Errorf("Failed to create New Request with Context: %w", err)
	}

	req.Header.Set("User-Agent", "go:lex:cmd777-with-"+utils.UUIDv4())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.reddit.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("DefaultClient failed to do the request: %w", err)
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Errorf("Failed to close response body: %w", closeerr)
		}
	}()

	var sub types.Subreddit

	if err = json.NewDecoder(resp.Body).Decode(&sub); err != nil {
		log.Errorf("json NewDecoder failed to decode JSON: %w", err)
	}

	return sub
}

func GetPosts(subreddit, after, flair string) types.Posts {
	url := fmt.Sprintf("https://www.reddit.com/r/%v", subreddit)
	if flair != "" {
		// stupid.
		// https://www.reddit.com/r/ModSupport/comments/hpf6na/filtering_by_flair_broken_for_some_users_on/
		url += fmt.Sprintf(`/search.json?raw_json=1&q=flair:"%v"&restrict_sr=1&sr_nsfw=1&include_over_18=1`, flair)
	} else {
		url += ".json?raw_json=1"
	}

	if after != "" {
		url += fmt.Sprintf("&after=%v", after)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		log.Errorf("Failed to create New Request with Context: %w", err)
	}

	req.Header.Set("User-Agent", "go:lex:cmd777-with-"+utils.UUIDv4())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.reddit.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("DefaultClient failed to do the request: %w", err)
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Errorf("Failed to close response body: %w", closeerr)
		}
	}()

	var posts types.Posts

	if err = json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		log.Errorf("json NewDecoder failed to decode JSON: %w", err)
	}

	return posts
}

func GetComments(subreddit, id string) (types.Posts, types.CommentsToUnmarshal) {
	url := fmt.Sprintf("https://www.reddit.com/r/%v/comments/%v.json?raw_json=1", subreddit, id)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		log.Errorf("Failed to create New Request with Context: %w", err)
	}

	req.Header.Set("User-Agent", "go:lex:cmd777-with-"+utils.UUIDv4())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.reddit.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("DefaultClient failed to do the request: %w", err)
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Errorf("Failed to close response body: %w", closeerr)
		}
	}()

	var commentsunmarshal types.CommentsToUnmarshal

	if err = json.NewDecoder(resp.Body).Decode(&commentsunmarshal.Data); err != nil {
		log.Errorf("json NewDecoder failed to decode JSON: %w", err)
	}

	var post types.Posts
	var comments types.Comments

	if err = json.Unmarshal(commentsunmarshal.Data[0], &post); err != nil {
		log.Errorf("json failed to unmarshal JSON: %w", err)
	}

	if err = json.Unmarshal(commentsunmarshal.Data[1], &comments); err != nil {
		log.Errorf("json failed to unmarshal JSON: %w", err)
	}

	internalDecode(&comments)

	return post, commentsunmarshal
}

func GetAccount(name, after string) types.Posts {
	url := fmt.Sprintf("https://www.reddit.com/u/%v.json?raw_json=1", name)

	if after != "" {
		url += fmt.Sprintf("&after=%v", after)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		log.Errorf("Failed to create New Request with Context: %w", err)
	}

	req.Header.Set("User-Agent", "go:lex:cmd777-with-"+utils.UUIDv4())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.reddit.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("DefaultClient failed to do the request: %w", err)
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Errorf("Failed to close response body: %w", closeerr)
		}
	}()

	var posts types.Posts

	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		log.Errorf("json NewDecoder failed to decode JSON: %w", err)
	}

	return posts
}

func internalDecode(comments *types.Comments) {
	/*
		CURRENT PLAN:
		First, we construct a similar type of that to the JSON's
		Then, we decode the Replies field to that, and iterate until it's over
			(instead of blindly guessing)
		Then, we do the same iteration in the HTML as well
			(this will be more flexible, as it will be easier to indent, and also may help provide a way to collapse it.)

		ISSUE:
		This is ancient code. DO NOT TOUCH material.
		I do not remember how it fully works.
	*/
	for _, v := range comments.Data.Children {
		var newdecoded types.Comments

		if err := json.Unmarshal(v.Data.Replies, &newdecoded); err == nil {
			// No decoding failure.
			internalDecode(&newdecoded)
		}
	}
}
