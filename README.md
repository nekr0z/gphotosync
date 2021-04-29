# gphotosync
An app to download pictures and videos stored in your Google Photos.

[![Build Status](https://travis-ci.org/nekr0z/gphotosync.svg?branch=master)](https://travis-ci.org/nekr0z/gphotosync) [![codecov](https://codecov.io/gh/nekr0z/gphotosync/branch/master/graph/badge.svg)](https://codecov.io/gh/nekr0z/gphotosync) [![Go Report Card](https://goreportcard.com/badge/github.com/nekr0z/gphotosync)](https://goreportcard.com/report/github.com/nekr0z/gphotosync)

This app is a fork of [Denis Vashchuk's project](https://gitlab.com/denis4net/gphotosync).

##### Table of Contents
* [How it works](#how-it-works)
  * [Usage](#usage)
* [Authentification](#authentification)
* [Building the app](#building-the-app)
* [Privacy considerations](#privacy-considerations)
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
|  |  +-- VIDEO-gphotosync-159b90e8a152aa00.MPG
|  |  +-- VIDEO-gphotosync-159e8e0cdc16be00.MPG
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
The appended value is by default a hex representation of Unix timestamp. You can use `-strategy id` command line option to append Google Photo ID of the media file instead of timestamp.

### Usage
If you are planning to use the app for scheduled backups (which is totally OK), please consider the following:
* It's not a good idea to schedule these backups to run on the hour (i.e. at 3:00, 10:00, 17:00, etc.), and it's even worse to schedule them for midnight. If enough people do that, Google servers will have to process a lot of requests from all of us at the same very second; eventially Google will ban the API key. The app already has a random delay (between 0 and 60 seconds) before it starts, but you can help things more if you schedule the runs for some random minute of the hour.
* If the limit for Google API requests per day is reached, the app will pause for a minute before retrying to request more files, then for 2 minutes, then for 4 minutes, and so on, and so on. If a lot of users use the same pre-compiled project key, the delays can get significant, so it's a good idea to use some kind of blocking so that your new backup doesn't start before the previous one is finished. Linux users are advised to use `setlock` in their `cron` jobs.

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

## Privacy considerations
If you're using the precompiled version of the app with builtin project credentials, the activity of your copy of the app will be included in what I see in my Google Developer Console (it shows the cumulative activity such as number of requests per hour or day, median delay in serving those requests, and so on; I'm not aware of any way to see any personal or personalized data there). Other than that no one (but Google and your ISP or VPN provider, of course) has any way to have any access to your stuff or collect any data about how you are using it. Hey, you can read the source code and see for yourself what it does, and you can recompile the code yourself just to make sure!

## Credits
This software includes the following software or parts thereof:
* [googleapi](https://google.golang.org/api/googleapi) Copyright 2018 Google
* [nekr0z/gphotoslibrary](https://github.com/nekr0z/gphotoslibrary) Copyright 2018 Google, parts copyright 2019 Evgeny Kuznetsov
* [OAuth2 for Go](https://github.com/golang/oauth2) Copyright 2014 The Go Authors
* [The Go Programming Language](https://golang.org) Copyright 2009 The Go Authors
