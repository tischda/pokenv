package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

var sut_pokenv pokenv

func TestProcessTestFile(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stdout)

	sut_pokenv = pokenv{registry: mock}
	sut_pokenv.processFile(REG_KEY_MACHINE, `data/setvar.txt`)
	expected := "valueline1"
	actual := mock.env["POKE_SECTION"]
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}

func TestCheckPath(t *testing.T) {
	var paths []string
	if runtime.GOOS == "windows" {
		paths = []string{
			`c:\Windows`,
			`c:\Windows\system32`,
			`%windir%`,
			`%windir%\system32`,
			`.`,
		}
	} else {
		// no variable expansion here
		paths = []string{
			`/etc`,
			`/usr/bin`,
			`/var`,
			`.`,
		}
	}
	for _, path := range paths {
		if isPathInvalid(path) {
			t.Errorf("Invalid path: %q", path)
		}
	}
}

// Inspired by https://talks.golang.org/2014/testing.slide#23
func TestParseAndCheckPaths(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		log.SetFlags(0)
		sut_pokenv = pokenv{registry: mock, checkPath: true}
		sut_pokenv.processFile(REG_KEY_MACHINE, `data/pathvars.txt`)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestParseAndCheckPaths")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	// capture output of process execution
	r, w, _ := os.Pipe()
	cmd.Stderr = w
	err := cmd.Run()
	w.Close()

	// check return code
	if e, ok := err.(*exec.ExitError); ok && e.Success() {
		t.Fatalf("Exptected exit status 1, but was: %v, ", err)
	}

	// now check that message is displayed
	captured, _ := io.ReadAll(r)
	actual := string(captured)
	expected := "Invalid path"

	if !strings.Contains(actual, expected) {
		t.Errorf("Expected: %s, but was: %s", expected, actual)
	}
}
