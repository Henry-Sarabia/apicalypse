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

// FuncOption is a functional option type used to set the options for an API query.
// FuncOption is the first-order function returned by the available functional options
// (e.g. Fields or Limit). For the full list of supported filters and their expected
// syntax, please visit: https://apicalypse.io/syntax/
type FuncOption func(*options) error

// Fields is a functional option for setting the included fields in the results from a query.
func Fields(fields ...string) FuncOption {
	return func(opt *options) error {
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
		opt.Filters["fields"] = f

		return nil
	}
}

// Exclude is a functional option for setting the excluded fields in the results from a query.
func Exclude(fields ...string) FuncOption {
	return func(opt *options) error {
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
		opt.Filters["exclude"] = f

		return nil
	}
}

// Where is a functional option for setting a custom data filter similar to SQL.
// If multiple filters are provided, they are AND'd together.
// For the full list of filters and more information, visit: https://apicalypse.io/syntax/
func Where(filters ...string) FuncOption {
	return func(opt *options) error {
		if len(filters) <= 0 {
			return ErrMissingInput
		}

		for _, f := range filters {
			if isBlank(f) {
				return ErrBlankArgument
			}
		}

		f := strings.Join(filters, " & ")
		opt.Filters["where"] = f

		return nil
	}
}

// Limit is a functional option for setting the number of items to return from a query.
// This usually has a maximum limit.
func Limit(n int) FuncOption {
	return func(opt *options) error {
		if n < 0 {
			return ErrNegativeInput
		}
		opt.Filters["limit"] = strconv.Itoa(n)

		return nil
	}
}

// Offset is a functional option for setting the index to start returning results from a query.
func Offset(n int) FuncOption {
	return func(opt *options) error {
		if n < 0 {
			return ErrNegativeInput
		}
		opt.Filters["offset"] = strconv.Itoa(n)

		return nil
	}
}

// Sort is a functional option for sorting the results of a query by a certain field's
// values and the use of "asc" or "desc" to sort by ascending or descending order.
func Sort(field, order string) FuncOption {
	return func(opt *options) error {
		if isBlank(field) || isBlank(order) {
			return ErrBlankArgument
		}

		opt.Filters["sort"] = field + " " + order
		return nil
	}
}

// Search is a functional option for searching for a value in a particular column of data.
// If the column is omitted, search will be performed on the default column.
func Search(column, term string) FuncOption {
	return func(opt *options) error {
		if isBlank(term) {
			return ErrBlankArgument
		}

		if !isBlank(column) {
			column = column + " "
		}

		opt.Filters["search"] = column + `"` + term + `"`
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
