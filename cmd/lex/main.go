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
		version.CheckForUpdates()
	} else {
		log.Warn("update checking disabled")
	}

	router.StartServer()
}
