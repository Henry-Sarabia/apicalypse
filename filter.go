package apicalypse

import (
	"github.com/pkg/errors"
	"strings"
)

// newFilters returns a filter map mutated by the provided FuncOption arguments.
// If no FuncOption's are provided, an empty map is returned.
func newFilters(funcOpts ...FuncOption) (map[string]string, error) {
	filters := map[string]string{}

	for _, f := range funcOpts {
		if err := f(filters); err != nil {
			return nil, errors.Wrap(err, "cannot create new options")
		}
	}

	return filters, nil
}

// toString returns the filters as a single string.
func toString(f map[string]string) string {
	if len(f) <= 0 {
		return ""
	}

	b := strings.Builder{}
	for k, v := range f {
		b.WriteString(k + " " + v + "; ")
	}

	return b.String()
}

// reader returns the filters as a *strings.Reader
// to satisfy the io.Reader interface.
func toReader(f map[string]string) *strings.Reader {
	s := toString(f)
	return strings.NewReader(s)
}
