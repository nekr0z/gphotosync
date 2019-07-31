package main

import (
	"flag"
	"log"
)

var (
	GoogleClientId     string
	GoogleClientSecret string
	version            string = "custom-build"
)

func main() {
	localLibArg := flag.String("lib", "", "local library path")
	flag.Parse()

	lib := Library{*localLibArg}

	if err := lib.Sync(); err != nil {
		log.Fatal(err)
	}
}
