package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	CSSMINURL = "https://www.toptal.com/developers/cssminifier/api/raw"
	JSMINURL  = "https://www.toptal.com/developers/javascript-minifier/api/raw"
)

func main() {
	MinifyCSS()
	MinifyJS()
}

func MinifyCSS() {
	fsdir, err := os.ReadDir("../css")
	if err != nil {
		log.Println(err)
	}

	for _, item := range fsdir {
		file, err := os.ReadFile(fmt.Sprintf("../css/%v", item.Name()))
		if err != nil {
			log.Println(err)
		}

		resp, err := http.PostForm(CSSMINURL, url.Values{
			"input": []string{string(file)},
		})
		if err != nil {
			log.Println(err)
		}

		defer resp.Body.Close()

		mincss, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		errw := os.WriteFile(fmt.Sprintf("../css/%v.min.css", strings.Split(item.Name(), ".")[0]), mincss, 0666)
		if errw != nil {
			log.Println(errw)
		}

		log.Printf("%v was successfully minimized.\r\n", item.Name())
	}
}

func MinifyJS() {
	fsdir, err := os.ReadDir("../js")
	if err != nil {
		log.Println(err)
	}

	for _, item := range fsdir {
		file, err := os.ReadFile(fmt.Sprintf("../js/%v", item.Name()))
		if err != nil {
			log.Println(err)
		}

		resp, err := http.PostForm(JSMINURL, url.Values{
			"input": []string{string(file)},
		})
		if err != nil {
			log.Println(err)
		}

		defer resp.Body.Close()

		minjs, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		errw := os.WriteFile(fmt.Sprintf("../js/%v.min.js", strings.Split(item.Name(), ".")[0]), minjs, 0666)
		if errw != nil {
			log.Println(errw)
		}

		log.Printf("%v was successfully minimized.\r\n", item.Name())
	}
}
