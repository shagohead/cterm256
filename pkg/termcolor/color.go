package termcolor

import (
	"fmt"
	"regexp"

	"github.com/lucasb-eyer/go-colorful"
)

var HEX *regexp.Regexp

func init() {
	HEX = regexp.MustCompile(`^#?[0-9a-fA-F]{6}$`)
}

type Color struct {
	src colorful.Color
	set bool
}

func (h Color) Nil() bool {
	return !h.set
}

func (h Color) String() string {
	l, a, b := h.src.Lab()
	return fmt.Sprintf("%s Lab: %0.2f %0.2f %0.2f", h.HEX(), l, a, b)
}

func (h Color) HEX() string {
	return h.src.Clamped().Hex()
}

func (h Color) RGB() (uint8, uint8, uint8) {
	return h.src.Clamped().RGB255()
}

func (h Color) Lightness() float64 {
	l, _, _ := h.src.Lab()
	return l
}

func FromHEX(hex string) Color {
	c, err := colorful.Hex(hex)
	if err != nil {
		panic(fmt.Sprintf("parsing %s: %v", hex, err))
	}
	return Color{c, true}
}

func color(l, a, b float64) Color {
	return Color{colorful.Lab(l, a, b), true}
}
