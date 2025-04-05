package filetype

import (
	"errors"
	"flag"
	"io"
	"strings"

	"github.com/shagohead/cterm256/pkg/termcolor"
)

type FileType interface {
	Parse(input io.Reader) (termcolor.Table, error)
	Support(name, ext string) bool
}

var ftypes = make(map[string]FileType)

func Register(name string, ftype FileType) {
	if _, ok := ftypes[name]; ok {
		panic("duplicate FileType name: " + name)
	}
	ftypes[name] = ftype
}

func RegisteredTypes() map[string]FileType {
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
	FileType FileType
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
