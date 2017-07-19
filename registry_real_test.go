// +build windows

package main

import (
	"testing"
)

func TestSetDeleteValue(t *testing.T) {

	var registry = realRegistry{}

	expected := "titi"

	// set value
	err := registry.SetString(REG_KEY_USER, "toto", expected)
	if err != nil {
		t.Errorf("Error in SetString", err)
	}

	// read value
	actual, err := registry.GetString(REG_KEY_USER, "toto")
	if err != nil {
		t.Errorf("Error in GetString", err)
	}
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}

	// delete value
	err = registry.DeleteValue(REG_KEY_USER, "toto")
	if err != nil {
		t.Errorf("Error in DeleteValue", err)
	}
}
