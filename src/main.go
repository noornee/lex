package main

import (
	"flag"
	"log"
	"main/logic/version"
	"main/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if !fiber.IsChild() {
		checkUpdate := flag.Bool("checkupdate", true, "enables automatic update checking")
		flag.Parse()

		if *checkUpdate {
			if okf, CVersion := version.CurrentVersion(); okf {
				log.Printf("Running LEX Version %d", CVersion)
				log.Println("Checking for updates... (to disable update checking, run the application with -checkupdate=false)")
				if ok, latest := version.CheckForUpdates(); !ok {
					log.Println("There was an error while attempting to check for updates, try again later.")
				} else if CVersion < latest {
					log.Printf("Your LEX version is outdated (version mismatch -> [gh:%d | local:%d])\r\n", latest, CVersion)
				} else {
					log.Println("You are running the latest version of LEX")
				}
			} else {
				log.Println("Failed to read local version file. (is it missing?)")
			}
		} else {
			log.Printf("update checking disabled")
		}
	}

	router.StartServer()
}
