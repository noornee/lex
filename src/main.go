package main

import (
	"log"
	"main/logic/version"
	"main/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if !fiber.IsChild() {
		log.SetFlags(log.LstdFlags | log.Lshortfile)

		log.Printf("Running LEX Version %d", version.LEX_VERSION)
		log.Println("Checking for updates...")
		ok, latest := version.CheckForUpdates()

		if !ok {
			log.Println("There was an error while attempting to check for updates, try again later.")
		} else if version.LEX_VERSION < latest {
			log.Printf("Your LEX version is outdated (version mismatch -> [gh:%d | local:%d])\r\n", latest, version.LEX_VERSION)
		} else {
			log.Println("You are running the latest version of LEX")
		}
	}

	router.StartServer()
}
