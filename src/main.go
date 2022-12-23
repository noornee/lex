package main

import (
	"log"
	"main/router"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	router.StartServer()
}
