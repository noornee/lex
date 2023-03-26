package types

type Subreddit struct {
	Data struct {
		Title               string  `json:"title"`
		PrimaryColor        string  `json:"primary_color"`
		DisplayName         string  `json:"display_name"`
		DisplayNamePrefixed string  `json:"display_name_prefixed"`
		ActiveUserCount     int64   `json:"accounts_active"`
		MemberCount         int64   `json:"subscribers"`
		Description         string  `json:"public_description"`
		CommunityIcon       string  `json:"community_icon"`
		Banner              string  `json:"banner_background_image"`
		Created             float64 `json:"created"`
		NSFW                bool    `json:"over18"`
	} `json:"data"`
}

type Posts struct {
	Data struct {
		After string `json:"after"`
		Dist  int64  `json:"dist"`
		//GeoFilter any    `json:"geo_filter"` todo
		Children []struct {
			Data struct {
				Title        string  `json:"title"`
				SelfText     string  `json:"selftext"`
				PostFlair    string  `json:"link_flair_text"`
				PostFlairHex string  `json:"link_flair_background_color"`
				UpvoteRatio  float64 `json:"upvote_ratio"`
				Ups          int64   `json:"ups"`
				Created      float64 `json:"created"`
				Pinned       bool    `json:"stickied"`
				Locked       bool    `json:"locked"`
				Archived     bool    `json:"archived"`
				TotalAwards  int64   `json:"total_awards_received"`
				Awardings    []struct {
					AwardSubType string `json:"award_sub_type"`
					Count        int64  `json:"count"`
					Name         string `json:"name"`
					ResizedIcons []struct {
						URL string `json:"url"`
					} `json:"resized_icons"`
				} `json:"all_awardings"`
				Author        string `json:"author"`
				AuthorFlair   string `json:"author_flair_text"`
				NSFW          bool   `json:"over_18"`
				Spoiler       bool   `json:"spoiler"`
				PostHint      string `json:"post_hint"`
				Distinguished string `json:"distinguished"`
				PostID        string `json:"id"`
				CommentCount  int64  `json:"num_comments"`
				Permalink     string `json:"permalink"`
				LinkURL       string `json:"url"`
				Subreddit     string `json:"subreddit"`

				// videos
				SecureMedia struct {
					RedditVideo struct {
						HLSURL      string `json:"hls_url"`
						FallbackURL string `json:"fallback_url"`
					} `json:"reddit_video"`
				} `json:"secure_media"`

				// images
				Preview struct {
					Images []struct {
						Source struct {
							URL string `json:"url"`
						} `json:"source"`
						Resolutions []struct {
							URL string `json:"url"`
						} `json:"resolutions"`
						Variants struct {
							MP4 struct {
								Source struct {
									URL string `json:"url"`
								} `json:"source"`
								Resolutions []struct {
									URL string `json:"url"`
								} `json:"resolutions"`
							} `json:"mp4"`
							GIF struct {
								Source struct {
									URL string `json:"url"`
								} `json:"source"`
							} `json:"gif"`
						} `json:"variants"`
					} `json:"images"`
					RedditVideoPreview struct {
						HLSURL      string `json:"hls_url"`
						FallbackURL string `json:"fallback_url"`
					} `json:"reddit_video_preview"`
					AutoChosenImageQuality  string
					AutoChosenPosterQuality string
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
					//todo: finish
					Permalink string `json:"permalink"`
				} `json:"crosspost_parent_list"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type InternalVData struct {
	Video                   bool
	Link                    string
	AutoChosenPosterQuality string
}

type InternalMetaData struct {
	//M string `json:"m"`
	P []struct {
		U string `json:"u"`
	} `json:"p"`
	S struct {
		U   string `json:"u"`
		MP4 string `json:"mp4"`
	} `json:"s"`
	//HLSUrl string `json:"hlsUrl"`
}
