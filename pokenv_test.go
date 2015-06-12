package main

import (
	"io/ioutil"
	"log"
	"testing"
)

var sut pokenv

func init() {
	sut = pokenv{
		environment: make(map[string][]string),
		registry:    mock,
	}
	log.SetOutput(ioutil.Discard)
}

func TestProcessLineValue(t *testing.T) {
	sut.currentVariable = "TESTING"
	sut.addCurrent("")

	sut.processLine(" value # comment")
	assertEquals(t, "value", sut.environment[sut.currentVariable][0])
}

func TestProcessLineSection(t *testing.T) {
	sut.processLine("[ A SECTION ]")
	assertEquals(t, "ASECTION", sut.currentVariable)
}

func TestProcessTestFile(t *testing.T) {
	sut.importEnv(PATH_MACHINE, `data/setvar.txt`)
	assertEquals(t, "valueline1", mock.env["POKE_SECTION"])
}

func assertEquals(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}

func TestCheckPath(t *testing.T) {
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
}
