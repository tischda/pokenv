// +build !windows

package main

var p pokenv = pokenv{
	registry:  mockRegistry{},
	pathcheck: check,
}
