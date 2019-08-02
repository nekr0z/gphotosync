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
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

type packFile struct {
	src string
	dst string
	mod os.FileMode
}

type target struct {
	goos   string
	goarch string
}

const appName = "gphotosync"

var (
	filename          string
	binName           string
	goos              string
	goarch            string
	keyID             string = "8BCDD592F2DE32D4"
	version           string
	btime             int64
	packFilesRequired = []packFile{
		{src: "LICENSE", dst: "LICENSE", mod: 0644},
		{src: "README.md", dst: "README.md", mod: 0644},
		{src: "SOURCE.txt", dst: "SOURCE.txt", mod: 0644},
	}
	packFiles = []packFile{}
	targets   = []target{
		{goos: "linux", goarch: "386"},
		{goos: "linux", goarch: "amd64"},
		{goos: "linux", goarch: "arm"},
		{goos: "windows", goarch: "386"},
		{goos: "windows", goarch: "amd64"},
		{goos: "darwin", goarch: "amd64"},
	}
)

func main() {
	sign := flag.Bool("s", false, "sign binary")
	tar := flag.Bool("t", false, "generate archive")
	buildAll := flag.Bool("a", false, "build archives for all supported systems")
	flag.StringVar(&goos, "goos", runtime.GOOS, "GOOS value")
	flag.StringVar(&goarch, "goarch", runtime.GOARCH, "GOARCH value")
	flag.Parse()
	version = getVersion()
	btime = buildTime()
	if *buildAll {
		*tar = true
	}
	if *tar {
		*sign = true
	}

	if *buildAll != true {
		targets = []target{
			{goos: goos, goarch: goarch},
		}
	}
	for _, tgt := range targets {
		packFiles = nil
		for _, pf := range packFilesRequired {
			packFiles = append(packFiles, pf)
		}
		build(tgt, *sign, *tar)
	}
}

func build(tgt target, sign bool, tar bool) {
	fmt.Printf("Building version %s for %s-%s\n", version, tgt.goos, tgt.goarch)
	fmt.Println("Building as of", time.Unix(btime, 0))

	if tgt.goos == "windows" {
		binName = appName + ".exe"
	} else {
		binName = appName
	}
	os.Setenv("GOOS", tgt.goos)
	os.Setenv("GOARCH", tgt.goarch)
	buildBinary(version, btime)

	if sign {
		signFile(binName, keyID)
	}

	if tar {
		filename = fmt.Sprintf("%s-%s-%s-%s", appName, tgt.goos, tgt.goarch, version)
		if tgt.goos == "windows" {
			buildZip()
		} else {
			buildTar()
		}
		fmt.Println("archive", filename, "created")
	}
}

func buildBinary(version string, t int64) {
	cred := Credentials{
		ID:     "",
		Secret: "",
	}
	err := ReadSecretJSON(".client_secret.json", &cred)
	if err != nil {
		fmt.Printf("couldn't read credentials from .client_secret.json: %s\n", err)
	}
	cmdline := fmt.Sprintf("go build -ldflags \"-X main.version=%s -X main.googleClientId=%s -X main.googleClientSecret=%s\" -v", version, cred.ID, cred.Secret)
	cmd := exec.Command("bash", "-c", cmdline)
	if err := cmd.Run(); err != nil {
		log.Fatalln("failed to build binary")
	}
	setFileTime(binName, t)
	packFiles = append(packFiles, packFile{binName, binName, 0755})
}

func buildZip() {
	filename = filename + ".zip"
	zf, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	for _, f := range packFiles {
		sf, err := os.Open(f.src)
		if err != nil {
			log.Fatal(err)
		}

		info, err := sf.Stat()
		if err != nil {
			log.Fatal(err)
		}

		h, err := zip.FileInfoHeader(info)
		if err != nil {
			log.Fatal(err)
		}

		h.Method = zip.Deflate

		w, err := zw.CreateHeader(h)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(w, sf)
		if err != nil {
			log.Fatal(err)
		}

		sf.Close()
	}
}

func buildTar() {
	for i := range packFiles {
		packFiles[i].dst = filename + "/" + packFiles[i].dst
	}
	filename = filename + ".tar.gz"
	fd, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	gw, err := gzip.NewWriterLevel(fd, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	tw := tar.NewWriter(gw)
	for _, f := range packFiles {
		sf, err := os.Open(f.src)
		if err != nil {
			log.Fatal(err)
		}
		info, err := sf.Stat()
		if err != nil {
			log.Fatal(err)
		}
		h := &tar.Header{
			Name:    f.dst,
			Size:    info.Size(),
			Mode:    int64(f.mod),
			ModTime: info.ModTime(),
		}
		err = tw.WriteHeader(h)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(tw, sf)
		if err != nil {
			log.Fatal(err)
		}
		sf.Close()
	}
	err = tw.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = gw.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = fd.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func setFileTime(f string, t int64) {
	cmd := exec.Command("touch", "-t", fmt.Sprint(time.Unix(t, 0).Format("200601021504.05")), f)
	if err := cmd.Run(); err != nil {
		log.Fatalln("failed to set time on", f)
	}
}

func signFile(f string, k string) {
	cmd := exec.Command("gpg", "--detach-sign", "--yes", "--passphrase", os.Getenv("GPG_PASSPHRASE"), "--pinentry-mode", "loopback", "-a", "-u", k, f)
	if err := cmd.Run(); err != nil {
		fmt.Println("signing", f, "failed")
		filename = filename + "-unsigned"
	} else {
		fmt.Println(f, "successfully signed with key", k)
		packFiles = append(packFiles, packFile{binName + ".asc", binName + ".asc", 0644})
	}
}

func getVersion() string {
	s, err := getString("git", "describe", "--always", "--dirty")
	versionRe := regexp.MustCompile(`^v[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}(\-[0-9]{1,3}\-g[0-9a-f]{5,15})?`)
	if err == nil {
		if versionRe.MatchString(s) {
			return s
		}
	}
	return "unknown"
}

func getString(c string, a ...string) (string, error) {
	cmd := exec.Command(c, a...)
	b, err := cmd.CombinedOutput()
	return string(bytes.TrimSpace(b)), err
}

func buildTime() int64 {
	s, err := getString("git", "show", "-s", "--format=%ct")
	if err == nil {
		if i, e := strconv.ParseInt(s, 10, 64); e == nil {
			return i
		}
	}
	return time.Now().Unix()
}
