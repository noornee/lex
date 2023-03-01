package version

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	LEX_VERSION = 299761
)

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
