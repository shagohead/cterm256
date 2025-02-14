package cterm

import (
	"io"

	"github.com/hsluv/hsluv-go"
)

// FIXME: Fix invalid RGB values

// Terminal emulator file type interface.
type FileType interface {
	Parse(input io.Reader) (*ColorScheme, error)
	Write(output io.Writer, theme *ColorScheme) error
	Support(name, ext string) bool
}

type ColorScheme struct {
	Indexed    [256]Color
	Background *Color
	Foreground *Color

	// Any colorscheme data specific for file format, which pass through from reading to writing.
	PassThrough any
}

// Formerly HSL, but actually these HSL is HSLuv.
type HSL struct {
	Hue        float64
	Saturation float64
	Lightness  float64
}

type RGB struct {
	Red   float64
	Green float64
	Blue  float64
}

type Color struct {
	hex   string
	hsl   HSL
	Valid bool
}

func (c Color) HEX() string {
	if c.hex != "" {
		return "#" + c.hex
	}
	return hsluv.HsluvToHex(c.hsl.Hue, c.hsl.Saturation, c.hsl.Lightness)
}

func (c Color) HSL() HSL {
	return c.hsl
}

func (c Color) RGB() RGB {
	r, g, b := hsluv.HsluvToRGB(c.hsl.Hue, c.hsl.Saturation, c.hsl.Lightness)
	return RGB{Red: r, Green: g, Blue: b}
}

// Blend color to c2 with t as transition within 0..1 range.
func (c Color) Blend(c2 Color, t float64) Color {
	return Color{
		hsl: HSL{
			Hue:        Blend(c.hsl.Hue, c2.hsl.Hue, t),
			Saturation: Blend(c.hsl.Saturation, c2.hsl.Saturation, t),
			Lightness:  Blend(c.hsl.Lightness, c2.hsl.Lightness, t),
		},
		Valid: true,
	}
}

func Blend(a, b, t float64) float64 {
	return a + t*(b-a)
}

func BlendBiDirectional(a, b, t, m float64) float64 {
	straight := b - a
	if straight < 0 {
		straight *= -1
	}
	if through0 := ShortestTo0(a, m) + ShortestTo0(b, m); through0 < straight {
		delta := through0 * t
		if a < m/2 {
			delta *= -1
		}
		a += delta
		switch {
		case a < 0:
			return m + a
		case a > m:
			return a - m
		default:
		}
		return a
	}
	return Blend(a, b, t)
}

func ShortestTo0(a, m float64) float64 {
	if a < m/2 {
		return a
	}
	return m - a
}

func (c Color) WithLigntness(l float64) Color {
	c.hex = ""
	c.hsl.Lightness = l
	return c
}

func (c Color) WithLigntnessDelta(d float64) Color {
	c.hex = ""
	c.hsl.Lightness += d
	return c
}

func (c Color) WithLigntnessMult(m float64) Color {
	c.hex = ""
	c.hsl.Lightness *= m
	return c
}

func ColorFromHEX(hex string) Color {
	h, s, l := hsluv.HsluvFromHex(hex)
	return Color{
		hex:   hex,
		hsl:   HSL{h, s, l},
		Valid: true,
	}
}

func ColorFromHSL(hsl HSL) Color {
	return Color{hsl: hsl, Valid: true}
}

func ColorFromRGB(rgb RGB) Color {
	h, s, l := hsluv.HsluvFromRGB(rgb.Red, rgb.Green, rgb.Blue)
	return Color{hsl: HSL{h, s, l}, Valid: true}
}
