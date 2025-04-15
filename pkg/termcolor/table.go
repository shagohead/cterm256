package termcolor

import (
	"fmt"
	"io"
)

// Color scheme API.
// Implemented in ftype packages.
type Table interface {
	// Color returns one of 256 colors from table.
	Color(number int) Color

	// SetColor [color] at index [number].
	SetColor(number int, color Color)

	// Return background color if exists.
	Background() Color

	// Return foreground color if exists.
	Foreground() Color

	Write(w Writer) error
}

type Writer interface {
	io.Writer
	io.StringWriter
}

// Generate 256 color palette based on first 8/16 + background & foreground.
func Generate(cs Table) error {
	for i := range 8 {
		if cs.Color(i) == nil {
			return fmt.Errorf("provided scheme missing color %d", i)
		}
	}
	background := cs.Background()
	if background == nil {
		fmt.Println("scheme missing background color; will use black/0 instead")
		background = cs.Color(0)
	}
	cs.SetColor(16, cs.Color(0).With(-1, -1, background.Lightness()))

	// Is it dark or light theme?
	isDark := cs.Color(1).Lightness() > cs.Color(16).Lightness()

	// R/g/b based gradients.
	for _, c := range []struct {
		brindex int    // Bright r/g/b
		targets [5]int // Indexes of tinted versions
	}{
		{9, [5]int{52, 88, 124, 160, 196}}, // Red
		{12, [5]int{17, 18, 19, 20, 21}},   // Green
		{10, [5]int{22, 28, 34, 40, 46}},   // Blue
	} {
		gradient(cs, isDark, c.brindex, c.targets)
	}

	contrast := maxContrast(isDark)
	cube(cs, contrast)

	foreground := cs.Foreground()
	if foreground == nil {
		white := cs.Color(7)
		foreground = white
		repr := "white/7"
		if brwhite := cs.Color(15); brwhite != nil {
			if contrast(brwhite.Lightness(), white.Lightness()) == brwhite.Lightness() {
				foreground = brwhite
				repr = "brwhite/15"
			}
		}
		fmt.Println("scheme missing foreground color; will use ", repr, " instead")
	}
	delta := 1.0 / 25
	for i := range 24 {
		s := delta * float64(i+1)
		cs.SetColor(232+i, background.Blend(foreground, s, s, s))
	}
	return nil
}

func maxContrast(isDark bool) func(a, b float64) float64 {
	if isDark {
		return func(a, b float64) float64 {
			if a > b {
				return a
			}
			return b
		}
	}
	return func(a, b float64) float64 {
		if a < b {
			return a
		}
		return b
	}
}

func gradient(cs Table, isDark bool, brindex int, targets [5]int) {
	bg := cs.Color(16)
	src := cs.Color(brindex)
	norm := cs.Color(brindex - 8)
	if isDark {
		if norm.Lightness() > src.Lightness() {
			src = norm
		}
	} else {
		if norm.Lightness() < src.Lightness() {
			src = norm
		}
	}
	for i, target := range targets {
		i := float64(i + 1)
		cs.SetColor(target, bg.Blend(src, 1, 0.1*i, 0.1*i))
	}
}

func cube(cs Table, maxContrast func(a, b float64) float64) {
	// Generate 6x6x6 colors cube
	for side := range 6 {
		perside := 5
		// RGB colors for every side.
		// Green is left/top on every side.
		// Red is left/bottom at first side + green.
		// Blue is right/top at first side + green.
		red := cs.Color(196)
		green := cs.Color(side*6 + 16)
		blue := cs.Color(21)
		if side > 0 {
			perside = 6
			s := float64(side)
			red = red.Blend(green, 0.1*s, 0.2*s, 0)
			blue = blue.Blend(green, 0.1*s, 0.2*s, 0)
		}
		for col := range perside {
			for row := range perside {
				cubeCell(cs, maxContrast, side, col, row, red, green, blue)
			}
		}
	}
}

func cubeCell(
	cs Table,
	maxContrast func(a, b float64) float64,
	side, col, row int,
	red, green, blue Color,
) {
	if row == 0 && col == 0 && side > 0 {
		return
	}
	idx := side*6 + row*36 + col
	if side == 0 {
		idx += 53
	} else {
		idx += 16
	}
	lightsrc := 52 + row*36 // Left column of red colors from first side.
	if side > 0 {
		lightsrc -= 36 // Move upper to left/top corner.
	}
	c := red.Blend(blue, 0.5+0.1*float64(col)-0.1*float64(row), 0.5, 0)
	light := maxContrast(
		cs.Color(lightsrc).Lightness(),
		cs.Color(16+col).Lightness(),
	)
	if side > 0 {
		light = maxContrast(light, green.Lightness())
	}
	cs.SetColor(idx, c.With(-1, -1, light))
}
