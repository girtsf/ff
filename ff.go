// Overcomplicated replacement for <find . -name "$1">.
//
// This was written mostly as an excuse to learn Go. So it's probably ugly and
// non-idiomatic.

package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

var caseInsensitive bool

// Highlight colorifies all occurences of pattern inside of needle
// and returns a string.
func highlight(pattern *regexp.Regexp, needle string) string {
	matches := pattern.FindAllStringIndex(needle, -1)
	out := ""
	prev := 0
	printer := color.New(color.FgRed).SprintFunc()
	for _, locs := range matches {
		txt := needle[locs[0]:locs[1]]
		out += needle[prev:locs[0]]
		out += printer(txt)
		prev = locs[1]
	}
	out += needle[prev:]
	return out
}

// buildWalkFun builds a filepath.Walkfunc which outputs all
// files and directories that match given pattern.
//
// If 'color' is true, the matching substrings are highlighted using
// ANSI colors.
func buildWalkFun(pattern string, writer io.Writer, color bool) filepath.WalkFunc {
	if (caseInsensitive) {
		pattern = "(?i)" + pattern
	}

	r := regexp.MustCompile(pattern)
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == "." || path == ".." {
			return nil
		}
		if r.MatchString(path) {
			out := path
			if color {
				out = highlight(r, path)
			}
			if info.IsDir() {
				out += "/"
			}
			fmt.Fprintf(writer, "%s\n", out)
		}
		return nil
	}
}

func init() {
	flag.BoolVar(&caseInsensitive, "i", false, "case insensitive")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <regexp>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	var needle string
	if flag.NArg() == 0 {
		needle = ""
	} else if flag.NArg() == 1 {
		needle = flag.Arg(0)
	} else {
		flag.Usage()
		return
	}

	walkFun := buildWalkFun(needle, os.Stdout, true)
	err := filepath.Walk(".", walkFun)
	if err != nil {
		panic(err)
	}
}
