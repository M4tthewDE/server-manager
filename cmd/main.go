package main

import (
	"log"

	"github.com/m4tthewde/server-manager/internal"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.Println("Starting server manager")

	internal.Run()
}
