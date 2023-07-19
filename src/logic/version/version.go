package version

import (
	"context"
	"embed"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
)

//go:embed internalversion.txt
var staticFile embed.FS

func CurrentVersion() (_ bool, _ int) {
	intver, err := staticFile.ReadFile("internalversion.txt")
	if err != nil {
		log.Errorf("failed to open internalversion file: %w", err)
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
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://raw.githubusercontent.com/cmd777/lex/main/VERSION.txt", http.NoBody)
	if err != nil {
		log.Errorf("Failed to create New Request with Context: %w", err)
		return false, 0
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("DefaultClient failed to do the request: %w", err)
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

func LaunchUpdater() bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			cmd = exec.Command("updater-amd64.exe")
		case "386":
			cmd = exec.Command("updater-i386.exe")
		default:
			return false
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			cmd = exec.Command("updater-amd64")
		case "386":
			cmd = exec.Command("updater-i386")
		default:
			return false
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			cmd = exec.Command("updater-amd64")
		case "arm64":
			cmd = exec.Command("updater-arm64")
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
