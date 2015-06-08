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

func TestProcessValue(t *testing.T) {
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
