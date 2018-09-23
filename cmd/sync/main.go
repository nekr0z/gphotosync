package main

import (
	"flag"
	"log"
)

var GoogleClientId string
var GoogleClientSecret string

func main() {
	localLibArg := flag.String("lib", "", "local library path")
	flag.Parse()

	lib := Library{*localLibArg}

	if err := lib.Sync(); err != nil {
		log.Fatal(err)
	}
}
