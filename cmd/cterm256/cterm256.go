package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/shagohead/cterm256/internal/config"
	"github.com/shagohead/cterm256/internal/filetypes"
	"github.com/shagohead/cterm256/internal/filetypes/kitty"
	"github.com/shagohead/cterm256/internal/generator"
	"github.com/shagohead/cterm256/internal/printer"
)

func main() {
	if err := run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

var (
	fileName     string
	fileType     = &filetypes.Flag{}
	configure    bool
	debugOutput  bool
	printColors  bool
	printCurrent bool
	overwrite    bool
	skipGen      bool
)

func run() error {
	fs := flag.NewFlagSet("cterm256", flag.ExitOnError)
	fs.Var(fileType, "t", "File type. Supported values: "+filetypes.Supported)
	fs.StringVar(&fileName, "f", "", "Source colorscheme file. If omits STDIN will be used")
	fs.BoolVar(&overwrite, "w", false, "Overwrite source colorscheme file")
	fs.BoolVar(&printColors, "print", false, "Print color table instead of colorscheme output")
	fs.BoolVar(&printCurrent, "print-current", false, "Print table with current terminal colors")
	fs.BoolVar(&debugOutput, "debug", false, "Print debug information")
	fs.BoolVar(&skipGen, "skip-gen", false, "Skip color table generation")
	fs.BoolVar(&configure, "configure", false, "Configure apps color schemes")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	if printCurrent {
		printer.PrintCurrent()
		return nil
	}
	if configure {
		return config.Configure()
	}
	ft := fileType.FileType
	if ft == nil {
		ft = &kitty.FileType{}
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
		if err := generator.Generate(scheme); err != nil {
			return err
		}
	}
	if debugOutput {
		fmt.Printf("%+v\n", scheme)
	}
	if printColors {
		printer.PrintScheme(scheme)
		return nil
	}
	var w io.Writer = os.Stdout
	if overwrite {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return err
		}
		w = file
	}
	if err := ft.Write(w, scheme); err != nil {
		return err
	}
	return nil
}
