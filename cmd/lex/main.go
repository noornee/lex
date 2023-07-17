//nolint:depguard // allow flag here.
package main

import (
	"flag"

	"github.com/cmd777/lex/src/logic/version"
	"github.com/cmd777/lex/src/router"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	checkUpdate := flag.Bool("checkupdate", true, "enables automatic update checking")
	flag.Parse()

	if *checkUpdate {
		if okf, cversion := version.CurrentVersion(); okf {
			log.Infof("Running LEX Version %d", cversion)
			log.Info("Checking for updates... (to disable update checking, run the application with -checkupdate=false)")
			if ok, latest := version.CheckForUpdates(); !ok {
				log.Error("There was an error while attempting to check for updates, try again later.")
			} else if cversion < latest {
				log.Warnf("Your LEX version is outdated (version mismatch -> [gh:%d | local:%d])\r\n", latest, cversion)
			} else {
				log.Info("You are running the latest version of LEX")
			}
		} else {
			log.Error("Failed to read local version file. (is it missing?)")
		}
	} else {
		log.Warn("update checking disabled")
	}

	router.StartServer()
}
