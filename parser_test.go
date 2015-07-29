package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"reflect"
)

var sutParser parser

func init() {
	sutParser = parser{}
}

func parseContents(contents string) varMap {
	return sutParser.processAllLines(strings.NewReader(contents))
}

func TestSimple(t *testing.T) {
	contents := `[SECTION-S]
	value1
	value2
	`
	expected := varMap{"SECTION-S": {"value1", "value2"}}
	assertDeepEquals(t, expected, parseContents(contents))
}

func TestOrdered(t *testing.T) {
	contents := `[SECTION-R]
	1
	2
	3
	4
	5
	6
	7
	8
	9
	10
	`
	expected := varMap{"SECTION-R": {"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}}
	assertDeepEquals(t, expected, parseContents(contents))
}


func TestOrphan(t *testing.T) {
	contents := `value1
	[SECTION-O]
	value2
	`
	expected := varMap{"SECTION-O": {"value2"}}
	assertDeepEquals(t, expected, parseContents(contents))
}

func TestDouble(t *testing.T) {
	contents := `[SECTION-D]
	value1
	[SECTION-D]
	value2
	`
	expected := varMap{"SECTION-D": {"value1", "value2"}}
	assertDeepEquals(t, expected, parseContents(contents))
}

func TestDoubleEmpty(t *testing.T) {
	contents := `[SECTION-DE]
	[SECTION-DE]
	value2
	`
	expected := varMap{"SECTION-DE": {"value2"}}
	assertDeepEquals(t, expected, parseContents(contents))

	contents = `[SECTION-DE]
	value1
	[SECTION-DE]`
	expected = varMap{"SECTION-DE": {"value1"}}
	assertDeepEquals(t, expected, parseContents(contents))

	contents = `[SECTION-DE]
	[SECTION-DE]`
	expected = varMap{"SECTION-DE": {}}
	assertDeepEquals(t, expected, parseContents(contents))
}

func TestEmptySingle(t *testing.T) {
	contents := `[SECTION-E]`
	expected := varMap{"SECTION-E": {}}
	assertDeepEquals(t, expected, parseContents(contents))
}

func TestEmptyNotLast(t *testing.T) {
	contents := `[SECTION-ENL]
	[SECTION-OTHER]
	value
	`
	expected := varMap{"SECTION-ENL": {}, "SECTION-OTHER": {"value"}}
	assertDeepEquals(t, expected, parseContents(contents))
}


func TestDuplicates(t *testing.T) {
	contents := `[SECTION-D]
	dupvalue
	dupvalue
	`
	// capture output
	r, w, _ := os.Pipe()
	log.SetOutput(w)
	defer func() { log.SetOutput(os.Stdout) }()

	// do work
	parseContents(contents)
	w.Close()

	// now check that message is displayed
	captured, _ := ioutil.ReadAll(r)

	actual := string(captured)
	expected := "duplicate entry"

	if !strings.Contains(actual, expected) {
		t.Errorf("Expected: %s, but was: %s", expected, actual)
	}
}

func TestNoDuplicates(t *testing.T) {
	contents := `[SECTION-A]
	nodupvalue
	[SECTION-B]
	nodupvalue`

	// capture output
	r, w, _ := os.Pipe()
	log.SetOutput(w)
	defer func() { log.SetOutput(os.Stdout) }()

	// do work
	parseContents(contents)
	w.Close()

	captured, _ := ioutil.ReadAll(r)
	actual := string(captured)
	if strings.Contains(actual, "duplicate entry") {
		t.Errorf("No duplicates expected.")
	}
}

func assertDeepEquals(t *testing.T, expected varMap, actual varMap) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}
