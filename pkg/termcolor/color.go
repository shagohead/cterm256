package termcolor

import (
	"regexp"

	"github.com/hsluv/hsluv-go"
)

var HEX *regexp.Regexp

func init() {
	HEX = regexp.MustCompile(`^#?[0-9a-fA-F]{6}$`)
}

type Color interface {
	HEX() string
	RGB() (int, int, int)

	Hue() float64
	Saturation() float64
	Lightness() float64

	// Blend current color with another.
	// Parameter values set transition steps from 0 to 1.
	// I.e. values which equals zero returns color with unchanged property.
	Blend(with Color, hue, sat, light float64) Color

	// Return color modified with argument values which less or eq zero.
	// I.e. values below zero ignored and returns as original color.
	With(hue, sat, light float64) Color
}

func FromHEX(hex string) Color {
	h, s, l := hsluv.HsluvFromHex(hex)
	return hsl{h, s, l}
}

type hsl struct {
	h float64
	s float64
	l float64
}

func zerodist(v float64) float64 {
	if v > 180 {
		return 360 - v
	}
	return v
}

// Blend implements Color.
func (h hsl) Blend(with Color, hue float64, sat float64, light float64) Color {
	if hue > 0 {
		azimuth := with.Hue() - h.h
		if azimuth < 0 {
			azimuth *= -1
		}
		through0 := zerodist(h.h) + zerodist(with.Hue())
		if through0 < azimuth {
			delta := through0 * hue
			if h.h < 180 {
				delta *= -1
			}
			h.h += delta
			switch {
			case h.h < 0:
				h.h += 360
			case h.h > 360:
				h.h -= 360
			}
		} else {
			h.h += hue * (with.Hue() - h.Hue())
		}
	}
	if sat > 0 {
		h.s += sat * (with.Saturation() - h.s)
	}
	if light > 0 {
		h.l += light * (with.Lightness() - h.l)
	}
	return h
}

// HEX implements Color.
func (h hsl) HEX() string {
	return hsluv.HsluvToHex(h.h, h.s, h.l)
}

// RGB implements Color.
func (h hsl) RGB() (int, int, int) {
	r, g, b := hsluv.HsluvToRGB(h.h, h.s, h.l)
	return int(r * 255), int(g * 255), int(b * 255)
}

// Hue implements Color.
func (h hsl) Hue() float64 {
	return h.h
}

// Saturation implements Color.
func (h hsl) Saturation() float64 {
	return h.s
}

// Lightness implements Color.
func (h hsl) Lightness() float64 {
	return h.l
}

// With implements Color.
func (h hsl) With(hue float64, sat float64, light float64) Color {
	if hue >= 0 {
		h.h = hue
	}
	if sat >= 0 {
		h.s = sat
	}
	if light >= 0 {
		h.l = light
	}
	return h
}
