package router_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/goccy/go-json"

	"github.com/cmd777/lex/src/logic/types"
	. "github.com/cmd777/lex/src/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"
)

// go test -v -race -run ^Test_RewriteURL$
func Test_RewriteURL(t *testing.T) {
	t.Parallel()

	urls := []string{
		"https://v.redd.it/v.mp4",
		"https://i.redd.it/i.png",
		"https://a.thumbs.redditmedia.com/a.png",
		"https://b.thumbs.redditmedia.com/b.png",
		"https://external-preview.redd.it/e.png",
		"https://preview.redd.it/p.png",
		"https://styles.redditmedia.com/style.css",
		"https://www.redditstatic.com/s.png",
		"https://i.imgur.com/sucks.png",
	}

	expect := []string{
		"/video/v.mp4",
		"/image/i.png",
		"/athumb/a.png",
		"/bthumb/b.png",
		"/external/e.png",
		"/preview/p.png",
		"/rstyle/style.css",
		"/rstatic/s.png",
		"/imgur/sucks.png",
	}

	for i := range urls {
		utils.AssertEqual(t, expect[i], RewriteURL(urls[i]))
	}
}

// go test -v -race -run ^Test_UGIDGen$
func Test_UGIDGen(t *testing.T) {
	t.Parallel()

	str := UGIDGen()

	if unicode.IsNumber(rune(str[0])) {
		t.Fatal("Test_UGIDGen: first rune cannot be a number")
	}
	if len(str) < 28 {
		t.Fatalf("Test_UGIDGen: wrong length: %d, wanted a minimum of 28", len(str))
	}
}

// go test -v -race -run ^Test_Sanitize$
func Test_Sanitize(t *testing.T) {
	t.Parallel()
	markdowns := []string{
		"### Hello, world!",
		`		
- Where?
- There!
		`,
		"**bold**",
		"[Google.com](https://google.com)",
		"<a href='https://google.com'>Google.com</a>",

		"<script>alert('hi')</script>",
		"[Google.com](javascript:alert('hi'))",
		"![Image](javascript:alert('hi'))",
		"<div onmouseover=\"alert('hi')\">hi</div>",
	}

	expect := []string{
		"<h3>Hello, world!</h3>",
		`
<ul>
<li>Where?</li>
<li>There!</li>
</ul>
		`,
		"<p><strong>bold</strong></p>",
		`<p><a href="https://google.com" rel="nofollow noreferrer noopener" target="_blank">Google.com</a></p>`,
		`<p><a href="https://google.com" rel="nofollow noreferrer noopener" target="_blank">Google.com</a></p>`,

		"<p></p>",
		`<p><a title="hi">Google.com</a>)</p>`,
		`<p><img alt="Image" title="hi"/>)</p>`,
		"<p><div>hi</div></p>",
	}
	for i := range markdowns {
		trimexpect := strings.TrimSpace(expect[i])
		trimsanitize := strings.TrimSpace(string(Sanitize(markdowns[i])))
		utils.AssertEqual(t, trimexpect, trimsanitize)
	}
}

// go test -v -race -run ^Test_QualifiesAsImg$
func Test_QualifiesAsImg(t *testing.T) {
	t.Parallel()

	test := []string{
		"https://example.com/a.gif",
		"https://example.com/is/a/long/url/image.png",
		"some-random-thing",
		"https://example.com/i.jpg",
		"https://example.com/a.gif/b.webp/c.png/d.jpeg",
		"https://example.com/i.webp",
		"https://imgur.sucks/i.gifv",
	}

	expect := []bool{
		true,
		true,
		false,
		true,
		true,
		false,
		false,
	}

	for i := range test {
		utils.AssertEqual(t, expect[i], QualifiesAsImg(test[i]))
	}
}

// go test -v -race -run ^Test_FmtEpochDate$
func Test_FmtEpochDate(t *testing.T) {
	t.Parallel()

	test := []float64{
		1493188556,
		1316752335,
		1514735842,
		1494144633,
		1195313356,
		1498503614,
		1443411730,
	}

	expect := []string{
		"Created Apr 26, 2017",
		"Created Sep 23, 2011",
		"Created Dec 31, 2017",
		"Created May 07, 2017",
		"Created Nov 17, 2007",
		"Created Jun 26, 2017",
		"Created Sep 28, 2015",
	}

	for i := range test {
		utils.AssertEqual(t, expect[i], FmtEpochDate(test[i]))
	}
}

// I cannot believe I'm writing a test for this.
// go test -v -race -run ^Test_Incrementbyone$
func Test_Incrementbyone(t *testing.T) {
	t.Parallel()

	test := []int{
		0,
		1,
		15,
		18,
		29,
		48,
		64,
		79,
		94,
		100,
	}

	expect := []int{
		1,
		2,
		16,
		19,
		30,
		49,
		65,
		80,
		95,
		101,
	}

	for i := range test {
		utils.AssertEqual(t, expect[i], Incrementbyone(test[i]))
	}
}

// Ah shit, here we go again.
// go test -v -race -run ^Test_Incrementbyonefloat$
func Test_Incrementbyonefloat(t *testing.T) {
	t.Parallel()

	test := []float64{
		0,
		1,
		15,
		18,
		29,
		48,
		64,
		79,
		94,
		100,
	}

	expect := []float64{
		1,
		2,
		16,
		19,
		30,
		49,
		65,
		80,
		95,
		101,
	}

	for i := range test {
		utils.AssertEqual(t, expect[i], Incrementbyonefloat(test[i]))
	}
}

// go test -v -race -run ^Test_FmtHumanDate$
func Test_FmtHumanDate(t *testing.T) {
	t.Parallel()

	test := []float64{
		float64(time.Now().Unix()),
		float64(time.Now().Add(time.Minute * -30).Unix()),
		float64(time.Now().Add(time.Hour * -1).Unix()),
		float64(time.Now().Add(time.Hour * -24).Unix()),
		float64(time.Now().Add(time.Hour * -168).Unix()),
		float64(time.Now().Add(time.Hour * -840).Unix()),
		float64(time.Now().Add(time.Hour * -10080).Unix()),
	}

	expect := []string{
		"now",
		"30 minutes ago",
		"1 hour ago",
		"1 day ago",
		"1 week ago",
		"1 month ago",
		"1 year ago",
	}

	for i := range test {
		utils.AssertEqual(t, expect[i], FmtHumanDate(test[i]))
	}
}

// go test -v -race -run ^Test_ToPercentage$
func Test_ToPercentage(t *testing.T) {
	t.Parallel()

	test := []float64{
		0.53,
		0.27,
		0.36,
		0.98,
		0.72,
		1,
	}

	expect := []string{
		"53",
		"27",
		"36",
		"98",
		"72",
		"100",
	}

	for i := range test {
		utils.AssertEqual(t, expect[i], ToPercentage(test[i]))
	}
}

// go test -v -race -run ^Test_SetcfgCookie$
func Test_SetcfgCookie(t *testing.T) {
	t.Parallel()

	testrouter := fiber.New()

	testrouter.Get("/", func(ctx *fiber.Ctx) error {
		SetcfgCookie(ctx, "index", "1")
		return ctx.SendString("Cookie Set")
	})

	testrouter.Get("/test", func(ctx *fiber.Ctx) error {
		SetcfgCookie(ctx, "test", "2")
		return ctx.SendString("Cookie Set")
	})

	testrouter.Post("/test2", func(ctx *fiber.Ctx) error {
		SetcfgCookie(ctx, "test2", "3")
		return ctx.SendString("Cookie Set")
	})

	resp, err := testrouter.Test(httptest.NewRequest(fiber.MethodGet, "/", http.NoBody))
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal("Test_SetcfgCookie: failed to close response body")
		}
	}()
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, "index", resp.Cookies()[0].Name)
	utils.AssertEqual(t, "1", resp.Cookies()[0].Value)

	resp, err = testrouter.Test(httptest.NewRequest(fiber.MethodGet, "/test", http.NoBody))
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal("Test_SetcfgCookie: failed to close response body")
		}
	}()
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, "test", resp.Cookies()[0].Name)
	utils.AssertEqual(t, "2", resp.Cookies()[0].Value)

	resp, err = testrouter.Test(httptest.NewRequest(fiber.MethodPost, "/test2", http.NoBody))
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal("Test_SetcfgCookie: failed to close response body")
		}
	}()
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, "test2", resp.Cookies()[0].Name)
	utils.AssertEqual(t, "3", resp.Cookies()[0].Value)
}

func Test_SortPostData(t *testing.T) {
	t.Parallel()

	var posts types.Posts

	contents, err := os.ReadFile("./testdata/test.json")
	utils.AssertEqual(t, nil, err)

	err = json.Unmarshal(contents, &posts)
	utils.AssertEqual(t, nil, err)

	SortPostData(&posts, 3)

	testpost := posts.Data.Children[1]

	utils.AssertEqual(t, false, testpost.Data.Archived)
	utils.AssertEqual(t, "jerf", testpost.Data.Author)
	utils.AssertEqual(t, float64(0), testpost.Data.CommentCount)
	utils.AssertEqual(t, "/r/golang/comments/147zgzw/reddit_api_protest/", testpost.Data.Permalink)
	utils.AssertEqual(t, "147zgzw", testpost.Data.PostID)
	utils.AssertEqual(t, "r/golang", testpost.Data.SubNamePref)
	utils.AssertEqual(t, "Reddit API Protest", testpost.Data.Title)
	utils.AssertEqual(t, float64(279), testpost.Data.Ups)
	utils.AssertEqual(t, float64(0.84), testpost.Data.UpvoteRatio)
	utils.AssertEqual(t, "", testpost.Data.Body)
}
