package alacritty

import (
	"os"
	"strings"
	"testing"

	"github.com/shagohead/cterm256/pkg/termcolor"
)

func TestParse(t *testing.T) {
	expected := func(colors map[int]string) func(t *testing.T, cs termcolor.Table) {
		return func(t *testing.T, cs termcolor.Table) {
			for idx, want := range colors {
				if got := strings.ToLower(cs.Color(idx).HEX()); got != want {
					t.Errorf("Color %d.HEX() = %s, want %s", idx, got, want)
				}
			}
		}
	}

	for _, tt := range []struct {
		name string
		data string
		want func(t *testing.T, cs termcolor.Table)
	}{
		{
			name: "colors.indexed_colors as empty array",
			data: `
			[colors]
			indexed_colors = []
			`,
		},
		{
			name: "colors.indexed_colors as array of tables. ver1",
			data: `
			[colors]
			indexed_colors = [
				{ index = 0, color = "#51576D" },
				{ index = 1, color = "#E78284" }
			]
			`,
			want: expected(map[int]string{0: "#51576d", 1: "#e78284"}),
		},
		{
			name: "colors.indexed_colors as array of tables. ver2",
			data: `
			[colors]
			[[colors.indexed_colors]]
			color = '#51576D'
			index = 0

			[[colors.indexed_colors]]
			color = '#E78284'
			index = 1
			`,
			want: expected(map[int]string{0: "#51576d", 1: "#e78284"}),
		},
		{
			name: "foreground/background",
			data: `
			[colors.primary]
			background = "#303446"
			bright_foreground = "#eeeeee"
			foreground = "#C6D0F5"
			`,
			want: func(t *testing.T, cs termcolor.Table) {
				for _, want := range []struct {
					name  string
					color termcolor.Color
					hex   string
				}{
					{name: "background", color: cs.Background(), hex: "#303446"},
					{name: "foreground", color: cs.Foreground(), hex: "#c6d0f5"},
				} {
					if want.color.Nil() {
						t.Fatalf("%s is nil", want.name)
					}
					if got := want.color.HEX(); got != want.hex {
						t.Errorf("%s.HEX() = %s, want %s", want.name, got, want.hex)
					}
				}
			},
		},
		{
			name: "colors.normal",
			data: `
			[colors.normal]
			black = "#51576D"
			blue = "#8CAAEE"
			cyan = "#81C8BE"
			green = "#A6D189"
			magenta = "#F4B8E4"
			red = "#E78284"
			white = "#B5BFE2"
			yellow = "#E5C890"
			`,
			want: expected(map[int]string{
				0: "#51576d",
				4: "#8caaee",
				6: "#81c8be",
				2: "#a6d189",
				5: "#f4b8e4",
				1: "#e78284",
				7: "#b5bfe2",
				3: "#e5c890",
			}),
		},
		{
			name: "colors.bright",
			data: `
			[colors.bright]
			black = "#626880"
			blue = "#8CAAEE"
			cyan = "#81C8BE"
			green = "#A6D189"
			magenta = "#F4B8E4"
			red = "#E78284"
			white = "#A5ADCE"
			yellow = "#E5C890"
			`,
			want: expected(map[int]string{
				8:  "#626880",
				12: "#8caaee",
				14: "#81c8be",
				10: "#a6d189",
				13: "#f4b8e4",
				9:  "#e78284",
				15: "#a5adce",
				11: "#e5c890",
			}),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cs, err := new(fileType).Parse(strings.NewReader(tt.data))
			if err != nil {
				t.Error("Parse():", err)
			}
			if tt.want != nil {
				tt.want(t, cs)
			}
		})
	}
}

func TestParseExampleFile(t *testing.T) {
	in, err := os.Open("testdata/alacritty.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		in.Close()
	})
	if _, err = new(fileType).Parse(in); err != nil {
		t.Fatal("Parse() fails:", err)
	}
}
