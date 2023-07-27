package types

// The Subreddit and Posts themselves
type Subreddit struct {
	Data struct {
		Title               string  `json:"title"`
		PrimaryColor        string  `json:"primary_color"`
		DisplayNamePrefixed string  `json:"display_name_prefixed"`
		Description         string  `json:"public_description"`
		CommunityIcon       string  `json:"community_icon"`
		Banner              string  `json:"banner_background_image"`
		ActiveUserCount     float64 `json:"accounts_active"`
		MemberCount         float64 `json:"subscribers"`
		Created             float64 `json:"created"`
		NSFW                bool    `json:"over18"`
	} `json:"data"`
}

type Posts struct {
	Data struct {
		After    string `json:"after"`
		Children []struct {
			Data internalPostData `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type internalPostData struct {
	Title         string `json:"title"`
	SelfText      string `json:"selftext"`
	Body          string `json:"body"`
	SubNamePref   string `json:"subreddit_name_prefixed"`
	PostFlair     string `json:"link_flair_text"`
	PostFlairHex  string `json:"link_flair_background_color"`
	Author        string `json:"author"`
	AuthorFlair   string `json:"author_flair_text"`
	PostHint      string `json:"post_hint"`
	Distinguished string `json:"distinguished"`
	PostID        string `json:"id"`
	Permalink     string `json:"permalink"`
	LinkURL       string `json:"url"`

	// videos
	SecureMedia struct {
		RedditVideo struct {
			HLSURL      string `json:"hls_url"`
			FallbackURL string `json:"fallback_url"`
		} `json:"reddit_video"`
	} `json:"secure_media"`

	// images
	Preview struct {
		AutoChosenImageQuality  string
		AutoChosenPosterQuality string
		RedditVideoPreview      struct {
			HLSURL      string `json:"hls_url"`
			FallbackURL string `json:"fallback_url"`
		} `json:"reddit_video_preview"`
		Images []struct {
			Variants struct {
				GIF struct {
					Source struct {
						URL string `json:"url"`
					} `json:"source"`
				} `json:"gif"`
				MP4 struct {
					Source struct {
						URL string `json:"url"`
					} `json:"source"`
					Resolutions []struct {
						URL string `json:"url"`
					} `json:"resolutions"`
				} `json:"mp4"`
			} `json:"variants"`
			Source struct {
				URL string `json:"url"`
			} `json:"source"`
			Resolutions []struct {
				URL string `json:"url"`
			} `json:"resolutions"`
		} `json:"images"`
	} `json:"preview"`

	// gallery
	MediaMetaData map[string]InternalMetaData `json:"media_metadata"`

	GalleryData struct {
		Items []struct {
			MediaID string `json:"media_id"`
		} `json:"items"`
	} `json:"gallery_data"`

	VMediaMetaData []InternalVData
	CrossPost      []struct {
		// todo: finish
		Permalink string `json:"permalink"`
	} `json:"crosspost_parent_list"`

	Awardings []struct {
		AwardSubType string `json:"award_sub_type"`
		Name         string `json:"name"`
		ResizedIcons []struct {
			URL string `json:"url"`
		} `json:"resized_icons"`
		Count float64 `json:"count"`
	} `json:"all_awardings"`

	UpvoteRatio     float64 `json:"upvote_ratio"`
	Ups             float64 `json:"ups"`
	Created         float64 `json:"created"`
	CommentCount    float64 `json:"num_comments"`
	Pinned          bool    `json:"stickied"`
	Locked          bool    `json:"locked"`
	Archived        bool    `json:"archived"`
	NSFW            bool    `json:"over_18"`
	Spoiler         bool    `json:"spoiler"`
	OriginalContent bool    `json:"is_original_content"`
}

type InternalCommentData struct {
	Replies     any    `json:"replies"`
	Author      string `json:"author"`
	AuthorFlair string `json:"author_flair_text"`
	Body        string `json:"body"`
	Permalink   string `json:"permalink"`
	VReplies    []InternalCommentData
	Depth       float64 `json:"depth"`
	Ups         float64 `json:"ups"`
	Created     float64 `json:"created"`
}

// MediaMetaData - Galleries
type InternalVData struct {
	Link                    string
	AutoChosenPosterQuality string
	Video                   bool
}

type InternalMetaData struct {
	// M string `json:"m"`
	S struct {
		U   string `json:"u"`
		MP4 string `json:"mp4"`
	} `json:"s"`
	P []struct {
		U string `json:"u"`
	} `json:"p"`
	// HLSUrl string `json:"hlsUrl"`
}

type Comments struct {
	Data struct {
		Children []struct {
			Data InternalCommentData `json:"data"`
		} `json:"children"`
	} `json:"data"`
}
