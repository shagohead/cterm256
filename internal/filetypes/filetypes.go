package filetypes

import (
	"errors"
	"flag"

	"github.com/shagohead/cterm256/internal/cterm"
	"github.com/shagohead/cterm256/internal/filetypes/kitty"
)

// FileType selector flag.
type Flag struct {
	Name     string
	FileType cterm.FileType
}

const Supported = "kitty"

// Set implements flag.Value.
func (f *Flag) Set(val string) error {
	f.Name = val
	switch val {
	case "kitty":
		f.FileType = &kitty.FileType{}
	default:
		return errors.New("unknown filetype. supported values: " + Supported)
	}
	return nil
}

// String implements flag.Value.
func (f *Flag) String() string {
	return f.Name
}

var _ flag.Value = (*Flag)(nil)
