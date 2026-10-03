// Package table parses the repo's marker-delimited markdown tables and the
// numeric cells inside them.
package table

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Table markers wrap every table in the repo.
const (
	BeginMarker = "<!-- begin table -->"
	EndMarker   = "<!-- end table -->"
)

// Table is a list of rows; each row is a list of cells.
type Table [][]string

// Located is a table's position and shape inside a file.
type Located struct {
	Begin int    `json:"begin"`
	End   int    `json:"end"`
	Lines string `json:"lines"`
	Cols  int    `json:"cols"`
	Rows  int    `json:"rows"`
}

var (
	delimiterRe = regexp.MustCompile(`^:?-{3,}:?$`)
	numberRe    = regexp.MustCompile(`(\d[\d.,]*)`)
	percentRe   = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*%`)
	imageRe     = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	linkRe      = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	codeRe      = regexp.MustCompile("`([^`]*)`")
	boldRe      = regexp.MustCompile(`\*\*([^*]+)\*\*`)
)

// SplitRow splits a table row into trimmed cells (leading/trailing pipes dropped).
func SplitRow(line string) []string {
	s := strings.TrimSpace(line)
	s = strings.TrimPrefix(s, "|")
	s = strings.TrimSuffix(s, "|")
	return strings.Split(s, "|")
}

// IsDelimiter reports whether a row is the `| --- |` separator.
func IsDelimiter(line string) bool {
	cells := SplitRow(line)
	if len(cells) == 0 {
		return false
	}
	for _, c := range cells {
		if !delimiterRe.MatchString(strings.TrimSpace(c)) {
			return false
		}
	}
	return true
}

// Parse returns the tables inside marker blocks.
func Parse(md string) []Table {
	var tables []Table
	var rows [][]string
	in := false
	for _, line := range strings.Split(md, "\n") {
		t := strings.TrimSpace(line)
		switch t {
		case BeginMarker:
			in = true
			rows = nil
			continue
		case EndMarker:
			in = false
			if len(rows) > 0 {
				tables = append(tables, rows)
			}
			continue
		}
		if !in || t == "" || IsDelimiter(t) {
			continue
		}
		rows = append(rows, SplitRow(t))
	}
	return tables
}

// Scan reads a file and returns the located marker-delimited tables.
func Scan(path string) ([]Located, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []Located
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	line := 0
	in := false
	cur := Located{}
	rows := 0
	cols := 0
	for sc.Scan() {
		line++
		t := strings.TrimSpace(sc.Text())
		switch t {
		case BeginMarker:
			in = true
			cur = Located{Begin: line}
			rows = 0
			cols = 0
			continue
		case EndMarker:
			if in {
				cur.End = line
				cur.Lines = "L" + strconv.Itoa(cur.Begin) + "-L" + strconv.Itoa(cur.End)
				cur.Cols = cols
				cur.Rows = rows
				out = append(out, cur)
			}
			in = false
			continue
		}
		if !in || t == "" || IsDelimiter(t) {
			continue
		}
		if n := len(SplitRow(t)); n > cols {
			cols = n
		}
		rows++
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// Clean strips bold markers and trims a cell.
func Clean(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "**", ""))
}

// BRL parses a Brazilian-real cell into a float (handles a leading minus).
func BRL(s string) (float64, bool) {
	s = Clean(s)
	neg := strings.HasPrefix(s, "−") || strings.HasPrefix(s, "-")
	m := numberRe.FindString(s)
	if m == "" {
		return 0, false
	}
	// Brazilian format: `.` is the thousands separator, `,` the decimal point.
	m = strings.ReplaceAll(m, ".", "")
	m = strings.ReplaceAll(m, ",", ".")
	v, err := strconv.ParseFloat(m, 64)
	if err != nil {
		return 0, false
	}
	if neg {
		v = -v
	}
	return v, true
}

// Percent parses a percentage cell into a float.
func Percent(s string) (float64, bool) {
	m := percentRe.FindStringSubmatch(Clean(s))
	if m == nil {
		return 0, false
	}
	v, err := strconv.ParseFloat(m[1], 64)
	return v, err == nil
}

// Approx reports whether a and b are within tol.
func Approx(a, b, tol float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= tol
}

// RenderedLength counts a cell's visible characters (markdown syntax stripped).
func RenderedLength(cell string) int {
	s := imageRe.ReplaceAllString(cell, "")
	s = linkRe.ReplaceAllString(s, "$1")
	s = codeRe.ReplaceAllString(s, "$1")
	s = boldRe.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "*", "")
	return utf8.RuneCountInString(strings.TrimSpace(s))
}
