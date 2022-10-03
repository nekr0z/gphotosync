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
// along with this program. If not, see <https://www.gnu.org/licenses/>.

//go:generate go run version_generate.go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"time"
)

var (
	googleClientId     string
	googleClientSecret string
	version            string = "custom-build"
)

const defaultRedirectURL = "http://127.0.0.1:8088"

func main() {
	localLibArg := flag.String("lib", "", "local library path")
	dedup := flag.String("strategy", "unixhex", "filename deduplication strategy: [unixhex|id]")
	flag.Parse()

	lib := Library{
		Path:         *localLibArg,
		Deduplicator: dedupUnixHex,
	}
	switch *dedup {
	case "id":
		lib.Deduplicator = dedupID
	case "unixhex":
		lib.Deduplicator = dedupUnixHex
	}

	fmt.Printf("gphotosync version %s\n", version)

	// if .client_secret.json exists in local lib path, use those credentials
	cred := credentials{
		ID:          googleClientId,
		Secret:      googleClientSecret,
		RedirectURL: defaultRedirectURL,
	}
	err := readSecretJSON(path.Join(*localLibArg, ".client_secret.json"), &cred)
	if err != nil {
		fmt.Printf("couldn't read credentials from .client_secret.json: %s\n", err)
	} else {
		fmt.Println("using custom credentials")
	}

	if cred.ID == "" || cred.Secret == "" {
		log.Fatal("no credentials available, can not continue")
	}

	// need to sleep for random number of milliseconds so as not to overload Google API
	// but will not do it if no token exists
	if _, err := os.Stat(lib.GetTokenPath()); !os.IsNotExist(err) {
		// will sleep between 0 and 1 minute
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(60000)
		fmt.Printf("Waiting for about %d seconds...\n", (n / 1000))
		time.Sleep(time.Duration(n) * time.Millisecond)
		fmt.Println("Ready to start now")
	}

	if err := lib.Sync(cred); err != nil {
		log.Fatal(err)
	}
}

type credentials struct {
	ID          string
	Secret      string
	RedirectURL string
}

// readSecretJSON reads the saved credentials from a JSON file
func readSecretJSON(file string, cr *credentials) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	var f map[string]interface{}
	err = json.Unmarshal(data, &f)
	if err != nil {
		return err
	}

	cred := f["web"].(map[string]interface{})

	cr.ID = cred["client_id"].(string)
	cr.Secret = cred["client_secret"].(string)
	uris := cred["redirect_uris"].([]interface{})
	if len(uris) == 1 {
		uri := uris[0].(string)
		cr.RedirectURL = uri
	}
	return nil
}
