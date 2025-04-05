package termcolor

import "testing"

func Test_HEX(t *testing.T) {
	for _, tt := range []struct {
		input string
		valid bool
	}{
		{input: "#000000", valid: true},
		{input: "#999999", valid: true},
		{input: "#123456", valid: true},
		{input: "000000", valid: true},
		{input: "00000", valid: false},
		{input: "0000000", valid: false},
		{input: "#ffffff", valid: true},
		{input: "#gggggg", valid: false},
	} {
		t.Run(tt.input, func(t *testing.T) {
			if v := HEX.MatchString(tt.input); v != tt.valid {
				t.Errorf("HEX.MatchString() = %v, want %v", v, tt.valid)
			}
		})
	}
}

func Test_hsl_Blend(t *testing.T) {
	tests := []struct {
		name  string
		from  hsl
		with  hsl
		hue   float64
		sat   float64
		light float64
		want  hsl
	}{
		{
			name: "hue/azimuth",
			from: hsl{h: 0},
			with: hsl{h: 180},
			hue:  0.5,
			want: hsl{h: 90},
		},
		{
			name: "hue/reverse-azimuth",
			from: hsl{h: 330},
			with: hsl{h: 180},
			hue:  0.5,
			want: hsl{h: 255},
		},
		{
			name: "hue/through-zero",
			from: hsl{h: 330},
			with: hsl{h: 20},
			hue:  0.9,
			want: hsl{h: 15},
		},
		{
			name: "hue/reverse-through-zero",
			from: hsl{h: 20},
			with: hsl{h: 330},
			hue:  0.9,
			want: hsl{h: 335},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.from.Blend(tt.with, tt.hue, tt.sat, tt.light)
			if got != tt.want {
				t.Errorf("Blend() = %v, want %v", got, tt.want)
			}
		})
	}
}
