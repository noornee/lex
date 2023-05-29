package router_test

import (
	"strings"
	"testing"
	"time"
	"unicode"

	. "main/router"
)

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
		if RewriteURL(urls[i]) != expect[i] {
			t.Fatalf("Test_RewriteURL: expected: %s, got: %s (#%d)", expect[i], RewriteURL(urls[i]), i)
		}
	}
}

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
		if trimexpect != trimsanitize {
			t.Fatalf("Test_Sanitize: expected: %s, got: %s (#%d)", trimexpect, trimsanitize, i)
		}
	}
}

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
		if QualifiesAsImg(test[i]) != expect[i] {
			t.Fatalf("Test_QualifiesAsImg: expected: %t, got: %t (#%d)", expect[i], QualifiesAsImg(test[i]), i)
		}
	}
}

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
		if FmtEpochDate(test[i]) != expect[i] {
			t.Fatalf("Test_FmtEpochDate: expected %s, got: %s (#%d)", expect[i], FmtEpochDate(test[i]), i)
		}
	}
}

// I cannot believe I'm writing a test for this.
func Test_Incrementbyone(t *testing.T) {
	t.Parallel()

	test := []int{
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
		if Incrementbyone(test[i]) != expect[i] {
			t.Fatalf("Test_Incrementbyone: expected %d, got: %d (#%d)", expect[i], Incrementbyone(test[i]), i)
		}
	}
}

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
		if FmtHumanDate(test[i]) != expect[i] {
			t.Fatalf("Test_FmtHumanDate: expected %s, got: %s (#%d)", expect[i], FmtHumanDate(test[i]), i)
		}
	}
}

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
		if ToPercentage(test[i]) != expect[i] {
			t.Fatalf("Test_ToPercentage: expected %s, got: %s (#%d)", expect[i], ToPercentage(test[i]), i)
		}
	}
}

// todo: addvartoctx, sortpostdata
