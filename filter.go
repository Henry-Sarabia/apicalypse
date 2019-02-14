package apicalypse

import (
	"github.com/pkg/errors"
	"strings"
)

// newFilters returns a filter map mutated by the provided Option arguments.
// If no Option's are provided, an empty map is returned.
func newFilters(funcOpts ...Option) (map[string]string, error) {
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
