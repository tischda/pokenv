package main

import (
	"log"
	"strings"
)

// If there is a comment in the line, return only
// content to the left of the comment marker '#'.
func trimComments(s string) string {
	if idx := strings.Index(s, "#"); idx != -1 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}

// Remove all spaces in line.
func trimSpaces(s string) string {
	return strings.Replace(s, " ", "", -1)
}

func checkFatal(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
