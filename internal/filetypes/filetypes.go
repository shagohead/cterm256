package filetypes

import (
	"errors"
	"flag"
	"strings"

	"github.com/shagohead/cterm256/internal/cterm"
)

var ftypes = make(map[string]cterm.FileType)

func Register(name string, ftype cterm.FileType) {
	if _, ok := ftypes[name]; ok {
		panic("duplicate FileType's name")
	}
	ftypes[name] = ftype
}

func RegisteredTypes() map[string]cterm.FileType {
	return ftypes
}

func RegisteredNames() string {
	s := strings.Builder{}
	n := len(ftypes)
	var i int
	for name := range ftypes {
		s.WriteString(name)
		if i++; i < n {
			s.WriteRune(' ')
		}
	}
	return s.String()
}

// FileType selector flag.
type Flag struct {
	Name     string
	FileType cterm.FileType
}

// Set implements flag.Value.
func (f *Flag) Set(val string) error {
	f.Name = val
	var ok bool
	f.FileType, ok = ftypes[val]
	if !ok {
		return errors.New("unknown filetype. supported values: " + RegisteredNames())
	}
	return nil
}

// String implements flag.Value.
func (f *Flag) String() string {
	return f.Name
}

var _ flag.Value = (*Flag)(nil)
