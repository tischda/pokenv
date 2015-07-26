package main

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
)

var sectionRegex = regexp.MustCompile(`^\[(.*)\]$`)

// TODO: reused in pokenv
type environment map[string][]string

type parser struct {
	env              environment
	setContainsValue map[string]bool
	currentVariable  string
}

func (p *parser) processAllLines(r io.Reader) environment {
	p.env = make(map[string][]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		p.processLine(strings.TrimSpace(line))
	}
	return p.env
}

func (p *parser) processLine(line string) {
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

func (p *parser) processSection(section string) {
	p.currentVariable = strings.Replace(section, " ", "", -1)

	// mark section as empty (required for deletion)
	if p.env[p.currentVariable] == nil {
		p.addToCurrentVariable("")
	}
}

func (p *parser) processValue(value string) {
	if p.setContainsValue[value] {
		log.Println("Warning: duplicate entry:", value)
	} else {

		// if this is first value, initialize and set value
		if firstValueIsEmpty(p.env[p.currentVariable]) {
			p.setContainsValue = make(map[string]bool)
			p.setFirst(value)
		} else {
			p.addToCurrentVariable(value)
		}
		p.setContainsValue[value] = true
	}
}

func (p *parser) setFirst(value string) {
	p.env[p.currentVariable][0] = value
}

// add value to current variable in environment
func (p *parser) addToCurrentVariable(value string) {
	p.env[p.currentVariable] = append(p.env[p.currentVariable], value)
}

// ------------------ utility

func trimComments(s string) string {
	if idx := strings.Index(s, "#"); idx != -1 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}

// TODO: this function is reused in pokenv
func firstValueIsEmpty(values []string) bool {
	return values[0] == ""
}
