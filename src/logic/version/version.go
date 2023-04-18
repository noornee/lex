//go:build release
// +build release

package version

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func CurrentVersion() (ok bool, current int) {
	intver, err := os.ReadFile("internalversion.txt")
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

func CheckForUpdates() (bool, int) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://raw.githubusercontent.com/cmd777/lex/main/VERSION.txt", nil)
	if err != nil {
		log.Println(err)
		return false, 0
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Println("Failed to close response body", closeerr)
		}
	}()

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
