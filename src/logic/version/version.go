package version

import (
	"context"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
)

//go:embed VERSION.txt
var staticVersion embed.FS

//go:embed updater_VERSION.txt
var updaterVersion embed.FS

const (
	versionPath = "https://raw.githubusercontent.com/cmd777/lex/main/src/logic/version/VERSION.txt"
	updaterPath = "https://raw.githubusercontent.com/cmd777/lex/main/src/logic/version/updater_VERSION.txt"
)

func CurrentVersion() (_ bool, _ int) {
	intver, err := staticVersion.ReadFile("VERSION.txt")
	if err != nil {
		log.Errorf("failed to read version file: %w", err)
		return false, 0
	}

	version, err := strconv.Atoi(string(intver))
	if err != nil {
		log.Errorf("failed to convert string to int: %w", err)
		return false, 0
	}

	return true, version
}

func CurrentUpdaterVersion() (_ bool, _ int) {
	intver, err := updaterVersion.ReadFile("updater_VERSION.txt")
	if err != nil {
		log.Errorf("failed to read updater version file: %w", err)
		return false, 0
	}

	version, err := strconv.Atoi(string(intver))
	if err != nil {
		log.Errorf("failed to convert string to int: %w", err)
		return false, 0
	}

	return true, version
}

func CheckForUpdates() (_ bool, _ int) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, versionPath, http.NoBody)
	if err != nil {
		log.Errorf("Failed to create New Request with Context: %w", err)
		return false, 0
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("DefaultClient failed to do the request: %w", err)
		return false, 0
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Errorf("Failed to close response body: %w", closeerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read response body: %w", err)
		return false, 0
	}

	version, err := strconv.Atoi(string(body))
	if err != nil {
		log.Errorf("failed to convert int to string: %w", err)
		return false, 0
	}

	return true, version
}

func CheckForUpdaterUpdates() (_ bool, _ int) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, updaterPath, http.NoBody)
	if err != nil {
		log.Errorf("Failed to create New Request with Context: %w", err)
		return false, 0
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("DefaultClient failed to do the request: %w", err)
		return false, 0
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Errorf("Failed to close response body: %w", closeerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read response body: %w", err)
		return false, 0
	}

	version, err := strconv.Atoi(string(body))
	if err != nil {
		log.Errorf("failed to convert int to string: %w", err)
		return false, 0
	}

	return true, version
}

// #nosec G204
func LaunchUpdater() bool {
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
			cmd = exec.Command(fmt.Sprintf("%s%cupdater-amd64.exe", wd, os.PathSeparator))
		case "386":
			wd, err := os.Getwd()
			if err != nil {
				log.Errorf("failed to get current working directory: %w", err)
				return false
			}
			cmd = exec.Command(fmt.Sprintf("%s%cupdater-i386.exe", wd, os.PathSeparator))
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
			cmd = exec.Command(fmt.Sprintf("%s%cupdater-amd64", wd, os.PathSeparator))
		case "386":
			wd, err := os.Getwd()
			if err != nil {
				log.Errorf("failed to get current working directory: %w", err)
				return false
			}
			cmd = exec.Command(fmt.Sprintf("%s%cupdater-i386", wd, os.PathSeparator))
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
			cmd = exec.Command(fmt.Sprintf("%s%cupdater-amd64", wd, os.PathSeparator))
		case "arm64":
			wd, err := os.Getwd()
			if err != nil {
				log.Errorf("failed to get current working directory: %w", err)
				return false
			}
			cmd = exec.Command(fmt.Sprintf("%s%cupdater-arm64", wd, os.PathSeparator))
		default:
			return false
		}
	default:
		log.Errorf("Unsupported OS: %s", runtime.GOOS)
		return false
	}

	if err := cmd.Start(); err != nil {
		log.Errorf("failed to start updater: %w", err)
		return false
	}

	return true
}
