package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"path"
	"time"
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

	if GoogleClientId == "" || GoogleClientSecret == "" {
		log.Fatal("no credentials available, can not continue")
	}

	// need to sleep for random number of milliseconds so as not to overload Google API
	// will sleep between 0 and 1 minute
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(60000)
	fmt.Printf("Waiting for about %d seconds...\n", (n / 1000))
	time.Sleep(time.Duration(n) * time.Millisecond)
	fmt.Println("Ready to start now")

	if err := lib.Sync(); err != nil {
		log.Fatal(err)
	}
}
