// +build !windows

package main

var p pokenv = pokenv{
	environment: make(map[string][]string),
	registry:    mockRegistry{},
	pathcheck:   check,
}
