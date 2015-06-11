package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

var sectionRegex = regexp.MustCompile(`^\[(.*)\]$`)

type pokenv struct {
	environment      map[string][]string
	setContainsValue map[string]bool
	currentVariable  string
	registry         Registry
}

func (p *pokenv) setEnv(key int, fileName string) {
	p.processFile(fileName)
	p.setVars(key)
}

func (p *pokenv) setVars(key int) {
	for variable, values := range p.environment {
		if p.firstValueIsEmpty(values) {
			log.Println("Deleting", variable)
			p.registry.DeleteValue(key, variable)
		} else {
			joined := strings.Join(values, ";")
			log.Printf("Setting `%s` to `%s`\n", variable, joined)
			p.registry.SetString(key, variable, joined)
		}
	}
}

func (p *pokenv) processFile(fileName string) {

	p.environment = make(map[string][]string)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		p.processLine(strings.TrimSpace(line))
	}
}

func (p *pokenv) processLine(line string) {
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
		p.processSection(match[1])
	} else {
		p.processValue(trimComments(line))
	}
}

func (p *pokenv) processSection(section string) {
	p.currentVariable = strings.Replace(section, " ", "", -1)

	// mark section as empty (required for deletion)
	p.addCurrent("")
}

func (p *pokenv) processValue(value string) {
	if p.setContainsValue[value] {
		log.Println("Fatal error: duplicate entry:", value)
		log.Fatalln("Aborting, no value is set.")
	} else {
		// if this is first value, initialize and set value
		if p.firstValueIsEmpty(p.environment[p.currentVariable]) {
			p.setContainsValue = make(map[string]bool)
			p.setFirst(value)
		} else {
			p.addCurrent(value)
		}
		p.setContainsValue[value] = true
	}
}

func (p *pokenv) firstValueIsEmpty(values []string) bool {
	return values[0] == ""
}

func (p *pokenv) setFirst(value string) {
	p.environment[p.currentVariable][0] = value
}

func (p *pokenv) addCurrent(value string) {
	p.environment[p.currentVariable] = append(p.environment[p.currentVariable], value)
}

func trimComments(s string) string {
	if idx := strings.Index(s, "#"); idx != -1 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}
