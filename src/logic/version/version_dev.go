//go:build !release
// +build !release

package version

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func CurrentVersion() (ok bool, current int) {
	intver, err := os.ReadFile("./logic/version/internalversion.txt")
	if err != nil {
		log.Println(err)
		return false, 0
	}

	version, err := strconv.Atoi(string(intver))
	if err != nil {
		log.Println(err)
		return false, 0
	}

	return true, version
}

func CheckForUpdates() (ok bool, latest int) {
	resp, err := http.Get("https://raw.githubusercontent.com/cmd777/lex/main/VERSION.txt")
	if err != nil {
		log.Println(err)
		return false, 0
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false, 0
	}

	version, err := strconv.Atoi(string(body))
	if err != nil {
		log.Println(err)
		return false, 0
	}

	return true, version
}
