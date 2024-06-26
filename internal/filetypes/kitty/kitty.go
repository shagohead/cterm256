package kitty

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/shagohead/cterm256/internal/cterm"
)

const (
	foregroundKeyword = `foreground`
	backgroundKeyword = `background`
)

type FileType struct{}

// Parse implements parser.FileType.
func (*FileType) Parse(input io.Reader) (*cterm.ColorScheme, error) {
	themeScanner := bufio.NewScanner(input)
	kvPattern, err := regexp.Compile(`(?P<key>[^\s]+)\s+#(?P<val>[^\s]+)`)
	if err != nil {
		return nil, err
	}
	idxPattern, err := regexp.Compile(`color(?P<idx>\d+)`)
	if err != nil {
		return nil, err
	}
	tm := &cterm.ColorScheme{
		Custom: make(map[string]cterm.Color),
	}
	for themeScanner.Scan() {
		if len(themeScanner.Bytes()) == 0 {
			continue
		}
		// TODO: Preserve commetaries
		matches := kvPattern.FindSubmatch(themeScanner.Bytes())
		if len(matches) < 1 {
			continue
		}
		var key, val []byte
		for i, m := range matches {
			switch kvPattern.SubexpNames()[i] {
			case "key":
				key = m
			case "val":
				val = m
			}
		}

		color := cterm.ColorFromHEX(string(val))

		idxMatch := idxPattern.FindSubmatch(key)
		if len(idxMatch) < 1 {
			custom := string(key)
			tm.Custom[custom] = color
			switch custom {
			case foregroundKeyword:
				tm.Foreground = &color
			case backgroundKeyword:
				tm.Background = &color
			}
			continue
		}

		var idx int
		for i, m := range idxMatch {
			switch idxPattern.SubexpNames()[i] {
			case "idx":
				idx, err = strconv.Atoi(string(m))
				if err != nil {
					return tm, fmt.Errorf("parsing index of color name `%s`: %w", m, err)
				}
			}
		}
		tm.Indexed[idx] = color
	}
	if err := themeScanner.Err(); err != nil {
		return tm, err
	}
	return tm, nil
}

// Write implements cterm.FileType.
func (f *FileType) Write(output io.Writer, theme *cterm.ColorScheme) error {
	for n, c := range theme.Custom {
		if err := writeColor(output, c, n); err != nil {
			return err
		}
	}
	for i, c := range theme.Indexed {
		if err := writeColor(output, c, "color", strconv.Itoa(i)); err != nil {
			return err
		}
	}
	return nil
}

func writeColor(w io.Writer, color cterm.Color, keyParts ...string) error {
	if ws, ok := w.(io.StringWriter); ok {
		return writeColorString(ws, color, keyParts...)
	} else {
		var ws bytes.Buffer
		if err := writeColorString(&ws, color, keyParts...); err != nil {
			return err
		}
		if _, err := w.Write(ws.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func writeColorString(ws io.StringWriter, color cterm.Color, keyParts ...string) error {
	keyParts = append(keyParts, " ", color.HEX(), "\n")
	for _, s := range keyParts {
		if _, err := ws.WriteString(s); err != nil {
			return err
		}
	}
	return nil
}

var _ cterm.FileType = (*FileType)(nil)
