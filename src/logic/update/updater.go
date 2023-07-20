package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/cmd777/lex/src/logic/version"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	UpdatePrep()
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

	log.Infof("Reading body (%d bytes)", resp.ContentLength)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read response body: %w", err)
		return
	}

	if closeerr := resp.Body.Close(); closeerr != nil {
		log.Errorf("Failed to close response body: %w", closeerr)
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
	if err := version.Unzip(targetfile, currentdir, "updater"); err != nil {
		log.Errorf("Failed to unzip file: %w", err)
	}

	log.Infof("Unzipped %s without any errors", targetfile)

	log.Infof("Removing %s from %s", targetfile, currentdir)

	if err := os.Remove(targetfile); err != nil {
		log.Errorf("failed to remove %s: %w", targetfile, err)
	}

	log.Info("LEX was updated successfully, starting LEX and exiting updater in 60 seconds.")

	timer := time.NewTimer(1 * time.Minute)

	for range timer.C {
		timer.Stop()
		if launchLEX() {
			os.Exit(0) //nolint:revive // We will exit here
		}
	}
}

// #nosec G204
func launchLEX() bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			wd, err := os.Getwd()
			if err != nil {
				log.Errorf("failed to get current working directory: %w", err)
				return false
			}
			cmd = exec.Command(fmt.Sprintf("%s%clex-amd64-windows.exe", wd, os.PathSeparator))
		case "386":
			wd, err := os.Getwd()
			if err != nil {
				log.Errorf("failed to get current working directory: %w", err)
				return false
			}
			cmd = exec.Command(fmt.Sprintf("%s%clex-i386-windows.exe", wd, os.PathSeparator))
		default:
			return false
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			wd, err := os.Getwd()
			if err != nil {
				log.Errorf("failed to get current working directory: %w", err)
				return false
			}
			cmd = exec.Command(fmt.Sprintf("%s%clex-amd64-linux", wd, os.PathSeparator))
		case "386":
			wd, err := os.Getwd()
			if err != nil {
				log.Errorf("failed to get current working directory: %w", err)
				return false
			}
			cmd = exec.Command(fmt.Sprintf("%s%clex-i386-linux", wd, os.PathSeparator))
		default:
			return false
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			wd, err := os.Getwd()
			if err != nil {
				log.Errorf("failed to get current working directory: %w", err)
				return false
			}
			cmd = exec.Command(fmt.Sprintf("%s%clex-amd64-darwin", wd, os.PathSeparator))
		case "arm64":
			wd, err := os.Getwd()
			if err != nil {
				log.Errorf("failed to get current working directory: %w", err)
				return false
			}
			cmd = exec.Command(fmt.Sprintf("%s%clex-arm64-darwin", wd, os.PathSeparator))
		default:
			return false
		}
	default:
		log.Errorf("Unsupported OS: %s", runtime.GOOS)
		return false
	}

	if err := cmd.Start(); err != nil {
		log.Errorf("failed to launch LEX: %w", err)
		return false
	}

	return true
}
