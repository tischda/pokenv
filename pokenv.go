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

var environment map[string][]string
var setContainsValue map[string]bool

var currentVariable string

func setEnv(key int, fileName string) {
	processFile(fileName)
	setVars(key)
}

func setVars(key int) {
	for variable, values := range environment {
		joined := strings.Join(values, ";")
		fmt.Printf("setting `%s` to `%s`\n", variable, joined)
		registry.SetString(key, variable, joined)
	}
}

func processFile(fileName string) {

	environment = make(map[string][]string)

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
		processValue(trimComments(line))
	}
}

func processSection(section string) {
	currentVariable = strings.Replace(section, " ", "", -1)
	setContainsValue = make(map[string]bool)
}

func processValue(value string) {
	if setContainsValue[value] {
		// log.Fatalln("Error, duplicate entry:", value)
		log.Println("Error, duplicate entry:", value)
	} else {
		setContainsValue[value] = true
		environment[currentVariable] = append(environment[currentVariable], value)
	}
}

func trimComments(s string) string {
	if idx := strings.Index(s, "#"); idx != -1 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}
