package main

import (
	"github.com/fatih/color"
	"regexp"
	"strings"
	"testing"
)

const (
	highlightStart = "\x1b[31m"
	highlightEnd   = "\x1b[0m"
)

var tests = []struct {
	regexp string
	input  string
	output string
}{
	{"a", "nope", "nope"},
	{"a", "a", "<a>"},
	{"a", "ba", "b<a>"},
	{"a", "ab", "<a>b"},
	{"a", "aba", "<a>b<a>"},
	{"a", "aa", "<a><a>"},
	{"a+", "aa", "<aa>"},
}

func TestHighlight(t *testing.T) {
	// Force color, even though we are not a tty.
	color.NoColor = false

	for i, test := range tests {
		r := regexp.MustCompile(test.regexp)
		actual := highlight(r, test.input)
		expected := strings.Replace(
			strings.Replace(test.output, "<", highlightStart, -1),
			">", highlightEnd, -1)
		if actual != expected {
			t.Errorf("%d: input: '%s', wanted '%s', got '%v'", i, test.input, expected, actual)
		}
	}
}
