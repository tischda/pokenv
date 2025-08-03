//go:build windows
// +build windows

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	os.Args = append(os.Args, "-version")

	// capture output of process execution
	old := os.Stdout
	defer func() { os.Stdout = old }()
	r, w, _ := os.Pipe()
	os.Stdout = w
	main()
	w.Close()

	// now check that version is displayed
	captured, _ := io.ReadAll(r)
	actual := string(captured)
	expected := fmt.Sprintf("%s version %s\n", PROG_NAME,version)

	if !strings.Contains(actual, expected) {
		t.Errorf("Expected: %s, but was: %s", expected, actual)
	}
}
