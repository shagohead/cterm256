package termcolor

import (
	"errors"
	"fmt"
	"io"

	"github.com/lucasb-eyer/go-colorful"
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

var errMissingBackground = errors.New("provided scheme missing (primary) background color")

// Generate 256 color palette based on first 8/16 + background & foreground.
func Generate(cs Table, warns io.Writer) error {
	for i := range 8 {
		if cs.Color(i).Nil() {
			return fmt.Errorf("provided scheme missing color %d", i)
		}
	}
	background := cs.Background()
	if background.Nil() {
		return errMissingBackground
	}

	// Is it dark or light theme?
	bglight := background.Lightness()
	isDark := cs.Color(1).Lightness() > bglight
	contrast := maxContrast(isDark)

	// Swap black and white colors for light theme if needed.
	if b, w := cs.Color(0), cs.Color(7); !isDark && w.Lightness() > b.Lightness() {
		cs.SetColor(0, w)
		cs.SetColor(7, b)
		if b, w = cs.Color(8), cs.Color(15); !b.Nil() && !w.Nil() {
			cs.SetColor(8, w)
			cs.SetColor(15, b)
		}
	}

	// Set 16 from 0 with background lightness.
	{
		_, a, b := cs.Color(0).src.Lab()
		l, _, _ := background.src.Lab()
		cs.SetColor(16, Color{colorful.Lab(l, a, b), true})
	}

	// Fix bright (and normal) colors: create, swap or change lightness if needed.
	tobright := makebright(bglight, isDark)
	for i := range 8 {
		norm, bright := cs.Color(i), cs.Color(i+8)
		if bright.Nil() {
			l, a, b := norm.src.Lab()
			cs.SetColor(i+8, color(tobright(l), a, b))
			continue
		}
		n, _, _ := norm.src.Lab()
		b, _, _ := bright.src.Lab()
		if n == b {
			l, a, b := bright.src.Lab()
			cs.SetColor(i+8, color(tobright(l), a, b))
			continue
		}
		// TODO: Fix schemes like Solarized, which uses bright as completly different colors.
		// Solarized bright colors are background/foreground gradient.
		// We already have 232-255 for this purpose.
		if contrast(n, b) == n {
			cs.SetColor(i, bright)
			cs.SetColor(i+8, norm)
		}
	}

	// R/g/b based gradients.
	for _, c := range []struct {
		brsource int    // Bright r/g/b
		targets  [5]int // Indexes of tinted versions
	}{
		{9, [5]int{52, 88, 124, 160, 196}}, // Red
		{10, [5]int{22, 28, 34, 40, 46}},   // Green
		{12, [5]int{17, 18, 19, 20, 21}},   // Blue
	} {
		gradient(cs, isDark, c.brsource, c.targets)
	}

	cube(cs, contrast)

	foreground := cs.Foreground()
	if foreground.Nil() {
		white := cs.Color(7)
		foreground = white
		repr := "white/7"
		if brwhite := cs.Color(15); !brwhite.Nil() {
			if contrast(brwhite.Lightness(), white.Lightness()) == brwhite.Lightness() {
				foreground = brwhite
				repr = "brwhite/15"
			}
		}
		fmt.Fprintln(warns, "scheme missing foreground color; will use ", repr, " instead")
	}
	delta := 1.0 / 25
	for i := range 24 {
		s := delta * float64(i+1)
		sl, sa, sb := background.src.Lab()
		dl, da, db := foreground.src.Lab()
		cs.SetColor(232+i, color(blend(sl, dl, s), blend(sa, da, s), blend(sb, db, s)))
	}
	return nil
}

func makebright(bg float64, dark bool) func(src float64) float64 {
	return func(src float64) float64 {
		// TODO: implement
		return src
	}
}

func blend(src, dst, scale float64) float64 {
	return src + (scale * (dst - src))
}

// maxContrast treats a and b as lightness values and returns one which have highest contrast with background.
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
		d := float64(i) * 0.2
		sl, _, _ := bg.src.Lab()
		dl, da, db := src.src.Lab()
		cs.SetColor(target, color(blend(sl, dl, d), da, db))
	}
}

// Generate 6x6x6 colors cube
func cube(cs Table, maxContrast func(a, b float64) float64) {
	for side := 1; side < 6; side++ {
		green := cs.Color(side*6 + 16)
		yl, ya, yb := green.src.Lab()

		// Blue+green columns in top row.
		for col := range 5 {
			blue := cs.Color(col + 17)
			target := side*6 + col + 17
			xl, xa, xb := blue.src.Lab()
			s := float64(side)*0.15 - float64(col)*0.05
			cs.SetColor(target, color(maxContrast(xl, yl), blend(xa, ya, s), blend(xb, yb, s)))
		}

		// Red+green rows in left column.
		for row := range 5 {
			red := cs.Color(row*36 + 52)
			target := side*6 + row*36 + 52
			xl, xa, xb := red.src.Lab()
			s := float64(side)*0.15 - float64(row)*0.05
			cs.SetColor(target, color(maxContrast(xl, yl), blend(xa, ya, s), blend(xb, yb, s)))
		}
	}

	// Mixes of top row and left column.
	for side := range 6 {
		for col := range 5 {
			blue := cs.Color(side*6 + col + 17)
			yl, ya, yb := blue.src.Lab()
			for row := range 5 {
				red := cs.Color(side*6 + row*36 + 52)
				target := side*6 + row*36 + col + 53
				xl, xa, xb := red.src.Lab()
				s := 0.5 + 0.1*float64(col) - 0.1*float64(row)
				cs.SetColor(target, color(blend(xl, yl, s), blend(xa, ya, s), blend(xb, yb, s)))
			}
		}
	}
}
