package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestMainPokenv(t *testing.T) {
	args := []string{"-version"}
	os.Args = append(os.Args, args...)

	expected := fmt.Sprintf("pokenv version %s\n", version)
	actual := captureOutput(main)

	if expected != actual {
		t.Errorf("Expected: %s, but was: %s", expected, actual)
	}
}

// captures Stdout and returns output of function f()
func captureOutput(f func()) string {
	// redirect output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	// reset output again
	w.Close()
	os.Stdout = old

	captured, _ := ioutil.ReadAll(r)
	return string(captured)
}
