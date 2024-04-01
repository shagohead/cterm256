package printer

import (
	"fmt"

	"github.com/shagohead/cterm256/internal/cterm"
)

func printTable(fn func(n int) string) {
	fmt.Printf("Standard%32s\n", "Bright")
	for n := 0; n < 16; n++ {
		if n == 8 {
			fmt.Print("\033[0m  ")
		}
		fmt.Printf("\033[48;%sm %02d ", fn(n), n)
	}
	fmt.Print("\033[0m\n\n")

	fmt.Println("216 colors 6x6x6 cube")
	var side, row, col int
	for {
		n := side*6 + row*36 + col + 16
		fmt.Printf("\033[48;%vm %03d ", fn(n), n)
		col++
		if col == 6 {
			side++
			col = 0
			if side%3 == 0 {
				row++
				if row < 6 {
					side -= 3
				} else {
					row = 0
					fmt.Print("\033[0m\n")
				}
			}
			if side == 6 {
				break
			}
			if side%3 != 0 {
				fmt.Print("\033[0m  ")
			} else {
				fmt.Print("\033[0m\n")
			}
		}
	}
	fmt.Print("\033[0m\nGrayscale\n")
	for n := 232; n < 256; n++ {
		if n == 244 {
			fmt.Print("\033[0m\n")
		}
		fmt.Printf("\033[48;%sm %02d ", fn(n), n)
	}
	fmt.Print("\033[0m\n\n")
}

func PrintScheme(scheme *cterm.ColorScheme) {
	printTable(func(n int) string {
		rgb := scheme.Indexed[n].RGB()
		r, g, b := int(rgb.Red*255), int(rgb.Green*255), int(rgb.Blue*255)
		return fmt.Sprintf("2;%v;%v;%v", r, g, b)
	})
}

func PrintCurrent() {
	printTable(func(n int) string {
		return fmt.Sprintf("5;%v", n)
	})
}
