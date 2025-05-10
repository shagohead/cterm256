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
		{input: "#FFFFFF", valid: true},
		{input: "#gggggg", valid: false},
	} {
		t.Run(tt.input, func(t *testing.T) {
			if v := HEX.MatchString(tt.input); v != tt.valid {
				t.Errorf("HEX.MatchString() = %v, want %v", v, tt.valid)
			}
		})
	}
}
