package generator

import (
	"github.com/shagohead/cterm256/internal/cterm"
)

func Generate(scheme *cterm.ColorScheme) error {
	scheme.Indexed[16] = scheme.Indexed[0]
	if scheme.Background != nil {
		scheme.Indexed[16] = scheme.Indexed[16].WithLigntness(scheme.Background.HSL().Lightness)
	}
	g := &generator{scheme: scheme}
	if err := g.gradient(9, [5]int{196, 160, 124, 88, 52}); err != nil {
		return err
	}
	if err := g.gradient(12, [5]int{21, 20, 19, 18, 17}); err != nil {
		return err
	}
	if err := g.gradient(10, [5]int{46, 40, 34, 28, 22}); err != nil {
		return err
	}
	for side := 0; side < 5; side++ {
		topLeft := side*6 + 22
		for i := 0; i < 5; i++ {
			g.scheme.Indexed[topLeft+i+1] = g.mix(topLeft, 17+i)
			left := 36*i + 52
			g.scheme.Indexed[topLeft+left-16] = g.mix(topLeft, left)
		}
	}
	for side := 0; side < 6; side++ {
		for row := 0; row < 5; row++ {
			rowc := side*6 + row*36 + 52
			for col := 0; col < 5; col++ {
				g.scheme.Indexed[side*6+row*36+col+53] = g.mix(rowc, side*6+col+17)
			}
		}
	}
	c := g.scheme.Indexed[16].WithLigntnessMult(1.05)
	step := (g.scheme.Indexed[7].HSL().Lightness - g.scheme.Indexed[0].HSL().Lightness) / 24
	for i := 0; i < 24; i++ {
		g.scheme.Indexed[232+i] = c.WithLigntnessDelta(step * float64(i))
	}
	return nil
}

type generator struct {
	scheme *cterm.ColorScheme
}

func (g *generator) gradient(src int, targets [5]int) error {
	start := g.scheme.Indexed[src].HSL().Lightness
	end := g.scheme.Indexed[16].HSL().Lightness
	rang := end - start
	step := rang / 4
	if rang > 0 {
		start *= 1.05
	} else {
		start *= 0.95
	}
	for i, target := range targets {
		c := g.scheme.Indexed[src]
		if i > 0 {
			c = c.WithLigntness(start + step*float64(i))
		}
		g.scheme.Indexed[target] = c
	}
	return nil
}

func (g *generator) idx(i int) cterm.Color {
	return g.scheme.Indexed[i]
}

func (g *generator) mix(a, b int) cterm.Color {
	ca := g.scheme.Indexed[a].RGB()
	cb := g.scheme.Indexed[b].RGB()
	return cterm.ColorFromRGB(cterm.RGB{
		Red:   (ca.Red + cb.Red) / 2,
		Green: (ca.Green + cb.Green) / 2,
		Blue:  (ca.Blue + cb.Blue) / 2,
	})
}
