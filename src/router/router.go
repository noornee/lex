package router

import (
	"fmt"
	"html/template"
	"main/logic"
	"main/logic/types"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	CommonExtNoNSH = blackfriday.NoIntraEmphasis | blackfriday.Tables | blackfriday.FencedCode | blackfriday.Autolink | blackfriday.Strikethrough | blackfriday.HeadingIDs | blackfriday.BackslashLineBreak | blackfriday.DefinitionLists
)

func StartServer() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(
		gzip.Gzip(gzip.BestCompression),
	)

	// region Load Files

	router.SetFuncMap(template.FuncMap{
		"contains": strings.Contains,
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

	router.LoadHTMLGlob("views/*")

	router.GET("/js/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.Header("Content-Type", "application/javascript")
		ctx.File(fmt.Sprintf("js/%v", id))
	})

	router.GET("/css/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.Header("Content-Type", "text/css")
		ctx.File(fmt.Sprintf("css/%v", id))
	})

	// endregion

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/config", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Will be implemented soon™")
	})

	// dev -> will probably keep this.
	router.POST("/byecookies", func(ctx *gin.Context) {
		ctx.Header("Set-Cookie", "nsfw_allowed=0")
		ctx.Redirect(http.StatusMovedPermanently, ctx.Request.Referer())
	})

	// for now, it's only purpose is to set cookies to nsfw subreddits (expand later™)
	router.POST("/config", func(ctx *gin.Context) {
		ctx.Header("Set-Cookie", "nsfw_allowed=1")
		ctx.Redirect(http.StatusMovedPermanently, ctx.Request.Referer())
	})

	router.GET("/r/:sub", func(ctx *gin.Context) {
		after := url.QueryEscape(ctx.Query("after"))
		sort := url.QueryEscape(ctx.Query("t"))
		subname := url.QueryEscape(ctx.Param("sub"))

		Sub := logic.GetSubredditData(subname)
		Posts := logic.GetPosts(after, sort, subname)

		if len(Posts.Data.Children) == 0 {
			ctx.String(http.StatusNotFound, "The subreddit 'r/%v' was banned, or doesn't exist. (Did you make a typo - exceeded the rate limit?)", subname)
			return
		}

		SortPostData(&Posts)

		nsfwallowed, _ := ctx.Cookie("nsfw_allowed")

		ctx.HTML(http.StatusOK, "sub.html", gin.H{
			"SubData":     Sub.Data,
			"Posts":       Posts.Data,
			"NSFWAllowed": nsfwallowed == "1" || !Sub.Data.NSFW,
		})
	})

	router.GET("/r/:sub/loadPosts", func(ctx *gin.Context) {
		subname := url.QueryEscape(ctx.Param("sub"))
		after := url.QueryEscape(ctx.Query("after"))
		sort := url.QueryEscape(ctx.Query("t"))

		Posts := logic.GetPosts(after, sort, subname)

		SortPostData(&Posts)

		ctx.HTML(http.StatusOK, "posts.html", gin.H{
			"Posts": Posts.Data,
		})
	})

	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

	// localhost:9090
	router.Run(":9090")
}

func SortPostData(Posts *types.Posts) {
	var wg sync.WaitGroup
	for i := 0; i < len(Posts.Data.Children); i++ {
		wg.Add(1)

		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			Post := Posts.Data.Children[i].Data
			if Post.Preview.Images != nil {
				Image := Post.Preview.Images[0]

				if Image.Resolutions != nil {
					Post.Preview.AutoChosenImageQuality = Image.Resolutions[int(math.Ceil(float64(len(Image.Resolutions)/2)))].URL
					Post.Preview.AutoChosenPosterQuality = Post.Preview.AutoChosenImageQuality
				} else {
					Post.Preview.AutoChosenImageQuality = Image.Source.URL
					Post.Preview.AutoChosenPosterQuality = Post.Preview.AutoChosenImageQuality
				}

				if strings.Contains(Image.Source.URL, ".gif") {
					if Image.Variants.MP4.Resolutions != nil {
						Post.Preview.AutoChosenImageQuality = Image.Variants.MP4.Resolutions[int(math.Ceil(float64(len(Image.Variants.MP4.Resolutions)/2)))].URL
					} else {
						Post.Preview.AutoChosenImageQuality = Image.Variants.MP4.Source.URL
					}
				}
			}

			if Post.SecureMedia != nil && Post.SecureMedia.RedditVideo != nil {
				Post.SecureMedia.RedditVideo.LQ = fmt.Sprintf("%v/DASH_360.mp4", Post.LinkURL)
				Post.SecureMedia.RedditVideo.MQ = fmt.Sprintf("%v/DASH_480.mp4", Post.LinkURL)
				Post.SecureMedia.RedditVideo.Audio = fmt.Sprintf("%v/DASH_audio.mp4", Post.LinkURL)
			}

			if Post.MediaMetaData != nil {
				MMD := make(map[string]string)

				for n := range Post.MediaMetaData {
					if Post.MediaMetaData[n].P != nil {
						MMD[n] = Post.MediaMetaData[n].P[int(math.Ceil(float64(len(Post.MediaMetaData[n].P)/2)))].U
					}
				}

				Post.VMediaMetaData = MMD
			}

			if len(Post.SelfText) != 0 {
				// invisible character, blackfriday doesn't recognize it, and just displays &#x200B; which is pretty distracting in some cases.
				Post.SelfText = strings.Replace(Post.SelfText, "&amp;#x200B;", "", -1)
			}
			Posts.Data.Children[i].Data = Post
		}(i, &wg)
	}
	wg.Wait()
}
