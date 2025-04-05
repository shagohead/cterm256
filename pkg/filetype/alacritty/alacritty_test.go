package alacritty

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/shagohead/cterm256/pkg/termcolor"
)

func TestParse(t *testing.T) {
	expected := func(colors map[int]string) func(t *testing.T, cs termcolor.Table) {
		return func(t *testing.T, cs termcolor.Table) {
			for idx, want := range colors {
				if got := cs.Color(idx).HEX(); got != want {
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
			want: expected(map[int]string{0: "#51576D", 1: "#E78284"}),
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
			want: expected(map[int]string{0: "#51576D", 1: "#E78284"}),
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
					{name: "foreground", color: cs.Foreground(), hex: "#C6D0F5"},
				} {
					if want.color == nil {
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
				0: "#51576D",
				1: "#8CAAEE",
				2: "#81C8BE",
				3: "#A6D189",
				4: "#F4B8E4",
				5: "#E78284",
				6: "#B5BFE2",
				7: "#E5C890",
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
				9:  "#8CAAEE",
				10: "#81C8BE",
				11: "#A6D189",
				12: "#F4B8E4",
				13: "#E78284",
				14: "#A5ADCE",
				15: "#E5C890",
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
