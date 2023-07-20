package version

import (
	"archive/zip"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

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

func currentAppVersion() (_ bool, _ int) {
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

func currentUpdaterVersion() (_ bool, _ int) {
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

func checkForAppUpdates() (_ bool, _ int) {
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

func checkForUpdaterUpdates() (_ bool, _ int) {
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
func launchUpdater() bool {
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

func CheckForUpdates() {
	if okf, cversion := currentAppVersion(); okf {
		log.Infof("Running LEX Version %d", cversion)
		log.Info("Checking for updates... (to disable update checking, run the application with -checkupdate=false)")

		if ok, latest := checkForAppUpdates(); !ok {
			log.Error("There was an error while attempting to check for updates, try again later.")
		} else if cversion < latest {
			log.Warnf("Your LEX version is outdated (version mismatch -> [gh:%d | local:%d])", latest, cversion)

			log.Info("Checking for updater updates...")

			if oku, cuversion := currentUpdaterVersion(); oku {
				log.Infof("Updater version is %d", cuversion)

				if okuc, nuversion := checkForUpdaterUpdates(); !okuc {
					log.Error("There was an error while attempting to check for updater updates, do you wish to use the old one? [Y/n]")
					if inputHandler() {
						if launchUpdater() {
							os.Exit(0)
						}
					}
				} else if cuversion < nuversion {
					log.Warnf("Updater is outdated (version mismatch -> [gh:%d | local:%d])", nuversion, cuversion)
					log.Warn("Do you wish to install the new updater now? [Y/n]")
					UpdateUpdater()
				} else {
					log.Info("Updater is up-to-date, do you wish to update LEX now? [Y/n]")
					if inputHandler() {
						if launchUpdater() {
							os.Exit(0)
						}
					}
				}
			}
		} else {
			log.Info("You are running the latest version of LEX")
		}
	}
}

func inputHandler() bool {
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
			return true
		} else {
			return false
		}
	case <-timer.C:
		log.Warn("Timeout exceeded, not updating")
		return false
	}
}

func Unzip(src, dst string, skipfile ...string) error {
	log.Infof("Reading zip file %s", src)
	archive, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to read zip file: %w", err)
	}
	defer func() {
		if err := archive.Close(); err != nil {
			log.Errorf("failed to close zip file: %w", err)
		}
	}()

	for _, file := range archive.File {
		if skipfile[0] != "" && strings.Contains(file.Name, skipfile[0]) {
			continue
		}

		log.Infof("Opening file %s", file.Name)
		reader, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		destpath, err := checkInvalidPath(dst, file.Name)
		if err != nil {
			return err
		}

		destpath = filepath.Clean(destpath)

		if file.FileInfo().IsDir() {
			log.Infof("Creating directory %s", destpath)
			if err := os.Mkdir(destpath, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		} else {
			log.Infof("Creating file %s", destpath)
			destfile, err := os.OpenFile(destpath, 0x00001|0x00040|0x00200, file.Mode()) // os.O_WRONLY|os.O_CREATE|os.O_TRUNC
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			log.Infof("Copying %d bytes from reader to %s", file.FileInfo().Size(), destfile.Name())
			if _, err := io.CopyN(destfile, reader, file.FileInfo().Size()); err != nil {
				return fmt.Errorf("failed to copy file data: %w", err)
			}
			log.Infof("File %s was created without any errors", destfile.Name())

			if err := destfile.Close(); err != nil {
				return fmt.Errorf("failed to close created file: %w", err)
			}
		}

		if err := reader.Close(); err != nil {
			return fmt.Errorf("failed to close reader: %w", err)
		}
	}

	return nil
}

func checkInvalidPath(k, v string) (string, error) {
	destpath := filepath.Join(k, v)
	if !strings.HasPrefix(destpath, filepath.Clean(k)) {
		return "", errors.New("illegal filepath")
	}
	return destpath, nil
}

func UpdateUpdater() {
	targetfile := fmt.Sprintf("updater-%s.zip", runtime.GOOS)
	targeturl := fmt.Sprintf("https://github.com/cmd777/lex/releases/download/updater/%s", targetfile)

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
	if err := Unzip(targetfile, currentdir); err != nil {
		log.Errorf("Failed to unzip file: %w", err)
	}

	log.Infof("Unzipped %s without any errors", targetfile)

	log.Infof("Removing %s from %s", targetfile, currentdir)

	if err := os.Remove(targetfile); err != nil {
		log.Errorf("failed to remove %s: %w", targetfile, err)
	}

	log.Info("Updater was updated successfully, starting the Updater and exiting LEX in 60 seconds.")

	timer := time.NewTimer(1 * time.Minute)

	for range timer.C {
		timer.Stop()
		if launchUpdater() {
			os.Exit(0)
		}
	}
}
