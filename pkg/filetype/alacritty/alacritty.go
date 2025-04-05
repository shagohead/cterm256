package alacritty

import (
	"errors"
	"fmt"
	"io"

	"github.com/pelletier/go-toml/v2"

	"github.com/shagohead/cterm256/pkg/filetype"
	"github.com/shagohead/cterm256/pkg/termcolor"
)

func init() {
	filetype.Register("alacritty", &fileType{})
}

type fileType struct{}

// Parse implements ftypes.FileType.
func (f *fileType) Parse(input io.Reader) (termcolor.Table, error) {
	cs := new(colorScheme)
	if err := toml.NewDecoder(input).Decode(&cs.config); err != nil {
		return nil, err
	}
	colorsv, ok := cs.config["colors"]
	if !ok {
		return nil, errors.New(`missing "colors" key`)
	}
	colors, ok := colorsv.(map[string]any)
	if !ok {
		return nil, fmt.Errorf(`colors: unexpected type %T`, colorsv)
	}
	if err := cs.parseIndexedSection(colors); err != nil {
		return nil, err
	}
	if err := cs.parseBaseColors(colors, "normal", 0); err != nil {
		return nil, err
	}
	if err := cs.parseBaseColors(colors, "bright", 8); err != nil {
		return nil, err
	}
	primaryv, ok := colors["primary"]
	if !ok {
		return cs, nil
	}
	primary, ok := primaryv.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("primary: unexpected type %T", primaryv)
	}
	if err := cs.parsePrimaryColor(primary, "background", &cs.background); err != nil {
		return nil, err
	}
	if err := cs.parsePrimaryColor(primary, "foreground", &cs.foreground); err != nil {
		return nil, err
	}
	return cs, nil
}

// Support implements ftypes.FileType.
func (f *fileType) Support(name string, ext string) bool {
	if ext == ".toml" {
		return true
	}
	return false
}

var _ filetype.FileType = (*fileType)(nil)

type colorScheme struct {
	indexed    [256]termcolor.Color
	background termcolor.Color
	foreground termcolor.Color
	config     map[string]any
}

// Write implements termcolor.Table.
func (cs *colorScheme) Write(w termcolor.Writer) error {
	var colors map[string]any
	if v, ok := cs.config["colors"]; ok {
		colors = v.(map[string]any)
	} else {
		colors = make(map[string]any)
		cs.config["colors"] = colors
	}

	for _, sec := range []string{"normal", "bright", "primary"} {
		if _, ok := colors[sec]; !ok {
			colors[sec] = make(map[string]any)
		}
	}

	if err := cs.setIndexedColors(colors); err != nil {
		return err
	}
	cs.setPrimaryColor(colors, "background", cs.background)
	cs.setPrimaryColor(colors, "foreground", cs.foreground)
	return toml.NewEncoder(w).Encode(cs.config)
}

func (cs *colorScheme) setPrimaryColor(colors map[string]any, key string, c termcolor.Color) {
	if c == nil {
		return
	}
	colors["primary"].(map[string]any)[key] = c.HEX()
}

func (cs *colorScheme) setIndexedColors(colors map[string]any) error {
	initial := make(map[int]map[string]any, 256)
	parseIndexedSection(colors, func(m map[string]any) error {
		n, err := parseIndexValue(m)
		if err != nil {
			return err
		}
		initial[n] = m
		return nil
	})

	indexed := make([]any, 0, 256)
	for n, c := range cs.indexed {
		if c == nil {
			if orig, ok := initial[n]; ok {
				indexed = append(indexed, orig)
			}
			continue
		}
		if n < 16 {
			bnum := n
			bname := "normal"
			if n > 7 {
				bnum -= 8
				bname = "bright"
			}
			colors[bname].(map[string]any)[baseColorName(bnum)] = c.HEX()
		}
		indexed = append(indexed, map[string]any{"index": n, "color": c.HEX()})
	}
	colors["indexed_colors"] = indexed
	return nil
}

func (cs *colorScheme) parsePrimaryColor(src map[string]any, key string, dst *termcolor.Color) error {
	val, ok := src[key]
	if !ok {
		return nil
	}
	col, ok := val.(string)
	if !ok {
		return fmt.Errorf("primary.%s: unexpected type %T", key, val)
	}
	*dst = termcolor.FromHEX(col)
	return nil
}

func parseIndexedSection(src map[string]any, cb func(map[string]any) error) error {
	indexed, ok := src["indexed_colors"]
	if !ok {
		return nil
	}
	switch indexed := indexed.(type) {
	case []any:
		for i, s := range indexed {
			m, ok := s.(map[string]any)
			if !ok {
				return fmt.Errorf("indexed_colors[%d]: unexpected type %T", i, s)
			}
			if err := cb(m); err != nil {
				return fmt.Errorf("indexed_colors[%d]: %v", i, err)
			}
		}
	case []map[string]any:
		for i, s := range indexed {
			if err := cb(s); err != nil {
				return fmt.Errorf("indexed_colors[%d]: %v", i, err)
			}
		}
	default:
		return fmt.Errorf("indexed_colors: unexpected type %T", indexed)
	}
	return nil
}

func (cs *colorScheme) parseIndexedSection(src map[string]any) error {
	return parseIndexedSection(src, func(m map[string]any) error {
		return cs.setIndexedColor(m)
	})
}

func parseIndexValue(src map[string]any) (int, error) {
	idxv, ok := src["index"]
	if !ok {
		return 0, errors.New(`missing "index" key`)
	}
	var idx int
	switch v := idxv.(type) {
	case int:
		idx = v
	case int32:
		idx = int(v)
	case int64:
		idx = int(v)
	default:
		return 0, fmt.Errorf("index: unexpected type %T", v)
	}
	return idx, nil
}

func (cs *colorScheme) setIndexedColor(src map[string]any) error {
	idx, err := parseIndexValue(src)
	if err != nil {
		return err
	}
	colv, ok := src["color"]
	if !ok {
		return errors.New(`missing "color" key`)
	}
	col, ok := colv.(string)
	if !ok {
		return fmt.Errorf("color: unexpected type %T", colv)
	}
	cs.indexed[idx] = termcolor.FromHEX(col)
	return nil
}

func (cs *colorScheme) parseBaseColors(src map[string]any, section string, offset int) error {
	colorsv, ok := src[section]
	if !ok {
		return nil
	}
	colors, ok := colorsv.(map[string]any)
	if !ok {
		return fmt.Errorf("%s: unexpected type %T", section, colorsv)
	}
	for name, val := range colors {
		idx := baseColorIndex(name)
		if idx < 0 {
			fmt.Printf("%s.%s: unknown color index\n", section, name)
			continue
		}
		col, ok := val.(string)
		if !ok {
			fmt.Printf("%s.%s: unexpected type %T\n", section, name, val)
			continue
		}
		cs.indexed[idx+offset] = termcolor.FromHEX(col)
	}
	return nil
}

func baseColorIndex(name string) int {
	for i, n := range baseColors {
		if n == name {
			return i
		}
	}
	return -1
}

func baseColorName(index int) string {
	if index < 8 {
		return baseColors[index]
	}
	return ""
}

var baseColors = []string{
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
}

// SetColor implements termcolor.Table.
func (cs *colorScheme) SetColor(number int, color termcolor.Color) {
	cs.indexed[number] = color
}

// Color implements termcolor.Table.
func (cs *colorScheme) Color(number int) termcolor.Color {
	if number > 255 || number < 0 {
		panic("color number out of bounds")
	}
	return cs.indexed[number]
}

// Background implements termcolor.Table.
func (cs *colorScheme) Background() termcolor.Color {
	return cs.background
}

// Foreground implements termcolor.Table.
func (cs *colorScheme) Foreground() termcolor.Color {
	return cs.foreground
}

var _ termcolor.Table = (*colorScheme)(nil)
