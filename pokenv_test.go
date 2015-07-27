package main

import (
	"io/ioutil"
	"log"
	"runtime"
	"testing"
)

var sut_pokenv pokenv

func init() {
	sut_pokenv = pokenv{
		registry: mock,
	}
	log.SetOutput(ioutil.Discard)
}

func TestProcessTestFile(t *testing.T) {
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
