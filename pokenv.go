package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var sectionRegex = regexp.MustCompile(`^\[(.*)\]$`)

var registry mockRegistry
var env map[string][]string
var set map[string]bool

var currentVariable string

func setEnv(path string, fileName string) {
	processFile(fileName)
	setVars(path)
}

func setVars(path string) {
	for variable, values := range env {
		joined := strings.Join(values, ";")
		fmt.Printf("setting `%s` to `%s`\n", variable, joined)
		registry.SetString(path, variable, joined)
	}
}

func processFile(fileName string) {

	env = make(map[string][]string)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		processLine(strings.TrimSpace(line))
	}
}

func processLine(line string) {
	// ignore blank line
	if line == "" {
		return
	}
	// ignore comment
	if strings.HasPrefix(line, "#") {
		return
	}
	// sections and values
	match := sectionRegex.FindStringSubmatch(line)
	if match != nil {
		processSection(match[1])
	} else {
		processValue(line)
	}
}

func processSection(section string) {
	currentVariable = strings.TrimSpace(section)
	set = make(map[string]bool)
}

func processValue(line string) {
	if set[line] {
		fmt.Println("duplicate:", line)
	} else {
		set[line] = true
		env[currentVariable] = append(env[currentVariable], line)
	}
}
