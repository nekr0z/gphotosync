package main

import (
	"testing"

	photoslibrary "evgenykuznetsov.org/go/gphotoslibrary"
)

var (
	testLib = Library{
		Path:         "/some/path/",
		Deduplicator: dedupUnixHex,
	}
	testMetadata = photoslibrary.MediaMetadata{
		CreationTime: "2006-01-02T15:04:05Z",
	}
	testItem = photoslibrary.MediaItem{
		Filename:      "20060102_150405.mp4",
		MediaMetadata: &testMetadata,
		Id:            "AAJ-kdRYAOoSLowCcySWIOXkQzV_HTy78NjW9Sfq5OLf596iz09YdIpL4vO3KVW1uMJ7zhtpHWf7KcpAzudOyHjgiZNgiPRGuQ",
	}
)

func TestDeduplicatePath(t *testing.T) {

	t.Run("-gphotosync-UnixHex.ext", func(t *testing.T) {
		got, _ := deduplicatePath(&testLib, &testItem)
		want := "/some/path/2006/01/20060102_150405-gphotosync-fc4a4d5fdf6b200.mp4"
		assertCorrectMessage(t, got, want)
	})

	t.Run("-gphotosync-id.ext", func(t *testing.T) {
		testLib.Deduplicator = dedupID
		got, _ := deduplicatePath(&testLib, &testItem)
		want := "/some/path/2006/01/20060102_150405-gphotosync-AAJ-kdRYAOoSLowCcySWIOXkQzV_HTy78NjW9Sfq5OLf596iz09YdIpL4vO3KVW1uMJ7zhtpHWf7KcpAzudOyHjgiZNgiPRGuQ.mp4"
		assertCorrectMessage(t, got, want)
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
