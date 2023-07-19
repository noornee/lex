package version

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

func StartUpdate() {
	log.Info("A new version was detected, do you wish to install it now? [Y/n]")

	inputchan := make(chan string, 1)
	timer := time.NewTimer(10 * time.Second)

	go func() {
		var input string
		fmt.Scanln(&input)
		inputchan <- input
	}()

	select {
	case in := <-inputchan:
		if strings.HasPrefix(strings.ToLower(in), "y") {
			log.Info("Starting update")
			UpdatePrep()
		} else {
			log.Info("Canceling update")
		}
		return
	case <-timer.C:
		log.Warn("Timeout exceeded, not updating")
		return
	}
}

func UpdatePrep() {
	targetfile := fmt.Sprintf("lex-%s.zip", runtime.GOOS)
	targeturl := fmt.Sprintf("https://github.com/cmd777/lex/releases/download/snapshot/%s", targetfile)

	log.Infof("Creating request to %s", targeturl)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, targeturl, http.NoBody)
	if err != nil {
		log.Errorf("Failed to create New Request with Context: %w", err)
		return
	}

	log.Info("Asking DefaultClient to do the request")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("DefaultClient failed to do the request: %w", err)
		return
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Errorf("Failed to close response body: %w", closeerr)
		}
	}()

	log.Infof("Reading body (%d bytes)", resp.ContentLength)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read response body: %w", err)
		return
	}

	log.Infof("Writing %d bytes to %s", len(body), targetfile)
	if err := os.WriteFile(targetfile, body, 0o600); err != nil {
		log.Errorf("failed to write to file: %w", err)
		return
	}

	// unzip targetfile
}
