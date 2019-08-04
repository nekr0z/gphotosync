package main

import (
	"testing"
	"time"
)

func TestAddTimestampToPath(t *testing.T) {
	testTime, _ := time.Parse("2006-01-02 15:04:05.000000000 MST -07:00", "1609-09-12 19:02:35.123456789 PDT +03:00")
	got := addTimestampToPath("somepath/whatever.tex", testTime)
	want := "somepath/whatever.tex-62359ca6994c1b15.tex"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
