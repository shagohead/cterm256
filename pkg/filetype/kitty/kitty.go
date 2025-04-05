package kitty

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/shagohead/cterm256/pkg/filetype"
	"github.com/shagohead/cterm256/pkg/termcolor"
)

func init() {
	filetype.Register("kitty", &fileType{})
}

type fileType struct{}

// Parse implements ftypes.FileType.
func (f *fileType) Parse(input io.Reader) (termcolor.Table, error) {
	var err error
	cs := &colorScheme{named: make(map[string]termcolor.Color)}
	scan := bufio.NewScanner(input)
	var ln int
	for scan.Scan() {
		line := ltrim(scan.Bytes())
		if l := len(line); l == 0 || line[0] == '#' {
			if l > 0 {
				cs.raw = append(cs.raw, line)
			}
			continue
		}
		err = scanLine(bytes.NewReader(line), cs)
		if errors.Is(err, errCannotParseLine) {
			cs.raw = append(cs.raw, line)
		} else if err != nil {
			return nil, fmt.Errorf("%d line: %v", ln, err)
		}
		ln++
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	if c, f := cs.named["cursor"], cs.named["foreground"]; c == nil && f != nil {
		cs.named["cursor"] = f
	}
	if c, b := cs.named["cursor_text_color"], cs.named["background"]; c == nil && b != nil {
		cs.named["cursor_text_color"] = b
	}
	return cs, nil
}

func ltrim(s []byte) []byte {
	for i, c := range s {
		switch c {
		case ' ', '\t':
			continue
		default:
			return s[i:]
		}
	}
	return nil
}

type keyword int

const (
	unknownKeyword keyword = iota
	namedColor
	indexedColor
)

var errCannotParseLine = errors.New("cannot parse line")

func scanLine(line io.Reader, cs *colorScheme) (err error) {
	scan := bufio.NewScanner(line)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return scan.Err()
	}
	switch word := scan.Text(); word {
	case "foreground",
		"background",
		"cursor",
		"cursor_text_color":
		c, err := scanColorValue(scan)
		if err != nil {
			return fmt.Errorf("%q: %w", word, err)
		}
		cs.named[word] = c
	default:
		if len(word) < 6 || !strings.HasPrefix(word, "color") {
			return errCannotParseLine
		}
		n, err := strconv.Atoi(word[5:])
		if err != nil {
			return fmt.Errorf("%q: parse number: %v", word, err)
		}
		c, err := scanColorValue(scan)
		if err != nil {
			return fmt.Errorf("%q: %w", word, err)
		}
		cs.indexed[n] = c
	}
	return scan.Err()
}

var (
	errMissingColorValue = errors.New("missing color value")
)

func scanColorValue(scan *bufio.Scanner) (termcolor.Color, error) {
	if !scan.Scan() {
		return nil, errMissingColorValue
	}
	v := scan.Text()
	if !termcolor.HEX.MatchString(v) {
		return nil, errCannotParseLine
	}
	return termcolor.FromHEX(scan.Text()), nil
}

// Support implements ftypes.FileType.
func (f *fileType) Support(name string, ext string) bool {
	if ext == ".conf" {
		return true
	}
	return false
}

var _ filetype.FileType = (*fileType)(nil)

type colorScheme struct {
	raw     [][]byte
	named   map[string]termcolor.Color
	indexed [256]termcolor.Color
}

// Write implements termcolor.Table.
func (cs *colorScheme) Write(w termcolor.Writer) error {
	s := &strings.Builder{}
	for _, r := range cs.raw {
		s.Write(r)
		s.WriteByte('\n')
	}
	for n, c := range cs.named {
		s.WriteString(n)
		writeColor(s, c)
	}
	for i, c := range cs.indexed {
		s.WriteString("color")
		s.WriteString(strconv.Itoa(i))
		writeColor(s, c)
	}
	_, err := w.WriteString(s.String())
	return err
}

func writeColor(s *strings.Builder, c termcolor.Color) {
	s.WriteRune(' ')
	s.WriteString(c.HEX())
	s.WriteRune('\n')
}

// SetColor implements termcolor.Table.
func (cs *colorScheme) SetColor(number int, color termcolor.Color) {
	cs.indexed[number] = color
}

// Background implements termcolor.Table.
func (cs *colorScheme) Background() termcolor.Color {
	return cs.named["background"]
}

// Color implements termcolor.Table.
func (cs *colorScheme) Color(number int) termcolor.Color {
	if number > 255 || number < 0 {
		panic("color number out of bounds")
	}
	return cs.indexed[number]
}

// Foreground implements termcolor.Table.
func (cs *colorScheme) Foreground() termcolor.Color {
	return cs.named["foreground"]
}

var _ termcolor.Table = (*colorScheme)(nil)
