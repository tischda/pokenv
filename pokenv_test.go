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

	expected := "value"
	actual := environment[currentVariable][0]

	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}

func TestProcessLineSection(t *testing.T) {
	environment = make(map[string][]string)

	processLine("[ A SECTION ]")

	expected := "ASECTION"
	actual := currentVariable

	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}

func TestProcessTestFile(t *testing.T) {
	setEnv(HKLM, `data/setvar.txt`)

	expected := "valueline1"
	actual := mock.env["POKE_SECTION"]

	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}
