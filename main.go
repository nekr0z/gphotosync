package main

import (
	"flag"
	"fmt"
	"log"
	"path"
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

	// if .client_secret.json exists in local lib path, use those credentials
	cred := credentials{
		id:     "",
		secret: "",
	}
	err := readSecret(path.Join(*localLibArg, ".client_secret.json"), &cred)
	if err != nil {
		fmt.Printf("couldn't read credentials from .client_secret.json: %s\n", err)
	} else {
		fmt.Println("using custom credentials")
		GoogleClientId = cred.id
		GoogleClientSecret = cred.secret
	}

	if err := lib.Sync(); err != nil {
		log.Fatal(err)
	}
}
