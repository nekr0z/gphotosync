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

// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var (
	filesToRelease = [...]string{
		"./*.tar.gz",
		"./*.zip",
	}
)

func main() {
	flag.Parse()
	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatalln("github token not available, can't work")
	}
	if os.Getenv("GITHUB_USER") == "" {
		log.Fatalln("github user not set, can't work")
	}

	version, err := getString("git", "describe")
	if err != nil {
		log.Fatalln(err)
	}

	// check version format
	re := regexp.MustCompile(`^v[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
	if !re.MatchString(version) {
		log.Fatalln("version", version, "doesn't make sense, giving up!")
	}

	log.Println("Trying to release version", version)
	process(version)
}

func process(version string) {
	// build
	fmt.Println("building...")
	cmd := exec.Command("go", "run", "build.go", "oauth.go", "-a")
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}

	// release packages
	var fileNames []string
	for _, glob := range filesToRelease {
		fn, err := filepath.Glob(glob)
		if err != nil {
			log.Fatalln(err)
		}
		fileNames = append(fileNames, fn...)
	}

	for _, fileName := range fileNames {
		args := []string{
			fileName,
			"release/",
		}
		cmd := exec.Command("cp", args...)
		if err := cmd.Run(); err != nil {
			fmt.Println("failed to copy", fileName)
		}
		fmt.Println(fileName, "copied successfully")
	}
}

func getString(c string, a ...string) (string, error) {
	cmd := exec.Command(c, a...)
	b, err := cmd.CombinedOutput()
	return string(bytes.TrimSpace(b)), err
}
