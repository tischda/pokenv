package main

import (
	"io/ioutil"
	"log"
	"runtime"
	"testing"
	"strings"
	"os"
)

var sut pokenv

func init() {
	sut = pokenv{
		environment: make(map[string][]string),
		registry:    mock,
	}
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
	defer func() {log.SetOutput(os.Stdout)}()

	// do work
	sut.processAllLines(strings.NewReader(contents))
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
	defer func() {log.SetOutput(os.Stdout)}()

	// do work
	sut.processAllLines(strings.NewReader(contents))
	w.Close()

	// now check that message is displayed
	captured, _ := ioutil.ReadAll(r)

	actual := string(captured)
	if strings.Contains(actual, "duplicate entry") {
		t.Errorf("No duplicates expected.")
	}
}

func TestProcessLineValue(t *testing.T) {
	sut.currentVariable = "TESTING"
	sut.addToCurrentVariable("")

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
