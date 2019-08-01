// Copyright (C) 2018  denis4net
// Copyright (C) 2019 Evgeny Kuznetsov (evgeny@kuznetsov.md)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along tihe this program. If not, see <https://www.gnu.org/licenses/>.

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
	googleClientId     string
	googleClientSecret string
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
		googleClientId = cred.id
		googleClientSecret = cred.secret
	}

	if googleClientId == "" || googleClientSecret == "" {
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
