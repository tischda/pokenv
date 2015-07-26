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
	t.Skip("not implemented")
	sut_pokenv.importFromFile(PATH_MACHINE, `data/setvar.txt`)
	assertEquals(t, "valueline1", mock.env["POKE_SECTION"])
}

func TestCheckPath(t *testing.T) {
	t.Skip("not implemented")
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
