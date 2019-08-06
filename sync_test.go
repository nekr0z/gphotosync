package main

import (
	"testing"
	"time"

	photoslibrary "github.com/nekr0z/gphotoslibrary"
)

func TestDeduplicatePath(t *testing.T) {

	t.Run("-gphotosync-UnixNano.ext", func(t *testing.T) {
		testTime, _ := time.Parse("2006-01-02 15:04:05.000000000 MST -07:00", "1609-09-12 19:02:35.123456789 PDT +03:00")
		got := deduplicatePath("somepath/whatever.tex", testTime)
		want := "somepath/whatever-gphotosync-62359ca6994c1b15.tex"
		assertCorrectMessage(t, got, want)
	})

	t.Run("-gphotosync-id.ext", func(t *testing.T) {
	})
}

func TestGetMediaPath(t *testing.T) {
	testLib := Library{Path: "/some/path/"}
	testMetadata := photoslibrary.MediaMetadata{
		CreationTime: "2006-01-02T15:04:05Z",
	}
	testItem := photoslibrary.MediaItem{
		Filename:      "20060102_150405.mp4",
		MediaMetadata: &testMetadata,
	}

	got, _ := getMediaPath(&testLib, &testItem)
	want := "/some/path/2006/01/20060102_150405.mp4"
	assertCorrectMessage(t, got, want)
}

func assertCorrectMessage(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
