package main

import (
	"io/ioutil"
	"runtime"
	"testing"
	"os"
	"os/exec"
	"strings"
	"log"
)

var sut_pokenv pokenv

func TestProcessTestFile(t *testing.T) {
	sut_pokenv = pokenv{registry: mock,	}
	sut_pokenv.importFromFile(PATH_MACHINE, `data/setvar.txt`)
	expected := "valueline1"
	actual := mock.env["POKE_SECTION"]
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}

func TestCheckPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		paths := []string{
			`c:\Windows`,
			`c:\Windows\system32`,
			`%windir%`,
			`%windir%\system32`,
			`.`,
		}
		for _, path := range paths {
			if isPathInvalid(path) {
				t.Errorf("Invalid path:", path)
			}
		}
	} else {
		t.Skip("Cannot test windows paths")
	}
}

// Inspired by https://talks.golang.org/2014/testing.slide#23
func TestParseAndCheckPaths(t *testing.T) {
	if runtime.GOOS == "windows" {

		if os.Getenv("BE_CRASHER") == "1" {
			log.SetFlags(0)
			sut_pokenv = pokenv{registry: mock, pathcheck: true, }
			sut_pokenv.importFromFile(PATH_MACHINE, `data/pathvars.txt`)
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
		captured, _ := ioutil.ReadAll(r)
		actual := string(captured)
		expected := "Invalid path"

		if !strings.Contains(actual, expected) {
			t.Errorf("Expected: %s, but was: %s", expected, actual)
		}
	} else {
		t.Skip("Cannot test windows paths")
	}
}

