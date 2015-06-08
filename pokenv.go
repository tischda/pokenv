package main

import (
	"bufio"
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
		if firstValueIsEmpty(values) {
			log.Println("Deleting", variable)
			registry.DeleteValue(key, variable)
		} else {
			joined := strings.Join(values, ";")
			log.Printf("Setting `%s` to `%s`\n", variable, joined)
			registry.SetString(key, variable, joined)
		}
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

	// mark section as empty (required for deletion)
	addCurrent("")
}

func processValue(value string) {
	if setContainsValue[value] {
		log.Println("Fatal error: duplicate entry:", value)
		log.Fatalln("Aborting, no value is set.")
	} else {
		// if this is first value, initialize and set value
		if firstValueIsEmpty(environment[currentVariable]) {
			setContainsValue = make(map[string]bool)
			setFirst(value)
		} else {
			addCurrent(value)
		}
		setContainsValue[value] = true
	}
}

func firstValueIsEmpty(values []string) bool {
	return values[0] == ""
}

func setFirst(value string) {
	environment[currentVariable][0] = value
}

func addCurrent(value string) {
	environment[currentVariable] = append(environment[currentVariable], value)
}

func trimComments(s string) string {
	if idx := strings.Index(s, "#"); idx != -1 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}
