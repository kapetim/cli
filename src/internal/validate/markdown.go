// Package validate holds repository validation commands.
//
// The markdown checks here are deliberately dependency-free: table structure
// is tracked with the begin/end markers and table cell width is measured on
// the rendered text (markdown syntax stripped). A goldmark-based engine with
// the full rule set is planned.
package validate

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	beginMarker = "<!-- begin table -->"
	endMarker   = "<!-- end table -->"
)

var (
	imageRe = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	linkRe  = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	codeRe  = regexp.MustCompile("`([^`]*)`")
	boldRe  = regexp.MustCompile(`\*\*([^*]+)\*\*`)
)

// renderedLength returns the length of a table cell once markdown syntax is
// stripped, counted in Unicode code points (matching a rendered cell width).
func renderedLength(cell string) int {
	s := cell
	s = imageRe.ReplaceAllString(s, "")
	s = linkRe.ReplaceAllString(s, "$1")
	s = codeRe.ReplaceAllString(s, "$1")
	s = boldRe.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "*", "")
	return utf8.RuneCountInString(strings.TrimSpace(s))
}

type issue struct {
	file string
	line int
	msg  string
}

func (i issue) Error() string {
	return fmt.Sprintf("%s:%d: %s", i.file, i.line, i.msg)
}

func splitRow(line string) []string {
	s := strings.TrimSpace(line)
	s = strings.TrimPrefix(s, "|")
	s = strings.TrimSuffix(s, "|")
	return strings.Split(s, "|")
}

func isDelimiterRow(line string) bool {
	s := strings.TrimSpace(line)
	if !strings.HasPrefix(s, "|") {
		return false
	}
	for _, c := range splitRow(s) {
		t := strings.Trim(strings.TrimSpace(c), ":")
		if t == "" || strings.Trim(t, "-") != "" {
			return false
		}
	}
	return true
}

func checkFile(path string, maxCell int) ([]issue, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var issues []issue
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	inTable := false
	line := 0
	for scanner.Scan() {
		line++
		s := strings.TrimSpace(scanner.Text())
		switch s {
		case beginMarker:
			if inTable {
				issues = append(issues, issue{path, line, "nested table begin marker"})
			}
			inTable = true
			continue
		case endMarker:
			if !inTable {
				issues = append(issues, issue{path, line, "table end without a begin marker"})
			}
			inTable = false
			continue
		}
		if !inTable || isDelimiterRow(s) {
			continue
		}
		for _, cell := range splitRow(s) {
			if n := renderedLength(cell); n > maxCell {
				issues = append(issues, issue{path, line, fmt.Sprintf("table cell renders to %d chars (max %d)", n, maxCell)})
			}
		}
	}
	if inTable {
		issues = append(issues, issue{path, line, "unclosed table begin marker"})
	}
	return issues, scanner.Err()
}

func collect(targets []string) ([]string, error) {
	var files []string
	for _, t := range targets {
		info, err := os.Stat(t)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			files = append(files, t)
			continue
		}
		err = filepath.WalkDir(t, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == ".git" || d.Name() == "node_modules" {
					return fs.SkipDir
				}
				return nil
			}
			if strings.EqualFold(filepath.Ext(d.Name()), ".md") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

// Run validates markdown table structure and rendered cell widths.
func Run(args []string) error {
	flags := flag.NewFlagSet("validate", flag.ContinueOnError)
	maxCell := flags.Int("max-cell", 99, "maximum rendered table cell character length")
	if err := flags.Parse(args); err != nil {
		return err
	}
	targets := flags.Args()
	if len(targets) == 0 {
		targets = []string{"."}
	}
	files, err := collect(targets)
	if err != nil {
		return err
	}
	var issues []issue
	for _, f := range files {
		found, err := checkFile(f, *maxCell)
		if err != nil {
			return err
		}
		issues = append(issues, found...)
	}
	for _, is := range issues {
		fmt.Fprintln(os.Stderr, is.Error())
	}
	if len(issues) > 0 {
		return fmt.Errorf("%d validation error(s)", len(issues))
	}
	fmt.Printf("cli validate: %d file(s) ok\n", len(files))
	return nil
}
