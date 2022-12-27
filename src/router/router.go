package router

import (
	"fmt"
	"html/template"
	"main/logic"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

func StartServer() {
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.BestCompression))

	// region Load Files

	router.SetFuncMap(template.FuncMap{
		"contains": strings.Contains,
		"notcontains": func(input, of string) bool {
			return !strings.Contains(input, of)
		},
		"sanitize": func(input string) template.HTML {
			Markdown := blackfriday.Run([]byte(input))
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
			input *= 100
			return fmt.Sprintf("%.0f", input)
		},
	})

	router.LoadHTMLGlob("views/*")

	router.GET("js/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.Header("Content-Type", "application/javascript")
		ctx.File(fmt.Sprintf("js/%v", id))
	})

	router.GET("css/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.Header("Content-Type", "text/css")
		ctx.File(fmt.Sprintf("css/%v", id))
	})

	// endregion

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
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

		for i := 0; i < len(Posts.Data.Children); i++ {
			if len(Posts.Data.Children[i].Data.Preview.Images) != 0 {
				Posts.Data.Children[i].Data.Preview.AutoChosenImageQuality = Posts.Data.Children[i].Data.Preview.Images[0].Resolutions[int(math.Round(float64(len(Posts.Data.Children[i].Data.Preview.Images[0].Resolutions)/2)))].URL
			}

			if len(Posts.Data.Children[i].Data.SecureMedia.RedditVideo.FallbackURL) != 0 {
				Posts.Data.Children[i].Data.SecureMedia.RedditVideo.LQ = fmt.Sprintf("%v/DASH_360.mp4", Posts.Data.Children[i].Data.LinkURL)
				Posts.Data.Children[i].Data.SecureMedia.RedditVideo.MQ = fmt.Sprintf("%v/DASH_480.mp4", Posts.Data.Children[i].Data.LinkURL)
				Posts.Data.Children[i].Data.SecureMedia.RedditVideo.Audio = fmt.Sprintf("%v/DASH_audio.mp4", Posts.Data.Children[i].Data.LinkURL)
			}

			if len(Posts.Data.Children[i].Data.MediaMetaData) != 0 {

				MMD := make(map[string]string)

				for n := range Posts.Data.Children[i].Data.MediaMetaData {
					MMD[n] = Posts.Data.Children[i].Data.MediaMetaData[n].P[int(math.Round(float64(len(Posts.Data.Children[i].Data.MediaMetaData[n].P)/2)))].U
				}

				Posts.Data.Children[i].Data.VMediaMetaData = MMD
			}

			if len(Posts.Data.Children[i].Data.SelfText) != 0 {
				// invisible character, blackfriday doesn't recognize it, and just displays &#x200B; which is pretty distracting in some cases.
				Posts.Data.Children[i].Data.SelfText = strings.Replace(Posts.Data.Children[i].Data.SelfText, "&amp;#x200B;", "", -1)
			}
		}

		ctx.HTML(http.StatusOK, "sub.html", gin.H{
			"SubData": Sub.Data,
			"Posts":   Posts.Data,
		})
	})

	router.GET("/r/:sub/loadPosts", func(ctx *gin.Context) {
		subname := url.QueryEscape(ctx.Param("sub"))
		after := url.QueryEscape(ctx.Query("after"))
		sort := url.QueryEscape(ctx.Query("t"))

		Posts := logic.GetPosts(after, sort, subname)

		for i := 0; i < len(Posts.Data.Children); i++ {
			if len(Posts.Data.Children[i].Data.Preview.Images) != 0 {
				Posts.Data.Children[i].Data.Preview.AutoChosenImageQuality = Posts.Data.Children[i].Data.Preview.Images[0].Resolutions[int(math.Round(float64(len(Posts.Data.Children[i].Data.Preview.Images[0].Resolutions)/2)))].URL
			}

			if len(Posts.Data.Children[i].Data.SecureMedia.RedditVideo.FallbackURL) != 0 {
				Posts.Data.Children[i].Data.SecureMedia.RedditVideo.LQ = fmt.Sprintf("%v/DASH_360.mp4", Posts.Data.Children[i].Data.LinkURL)
				Posts.Data.Children[i].Data.SecureMedia.RedditVideo.MQ = fmt.Sprintf("%v/DASH_480.mp4", Posts.Data.Children[i].Data.LinkURL)
				Posts.Data.Children[i].Data.SecureMedia.RedditVideo.Audio = fmt.Sprintf("%v/DASH_audio.mp4", Posts.Data.Children[i].Data.LinkURL)
			}

			if len(Posts.Data.Children[i].Data.MediaMetaData) != 0 {

				MMD := make(map[string]string)

				for n := range Posts.Data.Children[i].Data.MediaMetaData {
					MMD[n] = Posts.Data.Children[i].Data.MediaMetaData[n].P[int(math.Round(float64(len(Posts.Data.Children[i].Data.MediaMetaData[n].P)/2)))].U
				}

				Posts.Data.Children[i].Data.VMediaMetaData = MMD
			}

			if len(Posts.Data.Children[i].Data.SelfText) != 0 {
				// invisible character, blackfriday doesn't recognize it, and just displays &#x200B; which is pretty distracting in some cases.
				Posts.Data.Children[i].Data.SelfText = strings.Replace(Posts.Data.Children[i].Data.SelfText, "&amp;#x200B;", "", -1)
			}
		}

		ctx.HTML(http.StatusOK, "loadedposts.html", gin.H{
			"Posts": Posts.Data,
		})
	})

	// localhost:9090
	router.Run(":9090")
}
