package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	registry = mock
	log.SetOutput(ioutil.Discard)
}

func TestProcessLineValue(t *testing.T) {
	environment = make(map[string][]string)
	currentVariable = "TESTING"
	addCurrent("")

	processLine(" value # comment")
	assertEquals(t, "value", environment[currentVariable][0])
}

func TestProcessLineSection(t *testing.T) {
	environment = make(map[string][]string)
	processLine("[ A SECTION ]")
	assertEquals(t, "ASECTION", currentVariable)
}

func TestProcessTestFile(t *testing.T) {
	setEnv(HKLM, `data/setvar.txt`)
	assertEquals(t, "valueline1", mock.env["POKE_SECTION"])
}

func assertEquals(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}
