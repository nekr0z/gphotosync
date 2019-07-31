package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var GoogleClientId string
var GoogleClientSecret string

func main() {
	// need to sleep for random number of milliseconds so as not to overload Google API
	// will sleep between 0 and 1 minute
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(60000)
	fmt.Printf("Waiting for about %d seconds...\n", (n / 1000))
	time.Sleep(time.Duration(n) * time.Millisecond)
	fmt.Println("Ready to start now")

	localLibArg := flag.String("lib", "", "local library path")
	flag.Parse()

	lib := Library{*localLibArg}

	if err := lib.Sync(); err != nil {
		log.Fatal(err)
	}
}
