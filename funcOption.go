package apicalypse

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	// ErrMissingInput occurs when a function is called without input parameters (e.g. nil slice)
	ErrMissingInput = errors.New("missing input parameters")
	// ErrBlankArgument occurs when a function is called with a blank argument that should not be blank.
	ErrBlankArgument = errors.New("a provided argument is blank or empty")
	// ErrNegativeInput occurs when a function is called with a negative number that should not be negative.
	ErrNegativeInput = errors.New("input cannot be a negative number")
)

// FuncOption is a functional option type used to set the filters for an API query.
// FuncOption is the first-order function returned by the available functional options
// (e.g. Fields or Limit). For the full list of supported filters and their expected
// syntax, please visit: https://apicalypse.io/syntax/
type FuncOption func(map[string]string) error

// ComposeOptions composes multiple functional options into a single FuncOption.
// This is primarily used to create a single functional option that can be used
// repeatedly across multiple queries.
func ComposeOptions(funcOpts ...FuncOption) FuncOption {
	return func(filters map[string]string) error {
		for _, f := range funcOpts {
			if err := f(filters); err != nil {
				return errors.Wrap(err, "cannot compose functional options")
			}
		}
		return nil
	}
}

// Fields is a functional option for setting the included fields in the results from a query.
func Fields(fields ...string) FuncOption {
	return func(filters map[string]string) error {
		if len(fields) <= 0 {
			return ErrMissingInput
		}

		for _, f := range fields {
			if isBlank(f) {
				return ErrBlankArgument
			}
		}

		f := strings.Join(fields, ",")
		f = removeWhitespace(f)
		filters["fields"] = f

		return nil
	}
}

// Exclude is a functional option for setting the excluded fields in the results from a query.
func Exclude(fields ...string) FuncOption {
	return func(filters map[string]string) error {
		if len(fields) <= 0 {
			return ErrMissingInput
		}

		for _, f := range fields {
			if isBlank(f) {
				return ErrBlankArgument
			}
		}

		f := strings.Join(fields, ",")
		f = removeWhitespace(f)
		filters["exclude"] = f

		return nil
	}
}

// Where is a functional option for setting a custom data filter similar to SQL.
// If multiple filters are provided, they are AND'd together.
// For the full list of filters and more information, visit: https://apicalypse.io/syntax/
func Where(custom ...string) FuncOption {
	return func(filters map[string]string) error {
		if len(custom) <= 0 {
			return ErrMissingInput
		}

		for _, f := range custom {
			if isBlank(f) {
				return ErrBlankArgument
			}
		}

		c := strings.Join(custom, " & ")
		filters["where"] = c

		return nil
	}
}

// Limit is a functional option for setting the number of items to return from a query.
// This usually has a maximum limit.
func Limit(n int) FuncOption {
	return func(filters map[string]string) error {
		if n < 0 {
			return ErrNegativeInput
		}
		filters["limit"] = strconv.Itoa(n)

		return nil
	}
}

// Offset is a functional option for setting the index to start returning results from a query.
func Offset(n int) FuncOption {
	return func(filters map[string]string) error {
		if n < 0 {
			return ErrNegativeInput
		}
		filters["offset"] = strconv.Itoa(n)

		return nil
	}
}

// Sort is a functional option for sorting the results of a query by a certain field's
// values and the use of "asc" or "desc" to sort by ascending or descending order.
func Sort(field, order string) FuncOption {
	return func(filters map[string]string) error {
		if isBlank(field) || isBlank(order) {
			return ErrBlankArgument
		}

		filters["sort"] = field + " " + order
		return nil
	}
}

// Search is a functional option for searching for a value in a particular column of data.
// If the column is omitted, search will be performed on the default column.
func Search(column, term string) FuncOption {
	return func(filters map[string]string) error {
		if isBlank(term) {
			return ErrBlankArgument
		}

		if !isBlank(column) {
			column = column + " "
		}

		filters["search"] = column + `"` + term + `"`
		return nil
	}
}

// removeWhitespace returns the provided string with all of the whitespace removed.
// This includes spaces, tabs, newlines, returns, and form feeds.
func removeWhitespace(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, "")
}

// isBlank returns true if the provided string is empty or only consists of whitespace.
// Returns false otherwise.
func isBlank(s string) bool {
	if removeWhitespace(s) == "" {
		return true
	}

	return false
}
