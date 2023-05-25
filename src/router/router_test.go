package router_test

import (
	"strings"
	"testing"
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
			t.Fatalf("Test_RewriteURL: expected: %s, got: %s", expect[i], RewriteURL(urls[i]))
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
			t.Fatalf("Test_Sanitize: expected: %s, got: %s", trimexpect, trimsanitize)
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
			t.Fatalf("Test_QualifiesAsImg: expected: %v, got: %v", expect[i], QualifiesAsImg(test[i]))
		}
	}
}

// todo: sortpostdata
