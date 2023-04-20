package logic

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"main/logic/types"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/utils"
)

func GetSubredditData(subreddit string) types.Subreddit {
	url := fmt.Sprintf("https://www.reddit.com/r/%v/about.json?raw_json=1", subreddit)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("User-Agent", "go:lex:cmd777-with-"+utils.UUIDv4())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.reddit.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Println("Failed to close response body", closeerr)
		}
	}()

	var sub types.Subreddit

	err = json.NewDecoder(resp.Body).Decode(&sub)
	if err != nil {
		log.Println(err)
	}

	return sub
}

func GetPosts(subreddit, after, flair string) types.Posts {
	url := fmt.Sprintf("https://www.reddit.com/r/%v", subreddit)
	if len(flair) != 0 {
		// stupid.
		// https://www.reddit.com/r/ModSupport/comments/hpf6na/filtering_by_flair_broken_for_some_users_on/
		url += fmt.Sprintf(`/search.json?raw_json=1&q=flair:"%v"&restrict_sr=1&sr_nsfw=1&include_over_18=1`, flair)
	} else {
		url += ".json?raw_json=1"
	}

	if len(after) != 0 {
		url += fmt.Sprintf("&after=%v", after)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("User-Agent", "go:lex:cmd777-with-"+utils.UUIDv4())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.reddit.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Println("Failed to close response body", closeerr)
		}
	}()

	var posts types.Posts

	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		log.Println(err)
	}

	return posts
}

func GetComments(subreddit, id string) (types.Post, []types.InternalCommentData) {
	url := fmt.Sprintf("https://www.reddit.com/r/%v/comments/%v.json?raw_json=1", subreddit, id)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("User-Agent", "go:lex:cmd777-with-"+utils.UUIDv4())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.reddit.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Println("Failed to close response body", closeerr)
		}
	}()

	var commentsunmarshal types.CommentsToUnmarshal

	err = json.NewDecoder(resp.Body).Decode(&commentsunmarshal.Data)
	if err != nil {
		log.Println(err)
	}

	var post types.Post
	var comments types.Comments

	err = json.Unmarshal(commentsunmarshal.Data[0], &post)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(commentsunmarshal.Data[1], &comments)
	if err != nil {
		log.Println(err)
	}

	internalDecode(&comments)

	return post, comments.MReplies
}

func internalDecode(comments *types.Comments) {
	for _, v := range comments.Data.Children {
		comments.MReplies = append(comments.MReplies, v.Data)

		var newdecoded types.Comments

		if err := json.Unmarshal(v.Data.Replies, &newdecoded); err == nil {
			// No decoding failure.
			internalDecode(&newdecoded)
			comments.MReplies = append(comments.MReplies, newdecoded.MReplies...)
		}
	}
}
