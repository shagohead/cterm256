package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/shagohead/cterm256/pkg/filetype"
	_ "github.com/shagohead/cterm256/pkg/filetype/alacritty"
	_ "github.com/shagohead/cterm256/pkg/filetype/kitty"
	"github.com/shagohead/cterm256/pkg/printer"
	"github.com/shagohead/cterm256/pkg/termcolor"
)

func main() {
	if err := run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

var (
	fileName     string
	fileType     = &filetype.Flag{}
	debugColors  string
	printColors  bool
	printCurrent bool
	overwrite    bool
	skipGen      bool
	lightOutput  bool
)

const cmdMain = "cterm256"

func run() error {
	fs := flag.NewFlagSet(cmdMain, flag.ExitOnError)
	fs.Var(fileType, "t", "File type. Supported values: "+filetype.RegisteredNames())
	fs.StringVar(&fileName, "f", "", "Source colorscheme file. If omits STDIN will be used")
	fs.BoolVar(&overwrite, "w", false, "Overwrite source colorscheme file instead of writing to STDOUT")
	fs.BoolVar(&printColors, "print", false, "Print color table instead of colorscheme output")
	fs.BoolVar(&printCurrent, "print-current", false, "Print table with current terminal colors")
	fs.StringVar(&debugColors, "debug", "", "Print HSL data for specified colors (`number/b/bg/f/fg`), separetaed by comma and optionally prefixed with «-» for blank line prepending")
	fs.BoolVar(&skipGen, "skip-gen", false, "Skip color table generation")
	fs.BoolVar(&lightOutput, "light-stderr", false, "Write light/dark to STDERR")
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), `Usage: `+cmdMain+` [-h | --help]

Patch 8/16 terminal color scheme with generated 239 other ANSI colors.

`)
		fs.PrintDefaults()
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	if printCurrent {
		printer.PrintCurrent()
		return nil
	}
	ft := fileType.FileType
	if ft == nil {
		if fileName == "" {
			return errors.New("use -t or -f option")
		}
		ext := path.Ext(fileName)
		for name, ftype := range filetype.RegisteredTypes() {
			if ftype.Support(fileName, ext) {
				ft = ftype
				if !lightOutput {
					os.Stderr.WriteString("Type determined by file name: " + name + "\n")
				}
				break
			}
		}
		if ft == nil {
			return fmt.Errorf("cannot find supported file type of %s", fileName)
		}
	}
	var in io.Reader
	var file *os.File
	in = os.Stdin
	if fileName != "" {
		var err error
		ff := os.O_RDONLY
		if overwrite {
			ff = os.O_RDWR
		}
		file, err = os.OpenFile(fileName, ff, 0)
		if err != nil {
			return err
		}
		defer file.Close()
		in = file
	}
	scheme, err := ft.Parse(in)
	if err != nil {
		return err
	}
	if !skipGen {
		var warns io.Writer = os.Stderr
		if lightOutput {
			warns = new(noopWriter)
		}
		if err := termcolor.Generate(scheme, warns); err != nil {
			return err
		}
	}
	if debugColors != "" {
		for i, raw := range strings.Split(debugColors, ",") {
			if len(raw) > 0 && raw[0] == '-' {
				fmt.Fprint(os.Stderr, "\n")
				raw = raw[1:]
			}
			switch raw {
			case "b", "bg":
				fmt.Fprintf(os.Stderr, "bg: %s\n", scheme.Background())
			case "f", "fg":
				fmt.Fprintf(os.Stderr, "fg: %s\n", scheme.Foreground())
			default:
				n, err := strconv.Atoi(raw)
				if err != nil {
					return fmt.Errorf("debug flag[%d]: %v", i, err)
				}
				fmt.Fprintf(os.Stderr, "%d: %s\n", n, scheme.Color(n))
			}
		}
	}
	if printColors {
		printer.PrintScheme(scheme)
		return nil
	}
	if debugColors != "" {
		return nil
	}
	var w termcolor.Writer = os.Stdout
	if overwrite {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return err
		}
		w = file
	}
	if err := scheme.Write(w); err != nil {
		return err
	}
	if lightOutput {
		if scheme.Color(1).Lightness() > scheme.Color(16).Lightness() {
			os.Stderr.WriteString("dark")
		} else {
			os.Stderr.WriteString("light")
		}
	}
	return nil
}

type noopWriter struct{}

// Write implements io.Writer.
func (*noopWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
