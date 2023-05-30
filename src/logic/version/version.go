package version

import (
	"context"
	"embed"
	"io"
	"log"
	"net/http"
	"strconv"
)

//go:embed internalversion.txt
var staticFile embed.FS

func CurrentVersion() (_ bool, _ int) {
	intver, err := staticFile.ReadFile("internalversion.txt")
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

func CheckForUpdates() (_ bool, _ int) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://raw.githubusercontent.com/cmd777/lex/main/VERSION.txt", http.NoBody)
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
