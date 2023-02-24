package router

import (
	"fmt"
	"html/template"
	"log"
	"main/logic"
	"main/logic/types"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	ValidCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CommonExtNoNSH  = blackfriday.NoIntraEmphasis | blackfriday.Tables | blackfriday.FencedCode | blackfriday.Autolink | blackfriday.Strikethrough | blackfriday.HeadingIDs | blackfriday.BackslashLineBreak | blackfriday.DefinitionLists
)

var (
	SubCache = make(map[string]types.Subreddit)
)

func StartServer() {
	// region Template Engine

	TemplateEngine := html.New("./views", ".html")

	TemplateEngine.AddFuncMap(template.FuncMap{
		"contains": strings.Contains,
		"add": func(input int) int {
			return input + 1
		},
		"ugidgen": func() string {
			ubytes := make([]byte, 28)
			for i := 0; i < len(ubytes); i++ {
				ubytes[i] = ValidCharacters[rand.Intn(len(ValidCharacters))]
			}
			return string(ubytes)
		},
		"notcontains": func(input, of string) bool {
			return !strings.Contains(input, of)
		},
		"sanitize": func(input string) template.HTML {
			Markdown := blackfriday.Run([]byte(input), blackfriday.WithExtensions(CommonExtNoNSH))
			SHTML := bluemonday.UGCPolicy().SanitizeBytes(Markdown)
			return template.HTML(SHTML)
		},
		"replaceAmp": func(input string) string {
			return strings.Replace(input, "&amp;", "&", -1)
		},
		"fmtEpochDate": func(input float64) string {
			return time.Unix(int64(input), 0).Format("Created Jan 02, 2006")
		},
		"fmtHumanComma": func(input int64) string {
			return humanize.Comma(input)
		},
		"fmtHumanDate": func(input float64) string {
			return humanize.Time(time.Unix(int64(input), 0))
		},
		"toPercentage": func(input float64) string {
			return fmt.Sprintf("%.0f", input*100)
		},
	})

	// endregion

	router := fiber.New(fiber.Config{
		Prefork:     true,
		Views:       TemplateEngine,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	router.Use(
		logger.New(),
		recover.New(),
		helmet.New(helmet.Config{
			XSSProtection:      "1; mode=block",
			ContentTypeNosniff: "nosniff",
			XFrameOptions:      "DENY",
		}),
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),
	)

	// region Load Files

	router.Get("/js/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		ctx.Set("Content-Type", "application/javascript")
		return ctx.SendFile(fmt.Sprintf("js/%v", id))
	})

	router.Get("/css/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		ctx.Set("Content-Type", "text/css")
		return ctx.SendFile(fmt.Sprintf("css/%v", id))
	})

	// endregion

	router.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", nil)
	})

	router.Get("/config", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Will be implemented soon™")
	})

	// dev -> will probably keep this.
	router.Post("/byecookies", func(ctx *fiber.Ctx) error {
		ctx.Cookie(&fiber.Cookie{
			Name:     "nsfw_allowed",
			Value:    "0",
			Expires:  time.Now().Add((24 * time.Hour) * 365),
			Secure:   true,
			HTTPOnly: true,
			SameSite: "lax",
		})
		return ctx.RedirectBack("/config", http.StatusMovedPermanently)
	})

	// for now, it's only purpose is to set cookies to nsfw subreddits (expand later™)
	router.Post("/config", func(ctx *fiber.Ctx) error {
		ctx.Cookie(&fiber.Cookie{
			Name:     "nsfw_allowed",
			Value:    "1",
			Expires:  time.Now().Add((24 * time.Hour) * 365),
			Secure:   true,
			HTTPOnly: true,
			SameSite: "lax",
		})
		return ctx.RedirectBack("/config", http.StatusMovedPermanently)
	})

	router.Get("/r/:sub", func(ctx *fiber.Ctx) error {
		after := url.QueryEscape(ctx.Query("after"))
		sort := url.QueryEscape(ctx.Query("t"))
		subname := url.QueryEscape(ctx.Params("sub"))

		Posts := logic.GetPosts(after, sort, subname)

		if len(Posts.Data.Children) == 0 {
			return ctx.SendString(fmt.Sprintf("The subreddit 'r/%v' was banned, or doesn't exist. (Did you make a typo - exceeded the rate limit?)", subname))
		}

		// Cache subreddit data, so we don't have to keep making requests every single time.
		// This will store it in memory, which may not be the best, and a disk based cache would be better.
		var Sub types.Subreddit

		if scache, exists := SubCache[subname]; exists {
			Sub = scache
		} else {
			Sub = logic.GetSubredditData(subname)
			SubCache[subname] = Sub
		}

		SortPostData(&Posts)

		nsfwallowed := ctx.Cookies("nsfw_allowed")

		return ctx.Render("sub", fiber.Map{
			"SubData":     Sub.Data,
			"Posts":       Posts.Data,
			"NSFWAllowed": nsfwallowed == "1" || !Sub.Data.NSFW,
		})
	})

	router.Get("/r/:sub/loadPosts", func(ctx *fiber.Ctx) error {
		subname := url.QueryEscape(ctx.Params("sub"))
		after := url.QueryEscape(ctx.Query("after"))
		sort := url.QueryEscape(ctx.Query("t"))

		Posts := logic.GetPosts(after, sort, subname)

		SortPostData(&Posts)

		return ctx.Render("posts", fiber.Map{
			"Posts": Posts.Data,
		})
	})

	// NoRoute
	router.Use(func(ctx *fiber.Ctx) error {
		return ctx.Render("404", nil)
	})

	// localhost:9090
	log.Fatal(router.Listen(":9090"))
}

func SortPostData(Posts *types.Posts) {
	var dataChannel = make(chan types.Post, len(Posts.Data.Children))

	for i := 0; i < len(Posts.Data.Children); i++ {
		dataChannel <- Posts.Data.Children[i].Data
	}

	close(dataChannel)

	for i := 0; i < len(Posts.Data.Children); i++ {
		if Post, ok := <-dataChannel; ok {
			if Post.Preview.Images != nil {
				Image := Post.Preview.Images[0]

				if Image.Resolutions != nil {
					Mid := (len(Image.Resolutions) >> 1) + 1
					if Mid >= len(Image.Resolutions) {
						Mid = len(Image.Resolutions) - 1
					}
					Post.Preview.AutoChosenImageQuality = strings.Replace(Image.Resolutions[Mid].URL, "&amp;", "&", -1)
					Post.Preview.AutoChosenPosterQuality = strings.Replace(Post.Preview.AutoChosenImageQuality, "&amp;", "&", -1)
				} else {
					Post.Preview.AutoChosenImageQuality = strings.Replace(Image.Source.URL, "&amp;", "&", -1)
					Post.Preview.AutoChosenPosterQuality = strings.Replace(Post.Preview.AutoChosenImageQuality, "&amp;", "&", -1)
				}

				if strings.Contains(Image.Source.URL, ".gif") {
					if Image.Variants.MP4.Resolutions != nil {
						Mid := (len(Image.Variants.MP4.Resolutions) >> 1) + 1
						if Mid >= len(Image.Variants.MP4.Resolutions) {
							Mid = len(Image.Variants.MP4.Resolutions) - 1
						}
						Post.Preview.AutoChosenImageQuality = strings.Replace(Image.Variants.MP4.Resolutions[Mid].URL, "&amp;", "&", -1)
					} else {
						Post.Preview.AutoChosenImageQuality = strings.Replace(Image.Variants.MP4.Source.URL, "&amp;", "&", -1)
					}
				}
			}

			if Post.SecureMedia != nil && Post.SecureMedia.RedditVideo != nil {
				Post.SecureMedia.RedditVideo.LQ = fmt.Sprintf("%v/DASH_360.mp4", Post.LinkURL)
				Post.SecureMedia.RedditVideo.MQ = fmt.Sprintf("%v/DASH_480.mp4", Post.LinkURL)
			}

			if Post.MediaMetaData != nil {
				var MediaLinks []string

				if Post.GalleryData.Items != nil {
					for i := 0; i < len(Post.GalleryData.Items); i++ {
						ItemID := Post.GalleryData.Items[i].MediaID
						MediaData := Post.MediaMetaData[ItemID]
						if MediaData.P != nil {
							Mid := (len(MediaData.P) >> 1) + 1
							if Mid >= len(MediaData.P) {
								Mid = len(MediaData.P) - 1
							}
							MediaLinks = append(MediaLinks, strings.Replace(MediaData.P[Mid].U, "&amp;", "&", -1))
						}
					}
				} else {
					// range is random, therefore the images *may* be mixed up.
					// may, because there is a chance that images are in order, due to the randomness.
					// there is no way to sort this.
					for _, MediaData := range Post.MediaMetaData {
						if MediaData.P != nil {
							Mid := (len(MediaData.P) >> 1) + 1
							if Mid >= len(MediaData.P) {
								Mid = len(MediaData.P) - 1
							}
							MediaLinks = append(MediaLinks, strings.Replace(MediaData.P[Mid].U, "&amp;", "&", -1))
						}
					}
				}

				Post.VMediaMetaData = MediaLinks
			}

			if len(Post.SelfText) != 0 {
				// invisible character, blackfriday doesn't recognize it, and just displays &#x200B; which is pretty distracting in some cases.
				Post.SelfText = strings.Replace(Post.SelfText, "&amp;#x200B;", "", -1)
			}
			Posts.Data.Children[i].Data = Post
		}
	}
}
