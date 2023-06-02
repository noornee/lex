//nolint:depguard // allow flag here.
package main

import (
	"flag"
	"log"

	"github.com/cmd777/lex/src/logic/version"
	"github.com/cmd777/lex/src/router"
)

func main() {
	checkUpdate := flag.Bool("checkupdate", true, "enables automatic update checking")
	flag.Parse()

	if *checkUpdate {
		if okf, cversion := version.CurrentVersion(); okf {
			log.Printf("Running LEX Version %d", cversion)
			log.Println("Checking for updates... (to disable update checking, run the application with -checkupdate=false)")
			if ok, latest := version.CheckForUpdates(); !ok {
				log.Println("There was an error while attempting to check for updates, try again later.")
			} else if cversion < latest {
				log.Printf("Your LEX version is outdated (version mismatch -> [gh:%d | local:%d])\r\n", latest, cversion)
			} else {
				log.Println("You are running the latest version of LEX")
			}
		} else {
			log.Println("Failed to read local version file. (is it missing?)")
		}
	} else {
		log.Printf("update checking disabled")
	}

	router.StartServer()
}
