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
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	ValidCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CommonExtNoNSH  = blackfriday.NoIntraEmphasis | blackfriday.Tables | blackfriday.FencedCode | blackfriday.Autolink | blackfriday.Strikethrough | blackfriday.HeadingIDs | blackfriday.BackslashLineBreak | blackfriday.DefinitionLists

	JSCookie   = "JSEnabled"
	INFCookie  = "INFScroll"
	NSFWCookie = "NSFWAllowed"
	ResCookie  = "PreferredResolution"

	JSCookieValue   = "js_enabled"
	INFCookieValue  = "infscroll_enabled"
	NSFWCookieValue = "nsfw_allowed"
	ResCookieValue  = "preferred_resolution"
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
		csrf.New(csrf.Config{
			KeyLookup:      "form:csrf",
			ContextKey:     "csrf",
			CookieName:     "csrf_",
			CookieSecure:   true,
			CookieHTTPOnly: true,
			CookieSameSite: "lax",
			Expiration:     1 * time.Hour,
			KeyGenerator:   utils.UUID,
			ErrorHandler:   func(c *fiber.Ctx, err error) error { return fiber.ErrBadRequest },
		}),
		fiberrecover.New(),
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

	router.Get("/fonts/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		ctx.Set("Content-Type", "font/woff2")
		return ctx.SendFile(fmt.Sprintf("fonts/%v", id))
	})

	// endregion

	router.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", nil)
	})

	router.Get("/config", func(ctx *fiber.Ctx) error {
		jsenabled := ctx.Cookies(JSCookieValue)
		infscrollenabled := ctx.Cookies(INFCookieValue)
		nsfwallowed := ctx.Cookies(NSFWCookieValue)

		preferredres, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
		if err != nil {
			preferredres = 3
		}

		return ctx.Render("config", fiber.Map{
			JSCookie:   jsenabled == "1",
			INFCookie:  infscrollenabled == "1",
			NSFWCookie: nsfwallowed == "1",
			ResCookie:  preferredres,
			"csrf":     ctx.Locals("csrf"),
		})
	})

	router.Post("/config", func(ctx *fiber.Ctx) error {
		if ctx.FormValue("csrf") != ctx.Cookies("csrf_") {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if ctx.FormValue("EnableJS") == "on" {
			setcfgCookie(ctx, JSCookieValue, "1")
		} else if ctx.FormValue("EnableJS") == "off" {
			setcfgCookie(ctx, JSCookieValue, "0")
		}

		if ctx.FormValue("EnableInfScroll") == "on" {
			setcfgCookie(ctx, INFCookieValue, "1")
		} else if ctx.FormValue("EnableInfScroll") == "off" {
			setcfgCookie(ctx, INFCookieValue, "0")
		}

		if ctx.FormValue("AllowNSFW") == "on" {
			setcfgCookie(ctx, NSFWCookieValue, "1")
		} else if ctx.FormValue("AllowNSFW") == "off" {
			setcfgCookie(ctx, NSFWCookieValue, "0")
		}

		switch ctx.FormValue("PrefRes") {
		case "0":
			setcfgCookie(ctx, ResCookieValue, "0")
		case "1":
			setcfgCookie(ctx, ResCookieValue, "1")
		case "2":
			setcfgCookie(ctx, ResCookieValue, "2")
		case "3":
			setcfgCookie(ctx, ResCookieValue, "3")
		case "4":
			setcfgCookie(ctx, ResCookieValue, "4")
		case "5":
			setcfgCookie(ctx, ResCookieValue, "5")
		case "Source":
			/*
				ctx.Cookies returns a string, which we will convert to
				int via strconv.Atoi, but if we set the cookie value to
				"Source", then it will error out, so set it to a high
				value that doesn't exist, but is still valid.
			*/
			setcfgCookie(ctx, ResCookieValue, "11037")
		}

		return ctx.RedirectBack("/config", http.StatusMovedPermanently)
	})

	router.Get("/r/:sub", func(ctx *fiber.Ctx) error {
		after := url.QueryEscape(ctx.Query("after"))
		sort := url.QueryEscape(ctx.Query("t"))
		subname := url.QueryEscape(ctx.Params("sub"))

		Posts := logic.GetPosts(after, sort, subname)

		if len(Posts.Data.Children) == 0 {
			return ctx.Render("404", nil)
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

		ResolutionToUse, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
		if err != nil {
			ResolutionToUse = 3
		}

		SortPostData(&Posts, ResolutionToUse)

		jsenabled := ctx.Cookies(JSCookieValue)
		infscrollenabled := ctx.Cookies(INFCookieValue)
		nsfwallowed := ctx.Cookies(NSFWCookieValue)

		return ctx.Render("sub", fiber.Map{
			"SubName":  subname,
			"SubData":  Sub.Data,
			"Posts":    Posts.Data,
			JSCookie:   jsenabled == "1",
			INFCookie:  infscrollenabled == "1",
			NSFWCookie: nsfwallowed == "1" || !Sub.Data.NSFW,
			"csrf":     ctx.Locals("csrf"),
		})
	})

	router.Post("/loadPosts", func(ctx *fiber.Ctx) error {
		if ctx.FormValue("csrf") != ctx.Cookies("csrf_") {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		subname := url.QueryEscape(ctx.FormValue("sub"))
		after := url.QueryEscape(ctx.FormValue("after"))
		sort := url.QueryEscape(ctx.FormValue("t"))

		Posts := logic.GetPosts(after, sort, subname)

		ResolutionToUse, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
		if err != nil {
			ResolutionToUse = 3
		}

		SortPostData(&Posts, ResolutionToUse)

		infscrollenabled := ctx.Cookies(INFCookieValue)

		return ctx.Render("posts", fiber.Map{
			"SubName": subname,
			"Posts":   Posts.Data,
			INFCookie: infscrollenabled == "1",
			"csrf":    ctx.Locals("csrf"),
		})
	})

	// NoRoute
	router.Use(func(ctx *fiber.Ctx) error {
		return ctx.Render("404", nil)
	})

	// localhost:9090
	log.Fatal(router.Listen(":9090"))
}

func SortPostData(Posts *types.Posts, ResolutionToUse int) {
	for i, t := range Posts.Data.Children {
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					log.Println("Recovered from fatal panic...", rec)
				}
			}()

			Post := &t.Data

			if len(Post.Preview.Images) > 0 {
				Image := Post.Preview.Images[0]

				if len(Image.Resolutions) > 0 && ResolutionToUse != 11037 {
					if ResolutionToUse >= len(Image.Resolutions) {
						ResolutionToUse = len(Image.Resolutions) - 1
					}
					Post.Preview.AutoChosenImageQuality = Image.Resolutions[ResolutionToUse].URL
					Post.Preview.AutoChosenPosterQuality = Post.Preview.AutoChosenImageQuality
				} else {
					Post.Preview.AutoChosenImageQuality = Image.Source.URL
					Post.Preview.AutoChosenPosterQuality = Post.Preview.AutoChosenImageQuality
				}

				if strings.Contains(Image.Source.URL, ".gif") {
					if len(Image.Variants.MP4.Resolutions) > 0 && ResolutionToUse != 11037 {
						if ResolutionToUse >= len(Image.Variants.MP4.Resolutions) {
							ResolutionToUse = len(Image.Variants.MP4.Resolutions) - 1
						}
						Post.Preview.AutoChosenImageQuality = Image.Variants.MP4.Resolutions[ResolutionToUse].URL
					} else {
						Post.Preview.AutoChosenImageQuality = Image.Variants.MP4.Source.URL
					}
				}
			}

			if len(Post.MediaMetaData) > 0 {
				if len(Post.GalleryData.Items) > 0 {
					for j := 0; j < len(Post.GalleryData.Items); j++ {
						ItemID := Post.GalleryData.Items[j].MediaID
						MediaData := Post.MediaMetaData[ItemID]
						if ResolutionToUse >= len(MediaData.P) {
							ResolutionToUse = len(MediaData.P) - 1
						}

						Post.VMediaMetaData = append(Post.VMediaMetaData, vmediaappendor(MediaData, ResolutionToUse))
					}
				} else {
					// range is random, therefore the images *may* be mixed up.
					// may, because there is a chance that images are in order, due to the randomness.
					// there is no way to sort this.
					for _, MediaData := range Post.MediaMetaData {
						if ResolutionToUse >= len(MediaData.P) {
							ResolutionToUse = len(MediaData.P) - 1
						}
						Post.VMediaMetaData = append(Post.VMediaMetaData, vmediaappendor(MediaData, ResolutionToUse))
					}
				}
			}

			Posts.Data.Children[i].Data = *Post
		}()
	}
}

func vmediaappendor(MData types.InternalMetaData, ResolutionToUse int) types.InternalVData {
	IsVideo := len(MData.S.MP4) > 0
	var Poster, Source string

	if len(MData.P) > 0 {
		Poster = MData.P[ResolutionToUse].U
	}

	if IsVideo {
		Source = MData.S.MP4
	} else if ResolutionToUse == 11037 {
		Source = MData.S.U
	} else {
		Source = Poster
	}
	return types.InternalVData{
		Video:                   IsVideo,
		Link:                    Source,
		AutoChosenPosterQuality: Poster,
	}
}

func setcfgCookie(ctx *fiber.Ctx, cookiename, cookievalue string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     cookiename,
		Value:    cookievalue,
		Expires:  time.Now().Add((24 * time.Hour) * 365),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "lax",
	})
}
