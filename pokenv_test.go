package main

import (
	"io"
	"log"
	"os"
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
	paths := []string{
		`c:\Windows`,
		`c:\Windows\system32`,
		`%windir%`,
		`%windir%\system32`,
		`.`,
	}
	for _, path := range paths {
		if isPathInvalid(path) {
			t.Errorf("Invalid path: %q", path)
		}
	}
}

// TODO: re-enable when path checking is implemented
// Inspired by https://talks.golang.org/2014/testing.slide#23
/*
func TestParseAndCheckPaths(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		log.SetFlags(0)
		cfg := &Config{}
		sut_pokenv = pokenv{mock, cfg}
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
*/

func TestDeleteSectionRemovesVariable(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stdout)

	// Set up the mock environment with the variable to be deleted
	mock.env = map[string]string{
		"deleteme": "should_be_deleted",
	}

	sut_pokenv = pokenv{registry: mock}
	ret := sut_pokenv.processFile(REG_KEY_MACHINE, `data/deletesection.txt`)

	if ret != nil {
		t.Errorf("Expected no error, but got: %v", ret)
	}

	if _, exists := mock.env["deleteme"]; exists {
		t.Errorf("Expected 'deleteme' variable to be deleted, but it still exists")
	}
}
