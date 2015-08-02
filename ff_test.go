package main

import (
	"bytes"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
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

func makeExpected(input string) string {
	return strings.Replace(
		strings.Replace(input, "<", highlightStart, -1),
		">", highlightEnd, -1)
}

func TestHighlight(t *testing.T) {
	// Force color, even though we are not a tty.
	color.NoColor = false

	for i, test := range tests {
		r := regexp.MustCompile(test.regexp)
		actual := highlight(r, test.input)
		expected := makeExpected(test.output)
		if actual != expected {
			t.Errorf("%d: input: '%s', wanted '%s', got '%v'", i, test.input, expected, actual)
		}
	}
}

type fakeFileInfo struct {
	dir bool
}

func (f *fakeFileInfo) Name() string       { return "meh" }
func (f *fakeFileInfo) Sys() interface{}   { return nil }
func (f *fakeFileInfo) ModTime() time.Time { return time.Now() }
func (f *fakeFileInfo) IsDir() bool        { return f.dir }
func (f *fakeFileInfo) Size() int64        { return 0 }
func (f *fakeFileInfo) Mode() os.FileMode {
	if f.dir {
		return 0755 | os.ModeDir
	}
	return 0644
}

func TestWalkFunNoColor(t *testing.T) {
	buf := bytes.NewBufferString("")
	fun := buildWalkFun("a", buf, false)
	fun(".", nil, nil)
	fun("/meh/foo", &fakeFileInfo{false}, nil)
	fun("/meh/foodir", &fakeFileInfo{true}, nil)
	fun("/meh/yay", &fakeFileInfo{false}, nil)
	fun("/meh/yaydir", &fakeFileInfo{true}, nil)
	expected := "/meh/yay\n/meh/yaydir/\n"
	actual := buf.String()
	if expected != actual {
		t.Errorf("expected: '%s', got: %s", expected, actual)
	}
}

func TestWalkFunColor(t *testing.T) {
	buf := bytes.NewBufferString("")
	fun := buildWalkFun("a", buf, true)
	fun("/foo/hah", &fakeFileInfo{false}, nil)
	expected := makeExpected("/foo/h<a>h\n")
	actual := buf.String()
	if expected != actual {
		t.Errorf("expected: '%s', got: %s", expected, actual)
	}
}
