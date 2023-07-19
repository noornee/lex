package version

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

func StartUpdate() {
	log.Info("A new version was detected, do you wish to install it now? [Y/n]")

	inputchan := make(chan string, 1)
	timer := time.NewTimer(1 * time.Minute)

	go func() {
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			log.Errorf("Failed to scan input: %w", err)
		}
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

	log.Info("Getting current working directory")
	currentdir, err := os.Getwd()
	if err != nil {
		log.Errorf("failed to get the current working directory: %w", err)
	}

	log.Infof("Unzipping %s to %s", targetfile, currentdir)
	if err := Unzip(targetfile, currentdir); err != nil {
		log.Errorf("Failed to unzip file: %w", err)
	}

	log.Infof("Unzipped %s without any errors", targetfile)
}

func Unzip(src, dst string) error {
	log.Infof("Reading zip file %s", src)
	archive, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to read zip file: %w", err)
	}
	defer archive.Close()

	for _, file := range archive.File {
		log.Infof("Opening file %s", file.Name)
		reader, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer reader.Close()

		destpath := filepath.Join(dst, file.Name)

		if file.FileInfo().IsDir() {
			log.Infof("Creating directory %s", destpath)
			if err := os.Mkdir(destpath, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		} else {
			log.Infof("Creating file %s", destpath)
			destfile, err := os.OpenFile(destpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			defer destfile.Close()

			log.Infof("Copying data from reader to %s", destfile.Name())
			if _, err := io.Copy(destfile, reader); err != nil {
				return fmt.Errorf("failed to copy file data: %w", err)
			}
			log.Infof("File %s was created without any errors", destfile.Name())
		}
	}

	return nil
}
