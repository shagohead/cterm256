package cterm

import (
	"io"

	"github.com/hsluv/hsluv-go"
)

type FileType interface {
	Parse(input io.Reader) (*ColorScheme, error)
	Write(output io.Writer, theme *ColorScheme) error
}

type ColorScheme struct {
	Indexed    [256]Color
	Custom     map[string]Color
	Background *Color
	Foreground *Color
}

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
	hex string
	hsl HSL
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
	return RGB{
		Red:   r,
		Green: g,
		Blue:  b,
	}
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
		hex: hex,
		hsl: HSL{h, s, l},
	}
}

func ColorFromHSL(hsl HSL) Color {
	return Color{hsl: hsl}
}

func ColorFromRGB(rgb RGB) Color {
	h, s, l := hsluv.HsluvFromRGB(rgb.Red, rgb.Green, rgb.Blue)
	return Color{
		hsl: HSL{h, s, l},
	}
}
