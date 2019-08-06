package main

import (
	"testing"

	photoslibrary "github.com/nekr0z/gphotoslibrary"
)

var (
	testLib = Library{
		Path:         "/some/path/",
		Deduplicator: deduplicatePath,
	}
	testMetadata = photoslibrary.MediaMetadata{
		CreationTime: "2006-01-02T15:04:05Z",
	}
	testItem = photoslibrary.MediaItem{
		Filename:      "20060102_150405.mp4",
		MediaMetadata: &testMetadata,
	}
)

func TestDeduplicatePath(t *testing.T) {

	t.Run("-gphotosync-UnixNano.ext", func(t *testing.T) {
		got := deduplicatePath("somepath/whatever.tex", &testItem)
		want := "somepath/whatever-gphotosync-fc4a4d5fdf6b200.tex"
		assertCorrectMessage(t, got, want)
	})

	t.Run("-gphotosync-id.ext", func(t *testing.T) {
	})
}

func TestGetMediaPath(t *testing.T) {
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
