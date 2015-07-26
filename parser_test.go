package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

var sut_parser parser

func init() {
	sut_parser = parser{}
	log.SetOutput(ioutil.Discard)
}

func TestDuplicates(t *testing.T) {
	contents := `[SECTION]
	dupvalue
	dupvalue
	`
	// capture output
	r, w, _ := os.Pipe()
	log.SetOutput(w)
	defer func() { log.SetOutput(os.Stdout) }()

	// do work
	sut_parser.processAllLines(strings.NewReader(contents))
	w.Close()

	// now check that message is displayed
	captured, _ := ioutil.ReadAll(r)

	actual := string(captured)
	expected := "duplicate entry"

	if !strings.Contains(actual, expected) {
		t.Errorf("Expected: %s, but was: %s", expected, actual)
	}
}

func TestNoDuplicates(t *testing.T) {
	contents := `[SECTIONA]
	nodupvalue
	[SECTIONB]
	nodupvalue`

	// capture output
	r, w, _ := os.Pipe()
	log.SetOutput(w)
	defer func() { log.SetOutput(os.Stdout) }()

	// do work
	sut_parser.processAllLines(strings.NewReader(contents))
	w.Close()

	// now check that message is displayed
	captured, _ := ioutil.ReadAll(r)

	actual := string(captured)
	if strings.Contains(actual, "duplicate entry") {
		t.Errorf("No duplicates expected.")
	}
}

func TestProcessLineValue(t *testing.T) {
	sut_parser.currentVariable = "TESTING"
	sut_parser.addToCurrentVariable("")

	sut_parser.processLine(" value # comment")
	assertEquals(t, "value", sut_parser.env[sut_parser.currentVariable][0])
}

func TestProcessLineSection(t *testing.T) {
	sut_parser.processLine("[ A SECTION ]")
	assertEquals(t, "ASECTION", sut_parser.currentVariable)
}

func assertEquals(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}
