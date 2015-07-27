package main

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var sutParser parser

func init() {
	sutParser = parser{}
	log.SetOutput(ioutil.Discard)
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

func TestEmpty(t *testing.T) {
	contents := `[SECTION-E]`
	expected := varMap{"SECTION-E": {}}
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

// reflect.DeepEqual(m1, m2) breaks because
// order of values in maps is random.
func deepEqual(v1 varMap, v2 varMap) bool {
	if (v1 == nil) != (v2 == nil) {
		return false
	}
	if len(v1) != len(v2) {
		return false
	}
	for k, s1 := range v1 {
		s2 := v2[k]
		if len(s1) != len(s2) {
			return false
		}
		sort.Strings(s1)
		sort.Strings(s2)
		for i, v := range s1 {
			if s2[i] != v {
				return false
			}
		}
	}
	return true
}
