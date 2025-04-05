package termcolor

import (
	"fmt"
	"slices"
	"testing"
)

func Test_gradient(t *testing.T) {
	tests := []struct {
		name    string
		cs      *table
		isDark  bool
		brindex int
		targets [5]int
		expect  orderedMap
	}{
		{
			name: "dark-red",
			cs: &table{
				colors: map[int]hsl{
					16: {0, 20, 10},   // dark background
					9:  {12, 100, 50}, // bright red
					1:  {12, 100, 30}, // normal red
				},
			},
			isDark:  true,
			brindex: 9,
			targets: [5]int{52, 88, 124, 160, 196},
			expect: orderedMap{
				52:  {12, 28, 14},
				88:  {12, 36, 18},
				124: {12, 44, 22},
				160: {12, 52, 26},
				196: {12, 60, 30},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gradient(tt.cs, tt.isDark, tt.brindex, tt.targets)
			for _, i := range tt.expect.keys() {
				got, exist := tt.cs.colors[i]
				if !exist {
					t.Fatalf("Color #%d does not exists", i)
				}
				want := tt.expect[i]
				if g, w := fmt.Sprintf("%+v", got), fmt.Sprintf("%+v", want); g != w {
					t.Errorf("Color #%d = %s, want %s", i, g, w)
				}
			}
		})
	}
}

type orderedMap map[int]hsl

func (m orderedMap) keys() []int {
	s := make([]int, 0, len(m))
	for i := range m {
		s = append(s, i)
	}
	slices.Sort(s)
	return s
}

type table struct {
	colors map[int]hsl
	front  hsl
	back   hsl
}

// Background implements Table.
func (t *table) Background() Color {
	return t.back
}

// Color implements Table.
func (t *table) Color(number int) Color {
	return t.colors[number]
}

// Foreground implements Table.
func (t *table) Foreground() Color {
	return t.front
}

// SetColor implements Table.
func (t *table) SetColor(number int, color Color) {
	t.colors[number] = color.(hsl)
}

var _ Table = (*table)(nil)
