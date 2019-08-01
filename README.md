# gphotosync
An app to download pictures and videos stored in your Google Photos.

This app is a fork of [Denis Vashchuk's project](https://gitlab.com/denis4net/gphotosync).

##### Table of Contents
* [How it works](#how-it-works)
* [Authentification](#authentification)
* [Building the app](#building-the-app)
* [Credits](#credits)

## How it works
`gphotosync` downloads all the content of your Google Photos library to a local directory (current directory by default, otherwise to the one specified with `-lib [path]` command line option. It creates `year/month/` directory structure. If two files with the same filename happen to fit into the same month, the most recent one will be saved under its original filename, while the rest will have timestamp appended, so the directory structure will look something like this:
```
+-- 2019/
|  +-- 07/
|  |  +-- 01.jpg
|  |  +-- 02.jpg
|  |  +-- 20190704.JPG
|  +-- 05/
|  |  +-- VIDEO.MPG
|  |  +-- VIDEO.MPG-159b90e8a152aa00.MPG
|  |  +-- VIDEO.MPG-159e8e0cdc16be00.MPG
|  +-- 01/
|     +-- IMG512.JPG
+-- 2018/
|  +-- 12/
|  |  +-- IMAGE01.JPG
|  +-- 11/
|     +-- 0001.jpg
|     +-- 0004.jpg
...etc
         
```

## Authentification
You can use your own project's credentials for authentification: create a project using Photos Library API in [Google Developers Console](https://console.developers.google.com), download JSON credentials file from ID page, rename it to `.client_secret.json` and put in the directory where your local photos library will be downloaded. If no `.client_secret.json` is found in the working directory, the credentials supplied at build time will be used.

## Building the app
```
$ ./build.sh
```
or, if you're running a non-proper-shell-capable OS (i.e. Windows)
```
go run build.go oauth.go
```
If you want authentification credentials compiled in, have a `.client_secret.json` in repository directory at compile time (see [Authentification](#authentification) section for details).

## Credits
This software includes the following software or parts thereof:
* [googleapi](https://google.golang.org/api/googleapi) Copyright 2018 Google
* [nekr0z/gphotoslibrary](https://github.com/nekr0z/gphotoslibrary) Copyright 2018 Google, parts copyright 2019 Evgeny Kuznetsov
* [OAuth2 for Go](https://github.com/golang/oauth2) Copyright 2014 The Go Authors
* [The Go Programming Language](https://golang.org) Copyright 2009 The Go Authors
