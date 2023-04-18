package router

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"main/logic"
	"main/logic/types"

	"github.com/dustin/go-humanize"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	ValidCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CommonExtNoNSH  = blackfriday.NoIntraEmphasis | blackfriday.Tables | blackfriday.FencedCode | blackfriday.Autolink | blackfriday.Strikethrough | blackfriday.HeadingIDs | blackfriday.BackslashLineBreak | blackfriday.DefinitionLists

	JSCookie      = "JSEnabled"
	INFCookie     = "INFScroll"
	NSFWCookie    = "NSFWAllowed"
	ResCookie     = "PreferredResolution"
	GalleryCookie = "GalleryNav"
	USRCCookie    = "TrustUSrc"
	MathCookie    = "UseAdvMath"

	JSCookieValue      = "js_enabled"
	INFCookieValue     = "infscroll_enabled"
	NSFWCookieValue    = "nsfw_allowed"
	ResCookieValue     = "preferred_resolution"
	GalleryCookieValue = "gallery_navigation"
	USRCCookieValue    = "trust_unknownsources"
	MathCookieValue    = "advanced_math"
)

var (
	SubCache sync.Map

	CFGMap = map[string]string{
		"EnableJS":         JSCookieValue,
		"EnableInfScroll":  INFCookieValue,
		"AllowNSFW":        NSFWCookieValue,
		"PrefRes":          ResCookieValue,
		"EnableGalleryNav": GalleryCookieValue,
		"TrustUnknownSrc":  USRCCookieValue,
		"UseAdvancedMath":  MathCookieValue,
	}

	ValidImageExts = map[string]bool{
		".gif":  true,
		".png":  true,
		".jpg":  true,
		".jpeg": true,
	}

	RewritePath = map[string]string{
		"https://v.redd.it":                "/video",
		"https://i.redd.it":                "/image",
		"https://a.thumbs.redditmedia.com": "/athumb",
		"https://b.thumbs.redditmedia.com": "/bthumb",
		"https://external-preview.redd.it": "/external",
		"https://preview.redd.it":          "/preview",
		"https://styles.redditmedia.com":   "/rstyle",
		"https://www.redditstatic.com":     "/rstatic",
		"https://i.imgur.com":              "/imgur",
	}
)

func RewriteURL(input string) string {
	for k, v := range RewritePath {
		if strings.HasPrefix(input, k) {
			return v + input[len(k):]
		}
	}
	return input
}

func StartServer() {
	// region Template Engine

	templateEngine := html.New("./views", ".html")

	templateEngine.AddFuncMap(template.FuncMap{
		"contains":      strings.Contains,
		"sterilizepath": RewriteURL,
		"add": func(input int) int {
			return input + 1
		},
		"ugidgen": func() string {
			ubytes := make([]byte, 28)
			for i := 0; i < len(ubytes); i++ {
				ubytes[i] = ValidCharacters[rand.Intn(len(ValidCharacters))] //nolint:gosec,revive // We do not need to use crypto/rand here.
			}
			return string(ubytes)
		},
		"sanitize": func(input string) template.HTML {
			markdown := blackfriday.Run([]byte(input), blackfriday.WithExtensions(CommonExtNoNSH))
			sHTML := bluemonday.UGCPolicy().
				RequireNoFollowOnLinks(true).
				RequireNoReferrerOnLinks(true).
				AddTargetBlankToFullyQualifiedLinks(true).
				SanitizeBytes(markdown)
			return template.HTML(sHTML) //nolint:gosec,revive // bluemonday sanitizes this.
		},
		"qualifiesAsImg": func(input string) bool {
			return ValidImageExts[filepath.Ext(input)]
		},
		"fmtEpochDate": func(input float64) string {
			return time.Unix(int64(input), 0).Format("Created Jan 02, 2006")
		},
		"fmtHumanComma": humanize.Comma,
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
		Views:       templateEngine,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	router.Use(
		logger.New(),
		fiberrecover.New(),
		helmet.New(helmet.Config{
			XSSProtection:         "1; mode=block",
			ContentTypeNosniff:    "nosniff",
			XFrameOptions:         "DENY",
			ContentSecurityPolicy: "default-src 'self';form-action 'self';worker-src 'self' blob:;frame-ancestors 'none';script-src-attr 'self' 'unsafe-inline';style-src 'self' 'unsafe-inline';upgrade-insecure-requests",
			ReferrerPolicy:        "no-referrer",
		}),
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),
		func(ctx *fiber.Ctx) error {
			jsenabled := ctx.Cookies(JSCookieValue)
			infscrollenabled := ctx.Cookies(INFCookieValue)
			nsfwallowed := ctx.Cookies(NSFWCookieValue)
			gallerynav := ctx.Cookies(GalleryCookieValue)
			trustusrc := ctx.Cookies(USRCCookieValue)
			advmath := ctx.Cookies(MathCookieValue)

			ctx.Bind(fiber.Map{ //nolint:errcheck,gosec,revive // ctx.Bind always returns nil
				JSCookie:      jsenabled == "1",
				INFCookie:     infscrollenabled == "1",
				NSFWCookie:    nsfwallowed == "1",
				GalleryCookie: gallerynav == "1",
				USRCCookie:    trustusrc == "1",
				MathCookie:    advmath == "1",
			})

			return ctx.Next()
		},
	)

	router.Static("/js", "./js")
	router.Static("/css", "./css")
	router.Static("/fonts", "./fonts")

	router.Get("/:proxypath/*", func(ctx *fiber.Ctx) error {
		fullURL := ctx.Params("*")

		if index := strings.Index(ctx.OriginalURL(), "?"); index != 1 {
			fullURL += "?" + ctx.OriginalURL()[index+1:]
		}

		switch ctx.Params("proxypath") {
		case "video":
			if err := proxy.Do(ctx, "https://v.redd.it/"+fullURL); err != nil {
				return err
			}
		case "image":
			if err := proxy.Do(ctx, "https://i.redd.it/"+fullURL); err != nil {
				return err
			}
		case "athumb":
			if err := proxy.Do(ctx, "https://a.thumbs.redditmedia.com/"+fullURL); err != nil {
				return err
			}
		case "bthumb":
			if err := proxy.Do(ctx, "https://b.thumbs.redditmedia.com/"+fullURL); err != nil {
				return err
			}
		case "external":
			if err := proxy.Do(ctx, "https://external-preview.redd.it/"+fullURL); err != nil {
				return err
			}
		case "preview":
			if err := proxy.Do(ctx, "https://preview.redd.it/"+fullURL); err != nil {
				return err
			}
		case "rstyle":
			if err := proxy.Do(ctx, "https://styles.redditmedia.com/"+fullURL); err != nil {
				return err
			}
		case "rstatic":
			if err := proxy.Do(ctx, "https://www.redditstatic.com/"+fullURL); err != nil {
				return err
			}
		case "imgur":
			if err := proxy.Do(ctx, "https://i.imgur.com/"+fullURL); err != nil {
				return err
			}
		default:
			return ctx.Next()
		}

		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	router.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", nil)
	})

	router.Get("/config", func(ctx *fiber.Ctx) error {
		preferredres, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
		if err != nil {
			preferredres = 3
		}

		return ctx.Render("config", fiber.Map{
			ResCookie: preferredres,
		})
	})

	router.Post("/config", func(ctx *fiber.Ctx) error {
		for cookiekey, cookievalue := range CFGMap {
			switch formvalue := ctx.FormValue(cookiekey); formvalue {
			case "on":
				if cookiekey != "PrefRes" {
					setcfgCookie(ctx, cookievalue, "1")
				}
			case "off":
				if cookiekey != "PrefRes" {
					setcfgCookie(ctx, cookievalue, "0")
				}
			case "0", "1", "2", "3", "4", "5":
				if cookiekey == "PrefRes" {
					setcfgCookie(ctx, cookievalue, formvalue)
				}
			case "Source":
				if cookiekey == "PrefRes" {
					/*
						ctx.Cookies returns a string, which we will convert to
						int via strconv.Atoi, but if we set the cookie value to
						"Source", then it will error out, so set it to a high
						value that doesn't exist, but is still valid.
					*/

					setcfgCookie(ctx, cookievalue, "11037")
				}
			}
		}

		return ctx.RedirectBack("/config", fiber.StatusMovedPermanently)
	})

	router.Get("/r/:sub", func(ctx *fiber.Ctx) error {
		after := ctx.Query("after")
		flair := url.QueryEscape(ctx.Query("f"))
		subname := strings.ToLower(ctx.Params("sub"))

		posts := logic.GetPosts(subname, after, flair)

		if len(posts.Data.Children) == 0 {
			return ctx.Render("404", nil)
		}

		// Cache subreddit data, so we don't have to keep making requests every single time.
		// This will store it in memory, which may not be the best, and a disk based cache would be better.
		var sub types.Subreddit

		if scache, exists := SubCache.Load(subname); exists {
			sub = scache.(types.Subreddit) //nolint:errcheck,revive,forcetypeassert // This should not error out.
		} else {
			sub = logic.GetSubredditData(subname)
			SubCache.Store(subname, sub)
		}

		resolutionToUse, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
		if err != nil {
			resolutionToUse = 3
		}

		flairuesc, err := url.QueryUnescape(flair)
		if err != nil {
			log.Println(err)
		}

		SortPostData(&posts, resolutionToUse)

		return ctx.Render("sub", fiber.Map{
			"SubName":       subname,
			"SubData":       sub.Data,
			"Posts":         posts.Data,
			"FlairFiltered": flairuesc,
		})
	})

	router.Post("/loadPosts", func(ctx *fiber.Ctx) error {
		subname := ctx.FormValue("sub")
		after := ctx.FormValue("after")
		flair := url.QueryEscape(ctx.FormValue("flair"))

		posts := logic.GetPosts(subname, after, flair)

		resolutionToUse, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
		if err != nil {
			resolutionToUse = 3
		}

		flairuesc, err := url.QueryUnescape(flair)
		if err != nil {
			log.Println(err)
		}

		SortPostData(&posts, resolutionToUse)

		return ctx.Render("posts", fiber.Map{
			"SubName":       subname,
			"Posts":         posts.Data,
			"FlairFiltered": flairuesc,
		})
	})

	// NoRoute
	router.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).Render("404", nil)
	})

	// localhost:9090
	log.Fatal(router.Listen(":9090")) //nolint:revive // No
}

func SortPostData(posts *types.Posts, resolutionToUse int) {
	origRes := resolutionToUse
	for i := range posts.Data.Children {
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					log.Println("Recovered from fatal panic...", rec)
				}
			}()

			post := &posts.Data.Children[i].Data

			if len(post.Preview.Images) > 0 {
				image := post.Preview.Images[0]

				if len(image.Resolutions) > 0 && resolutionToUse != 11037 {
					if resolutionToUse >= len(image.Resolutions) {
						resolutionToUse = len(image.Resolutions) - 1
					}
					post.Preview.AutoChosenImageQuality = RewriteURL(image.Resolutions[resolutionToUse].URL)
					post.Preview.AutoChosenPosterQuality = post.Preview.AutoChosenImageQuality
				} else {
					post.Preview.AutoChosenImageQuality = RewriteURL(image.Source.URL)
					post.Preview.AutoChosenPosterQuality = post.Preview.AutoChosenImageQuality
				}

				resolutionToUse = origRes

				if strings.Contains(image.Source.URL, ".gif") {
					if len(image.Variants.MP4.Resolutions) > 0 && resolutionToUse != 11037 {
						if resolutionToUse >= len(image.Variants.MP4.Resolutions) {
							resolutionToUse = len(image.Variants.MP4.Resolutions) - 1
						}
						post.Preview.AutoChosenImageQuality = RewriteURL(image.Variants.MP4.Resolutions[resolutionToUse].URL)
					} else {
						post.Preview.AutoChosenImageQuality = RewriteURL(image.Variants.MP4.Source.URL)
					}
				}

				resolutionToUse = origRes
			}

			if len(post.MediaMetaData) > 0 {
				if len(post.GalleryData.Items) > 0 {
					for j := 0; j < len(post.GalleryData.Items); j++ {
						itemID := post.GalleryData.Items[j].MediaID
						mediaData := post.MediaMetaData[itemID]
						if resolutionToUse >= len(mediaData.P) {
							resolutionToUse = len(mediaData.P) - 1
						}

						post.VMediaMetaData = append(post.VMediaMetaData, vmediaappendor(mediaData, resolutionToUse))
						resolutionToUse = origRes
					}
				} else {
					// range is random, therefore the images *may* be mixed up.
					// may, because there is a chance that images are in order, due to the randomness.
					// there is no way to sort this.
					for _, MediaData := range post.MediaMetaData {
						if resolutionToUse >= len(MediaData.P) {
							resolutionToUse = len(MediaData.P) - 1
						}
						post.VMediaMetaData = append(post.VMediaMetaData, vmediaappendor(MediaData, resolutionToUse))
						resolutionToUse = origRes
					}
				}
			}

			if len(post.LinkURL) > 0 {
				post.LinkURL = RewriteURL(post.LinkURL)
			}

			post.SelfText = strings.ReplaceAll(post.SelfText, "&#x200B;", "")

			posts.Data.Children[i].Data = *post
		}()
	}
}

func vmediaappendor(mData types.InternalMetaData, resolutionToUse int) types.InternalVData {
	isVideo := len(mData.S.MP4) > 0
	var poster, source string

	if len(mData.P) > 0 {
		poster = RewriteURL(mData.P[resolutionToUse].U)
	}

	if isVideo {
		source = RewriteURL(mData.S.MP4)
	} else if resolutionToUse == 11037 {
		source = RewriteURL(mData.S.U)
	} else {
		source = poster
	}
	return types.InternalVData{
		Video:                   isVideo,
		Link:                    source,
		AutoChosenPosterQuality: poster,
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
