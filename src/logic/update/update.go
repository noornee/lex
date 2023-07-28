/*
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
	log.Info("LEX was updated successfully")
}
*/

package update

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/cmd777/lex/src/logic/version"
	"github.com/gofiber/fiber/v2/log"
)

const (
	versionPath = "https://raw.githubusercontent.com/cmd777/lex/main/src/logic/version/version.go"
)

var SemRegex = regexp.MustCompile(`(?m)v(?P<Major>\d+)\.(?P<Minor>\d+)\.(?P<Patch>\d+)`)

func checkForAppUpdates() (_ bool, _ []string) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, versionPath, http.NoBody)
	if err != nil {
		log.Errorf("Failed to create New Request with Context: %w", err)
		return false, nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("DefaultClient failed to do the request: %w", err)
		return false, nil
	}

	defer func() {
		if closeerr := resp.Body.Close(); closeerr != nil {
			log.Errorf("Failed to close response body: %w", closeerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read response body: %w", err)
		return false, nil
	}

	return true, SemRegex.FindStringSubmatch(string(body))
}

func CheckForUpdates() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Warnf("Recovered from a fatal panic: %w", rec)
		}
	}()

	log.Infof("Running LEX %s", version.VERSION)
	log.Info("Checking for updates... (to disable update checking, run the application with -checkupdate=false)")

	if ok, submatches := checkForAppUpdates(); ok {
		oklocal, localversion := ReadSemVer(SemRegex.FindStringSubmatch(version.VERSION)[1:])
		if !oklocal {
			log.Info("failed to read local version semver.")
			return
		}

		okgit, gitversion := ReadSemVer(submatches[1:])
		if !okgit {
			log.Error("failed to read git version semver.")
			return
		}

		localmajor := localversion[0]
		localminor := localversion[1]
		localpatch := localversion[2]

		gitmajor := gitversion[0]
		gitminor := gitversion[1]
		gitpatch := gitversion[2]

		switch {
		case gitmajor > localmajor:
			log.Infof("There is a new major version available: %d (%s)", gitmajor, submatches[0])
		case gitminor > localminor:
			log.Infof("There is a new minor version available: %d (%s)", gitminor, submatches[0])
		case gitpatch > localpatch:
			log.Infof("There is a new patch available: %d (%s)", gitpatch, submatches[0])
		default:
			log.Infof("You are running the latest version of LEX. (%s)", version.VERSION)
		}
	} else {
		log.Error("There was an error while attempting to check for updates, try again later.")
	}
}

func ReadSemVer(ver []string) (_ bool, _ []int) {
	major, err := strconv.Atoi(ver[0])
	if err != nil {
		log.Errorf("failed to convert major string to int: %w", err)
		return false, nil
	}
	minor, err := strconv.Atoi(ver[1])
	if err != nil {
		log.Errorf("failed to convert minor string to int: %w", err)
		return false, nil
	}
	patch, err := strconv.Atoi(ver[2])
	if err != nil {
		log.Errorf("failed to convert patch string to int: %w", err)
		return false, nil
	}

	return true, []int{major, minor, patch}
}

/* func inputHandler() bool {
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
		}
		return false
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
		if len(skipfile) > 0 {
			if strings.Contains(file.Name, skipfile[0]) {
				continue
			}
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
*/
