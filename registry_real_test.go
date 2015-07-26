// +build windows

package main

import (
	"testing"
)

func TestSetDeleteValue(t *testing.T) {

	var registry = realRegistry{}

	expected := "titi"

	// set value
	err := registry.SetString(PATH_USER, "toto", expected)
	if err != nil {
		t.Errorf("Error in SetString", err)
	}

	// read value
	actual, err := registry.GetString(PATH_USER, "toto")
	if err != nil {
		t.Errorf("Error in GetString", err)
	}
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}

	// delete value
	err = registry.DeleteValue(PATH_USER, "toto")
	if err != nil {
		t.Errorf("Error in DeleteValue", err)
	}
}
