package router

import (
	"embed"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cmd777/lex/src/logic"
	"github.com/cmd777/lex/src/logic/types"
	"github.com/valyala/fasthttp"

	"github.com/dustin/go-humanize"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

var (
	//go:embed views
	viewFS embed.FS

	//go:embed js
	jsFS embed.FS

	//go:embed css
	cssFS embed.FS

	//go:embed fonts
	fontsFS embed.FS

	//go:embed favicon.ico
	faviconFS embed.FS

	jsonCache sync.Map

	resolutionHelper = []struct {
		Name  string
		Value int
	}{
		{Name: "Very Low", Value: 1},
		{Name: "Low", Value: 2},
		{Name: "Average", Value: 3},
		{Name: "Good", Value: 4},
		{Name: "Very Good", Value: 5},
		{Name: "Source", Value: 11037},
	}

	themeHelper = []struct {
		Name  string
		Value string
	}{
		{Name: "Dark (default)", Value: "dark"},
		{Name: "Dracula", Value: "dracula"},
		{Name: "Solarized", Value: "solarized"},
		{Name: "Abyss", Value: "abyss"},
		{Name: "Monokai", Value: "monokai"},
		{Name: "Red", Value: "redtheme"},
		{Name: "Light", Value: "lighttheme"},
	}
)

const (
	ValidCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CommonExtNoNSH  = blackfriday.NoIntraEmphasis | blackfriday.Tables | blackfriday.FencedCode | blackfriday.Autolink | blackfriday.Strikethrough | blackfriday.HeadingIDs | blackfriday.BackslashLineBreak | blackfriday.DefinitionLists

	JSCookie      = "JSEnabled"
	INFCookie     = "INFScroll"
	NSFWCookie    = "NSFWAllowed"
	ResCookie     = "PreferredResolution"
	ThemeCookie   = "PreferredTheme"
	GalleryCookie = "GalleryNav"
	USRCCookie    = "TrustUSrc"
	MathCookie    = "UseAdvMath"
	AwardCookie   = "DisableAwards"
	CommentCookie = "DisableComments"

	JSCookieValue      = "js_enabled"
	INFCookieValue     = "infscroll_enabled"
	NSFWCookieValue    = "nsfw_allowed"
	ResCookieValue     = "preferred_resolution"
	ThemeCookieValue   = "preferred_theme"
	GalleryCookieValue = "gallery_navigation"
	USRCCookieValue    = "trust_unknownsources"
	MathCookieValue    = "advanced_math"
	AwardCookieValue   = "disable_awards"
	CommentCookieValue = "disable_comments"

	MaxResolution = 11037

	CleanTimer = 2 * time.Minute
)

func RewriteURL(input string) string {
	switch {
	case strings.HasPrefix(input, "https://v.redd.it"):
		return "/video" + input[17:]
	case strings.HasPrefix(input, "https://i.redd.it"):
		return "/image" + input[17:]
	case strings.HasPrefix(input, "https://a.thumbs.redditmedia.com"):
		return "/athumb" + input[32:]
	case strings.HasPrefix(input, "https://b.thumbs.redditmedia.com"):
		return "/bthumb" + input[32:]
	case strings.HasPrefix(input, "https://external-preview.redd.it"):
		return "/external" + input[32:]
	case strings.HasPrefix(input, "https://preview.redd.it"):
		return "/preview" + input[23:]
	case strings.HasPrefix(input, "https://styles.redditmedia.com"):
		return "/rstyle" + input[30:]
	case strings.HasPrefix(input, "https://www.redditstatic.com"):
		return "/rstatic" + input[28:]
	case strings.HasPrefix(input, "https://i.imgur.com"):
		return "/imgur" + input[19:]
	default:
		return input
	}
}

func UGIDGen() string {
	ubytes := make([]byte, 28)
	for i := 0; i < len(ubytes); i++ {
		ubytes[i] = ValidCharacters[rand.Intn(len(ValidCharacters))] //nolint:gosec // We do not need to use crypto/rand here.
	}
	return string(ubytes)
}

func Sanitize(input string) template.HTML {
	markdown := blackfriday.Run([]byte(input), blackfriday.WithExtensions(CommonExtNoNSH))
	sHTML := bluemonday.UGCPolicy().
		RequireNoFollowOnLinks(true).
		RequireNoReferrerOnLinks(true).
		AddTargetBlankToFullyQualifiedLinks(true).
		SanitizeBytes(markdown)
	return template.HTML(sHTML) //nolint:gosec // bluemonday sanitizes this.
}

func QualifiesAsImg(input string) bool {
	switch filepath.Ext(input) {
	case ".gif":
		return true
	case ".png":
		return true
	case ".jpg":
		return true
	case ".jpeg":
		return true
	default:
		return false
	}
}

func FmtEpochDate(input float64) string {
	return time.Unix(int64(input), 0).Format("Created Jan 02, 2006")
}

func Incrementbyone(input int) int {
	return input + 1
}

func FmtHumanDate(input float64) string {
	return humanize.Time(time.Unix(int64(input), 0))
}

func ToPercentage(input float64) string {
	return fmt.Sprintf("%.0f", input*100)
}

func StartServer() {
	go backgroundJanitor()
	startTime := time.Now()

	var subCache sync.Map
	defClient := &fasthttp.Client{
		ReadBufferSize:           8192,
		NoDefaultUserAgentHeader: true,
	}

	cfgMap := map[string]string{
		"EnableJS":         JSCookieValue,
		"EnableInfScroll":  INFCookieValue,
		"AllowNSFW":        NSFWCookieValue,
		"EnableGalleryNav": GalleryCookieValue,
		"TrustUnknownSrc":  USRCCookieValue,
		"UseAdvancedMath":  MathCookieValue,
		"BlockAwards":      AwardCookieValue,
		"DontLoadComments": CommentCookieValue,
	}

	// region Template Engine

	templateEngine := html.NewFileSystem(http.FS(viewFS), ".html")

	templateEngine.AddFuncMap(template.FuncMap{
		"contains":       strings.Contains,
		"sterilizepath":  RewriteURL,
		"add":            Incrementbyone,
		"ugidgen":        UGIDGen,
		"sanitize":       Sanitize,
		"qualifiesAsImg": QualifiesAsImg,
		"fmtEpochDate":   FmtEpochDate,
		"fmtHumanComma":  humanize.Comma,
		"fmtHumanDate":   FmtHumanDate,
		"toPercentage":   ToPercentage,
	})
	// endregion

	router := fiber.New(fiber.Config{
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
			disableawards := ctx.Cookies(AwardCookieValue)
			disablecomments := ctx.Cookies(CommentCookieValue)

			theme := "dark"
			for i := range themeHelper {
				if ctx.Cookies(ThemeCookieValue) == themeHelper[i].Value {
					theme = ctx.Cookies(ThemeCookieValue)
					break
				}
			}

			ctx.Bind(fiber.Map{ //nolint:errcheck,gosec // ctx.Bind always returns nil
				JSCookie:      jsenabled == "1",
				INFCookie:     infscrollenabled == "1",
				NSFWCookie:    nsfwallowed == "1",
				GalleryCookie: gallerynav == "1",
				USRCCookie:    trustusrc == "1",
				MathCookie:    advmath == "1",
				AwardCookie:   disableawards == "1",
				CommentCookie: disablecomments == "1",
				ThemeCookie:   theme,
			})

			return ctx.Next()
		},
	)

	router.Use("/js", filesystem.New(filesystem.Config{
		Root:       http.FS(jsFS),
		PathPrefix: "js",
	}))

	router.Use("/css", filesystem.New(filesystem.Config{
		Root:       http.FS(cssFS),
		PathPrefix: "css",
	}))

	router.Use("/fonts", filesystem.New(filesystem.Config{
		Root:       http.FS(fontsFS),
		PathPrefix: "fonts",
	}))

	router.Use("/favicon.ico", filesystem.New(filesystem.Config{
		Root:       http.FS(faviconFS),
		PathPrefix: "favicon.ico",
	}))

	router.Get("/video/*", func(ctx *fiber.Ctx) error {
		if err := proxy.Do(ctx, "https://v.redd.it/"+ctx.OriginalURL()[7:], defClient); err != nil {
			return err
		}

		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	router.Get("/image/*", func(ctx *fiber.Ctx) error {
		ctx.Request().Header.Set("Accept", "image/avif,image/webp,*/*")

		if err := proxy.Do(ctx, "https://i.redd.it/"+ctx.OriginalURL()[7:], defClient); err != nil {
			return err
		}

		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	router.Get("/athumb/*", func(ctx *fiber.Ctx) error {
		if err := proxy.Do(ctx, "https://a.thumbs.redditmedia.com/"+ctx.OriginalURL()[8:], defClient); err != nil {
			return err
		}

		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	router.Get("/bthumb/*", func(ctx *fiber.Ctx) error {
		if err := proxy.Do(ctx, "https://b.thumbs.redditmedia.com/"+ctx.OriginalURL()[8:], defClient); err != nil {
			return err
		}

		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	router.Get("/external/*", func(ctx *fiber.Ctx) error {
		ctx.Request().Header.Set("Accept", "image/avif,image/webp,*/*")

		if err := proxy.Do(ctx, "https://external-preview.redd.it/"+ctx.OriginalURL()[10:], defClient); err != nil {
			return err
		}

		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	router.Get("/preview/*", func(ctx *fiber.Ctx) error {
		ctx.Request().Header.Set("Accept", "video/webm,video/ogg,video/*;q=0.9,application/ogg;q=0.7,audio/*;q=0.6,*/*;q=0.5")

		if err := proxy.Do(ctx, "https://preview.redd.it/"+ctx.OriginalURL()[9:], defClient); err != nil {
			return err
		}

		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	router.Get("/rstyle/*", func(ctx *fiber.Ctx) error {
		if err := proxy.Do(ctx, "https://styles.redditmedia.com/"+ctx.OriginalURL()[8:], defClient); err != nil {
			return err
		}

		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	router.Get("/rstatic/*", func(ctx *fiber.Ctx) error {
		if err := proxy.Do(ctx, "https://www.redditstatic.com/"+ctx.OriginalURL()[9:], defClient); err != nil {
			return err
		}

		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	router.Get("/imgur/*", func(ctx *fiber.Ctx) error {
		if err := proxy.Do(ctx, "https://i.imgur.com/"+ctx.OriginalURL()[7:], defClient); err != nil {
			return err
		}

		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	router.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("views/index", fiber.Map{})
	})

	router.Get("/config", func(ctx *fiber.Ctx) error {
		preferredres, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
		if err != nil {
			preferredres = 3
		}

		return ctx.Render("views/config", fiber.Map{
			"ResHelper":   resolutionHelper,
			"ThemeHelper": themeHelper,
			ResCookie:     preferredres,
		})
	})

	router.Post("/config", func(ctx *fiber.Ctx) error {
		for cookiekey, cookievalue := range cfgMap {
			switch formvalue := ctx.FormValue(cookiekey); formvalue {
			case "on":
				SetcfgCookie(ctx, cookievalue, "1")
			case "off":
				SetcfgCookie(ctx, cookievalue, "0")
			}
		}

		for i := range resolutionHelper {
			resStr := strconv.Itoa(resolutionHelper[i].Value)
			if ctx.FormValue("PrefRes") == resStr {
				SetcfgCookie(ctx, ResCookieValue, resStr)
				break
			}
		}

		for k := range themeHelper {
			if ctx.FormValue("PrefTheme") == themeHelper[k].Value {
				SetcfgCookie(ctx, ThemeCookieValue, themeHelper[k].Value)
				break
			}
		}

		return ctx.RedirectBack("/config", fiber.StatusMovedPermanently)
	})

	router.Get("/r/:sub", func(ctx *fiber.Ctx) error {
		after := ctx.Query("after")
		flair := url.QueryEscape(ctx.Query("f"))
		subname := utils.ToLower(ctx.Params("sub"))

		var posts types.Posts

		if postcache, exists := jsonCache.Load(fmt.Sprintf("%s-%s-%s", subname, after, flair)); exists {
			posts = postcache.(types.Posts) //nolint:errcheck,forcetypeassert // if this fails, we are fucked anyway.
		} else {
			posts = logic.GetPosts(subname, after, flair)
			jsonCache.Store(fmt.Sprintf("%s-%s-%s", subname, after, flair), posts)
		}

		if len(posts.Data.Children) == 0 {
			return ctx.Render("views/404", fiber.Map{})
		}

		// Cache subreddit data, so we don't have to keep making requests every single time.
		// This will store it in memory, which may not be the best, and a disk based cache would be better.
		var sub types.Subreddit

		if scache, exists := subCache.Load(subname); exists {
			sub = scache.(types.Subreddit) //nolint:errcheck,forcetypeassert // This should not error out.
		} else {
			sub = logic.GetSubredditData(subname)
			subCache.Store(subname, sub)
		}

		resolutionToUse, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
		if err != nil {
			resolutionToUse = 3
		}

		flairuesc, err := url.QueryUnescape(flair)
		if err != nil {
			log.Errorf("Failed to unescape flair: %w", err)
		}

		SortPostData(&posts, resolutionToUse)

		return ctx.Render("views/sub", fiber.Map{
			"SubName":       subname,
			"SubData":       sub.Data,
			"Posts":         posts.Data,
			"FlairFiltered": flairuesc,
		})
	})

	router.Get("/r/:sub/comments/:id/*", func(ctx *fiber.Ctx) error {
		subname := utils.ToLower(ctx.Params("sub"))
		cid := ctx.Params("id")

		post, comm := logic.GetComments(subname, cid)

		resolutionToUse, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
		if err != nil {
			resolutionToUse = 3
		}

		SortPostData(&post, resolutionToUse)

		return ctx.Render("views/comments", fiber.Map{
			"Posts":    post.Data,
			"Comments": comm.Data.Children,
		})
	})

	router.Get("/u/:user", func(ctx *fiber.Ctx) error {
		username := ctx.Params("user")
		after := ctx.Query("after")

		post := logic.GetAccount(username, after)

		resolutionToUse, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
		if err != nil {
			resolutionToUse = 3
		}

		SortPostData(&post, resolutionToUse)

		return ctx.Render("views/account", fiber.Map{
			"Posts":    post.Data,
			"username": username,
		})
	})

	router.Post("/loadPosts", func(ctx *fiber.Ctx) error {
		username := ctx.FormValue("user")
		after := ctx.FormValue("after")
		if username != "" {
			posts := logic.GetAccount(username, after)

			resolutionToUse, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
			if err != nil {
				resolutionToUse = 3
			}

			SortPostData(&posts, resolutionToUse)

			return ctx.Render("views/ucomm", fiber.Map{
				"username": username,
				"Posts":    posts.Data,
			})
		}

		subname := ctx.FormValue("sub")
		flair := url.QueryEscape(ctx.FormValue("flair"))

		var posts types.Posts

		if postcache, exists := jsonCache.Load(fmt.Sprintf("%s-%s-%s", subname, after, flair)); exists {
			posts = postcache.(types.Posts) //nolint:errcheck,forcetypeassert // if this fails, we are fucked anyway.
		} else {
			posts = logic.GetPosts(subname, after, flair)
			jsonCache.Store(fmt.Sprintf("%s-%s-%s", subname, after, flair), posts)
		}

		resolutionToUse, err := strconv.Atoi(ctx.Cookies(ResCookieValue))
		if err != nil {
			resolutionToUse = 3
		}

		flairuesc, err := url.QueryUnescape(flair)
		if err != nil {
			log.Errorf("Failed to unescape flair: %w", err)
		}

		SortPostData(&posts, resolutionToUse)

		return ctx.Render("views/posts", fiber.Map{
			"SubName":       subname,
			"Posts":         posts.Data,
			"FlairFiltered": flairuesc,
		})
	})

	// NoRoute
	router.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).Render("views/404", fiber.Map{})
	})

	// Shutting down
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan

		log.Infof("\n\n\nSession lasted %s\nSee you next time!\n\n\n", time.Since(startTime).String())

		if err := router.Shutdown(); err != nil {
			log.Errorf("Failed to gracefully Shutdown router: %w", err)
		}
	}()

	// localhost:9090
	log.Fatal(router.Listen(":9090")) //nolint:revive // No
}

func SortPostData(posts *types.Posts, resolutionToUse int) {
	origRes := resolutionToUse
	for i := range posts.Data.Children {
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					log.Warnf("Recovered from panic: %w", rec)
				}
			}()

			post := &posts.Data.Children[i].Data

			if len(post.Preview.Images) > 0 {
				image := post.Preview.Images[0]

				if len(image.Resolutions) > 0 && resolutionToUse != MaxResolution {
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
					if len(image.Variants.MP4.Resolutions) > 0 && resolutionToUse != MaxResolution {
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

	switch isVideo {
	case true:
		source = RewriteURL(mData.S.MP4)
	case false:
		if resolutionToUse == MaxResolution {
			source = RewriteURL(mData.S.U)
		} else {
			source = poster
		}
	}
	return types.InternalVData{
		Video:                   isVideo,
		Link:                    source,
		AutoChosenPosterQuality: poster,
	}
}

func SetcfgCookie(ctx *fiber.Ctx, cookiename, cookievalue string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     cookiename,
		Value:    cookievalue,
		Expires:  time.Now().Add((24 * time.Hour) * 365),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "lax",
	})
}

func backgroundJanitor() {
	timer := time.NewTimer(CleanTimer)
	for range timer.C {
		jsonCache = sync.Map{}
		log.Info("Scheduled maintenance: jsonCache emptied.")
		if !timer.Stop() {
			timer.Reset(CleanTimer)
		}
	}
}
