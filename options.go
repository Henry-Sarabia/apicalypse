package apicalypse

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrMissingInput = errors.New("missing input parameters")
var ErrNegativeInput = errors.New("input cannot be a negative number")

// options contains the optional filters for a custom API query.
// The options type is not accessed directly but instead mutated
// using the functional options that return a FuncOption.
type options struct {
	Filters map[string]string
}

// FuncOption is a functional option type used to set the options for an API query.
// FuncOption is the first-order function returned by the available functional options
// (e.g. Fields or Limit).
type FuncOption func(*options) error

// newOpt returns a new options object mutated by the provided FuncOption arguments.
// If no FuncOption's are provided, an empty options object is returned.
func newOpt(funcOpts ...FuncOption) (*options, error) {
	opt := &options{Filters: map[string]string{}}

	for _, f := range funcOpts {
		if err := f(opt); err != nil {
			return nil, err
		}
	}

	return opt, nil
}

// ComposeOptions composes multiple functional options into a single FuncOption.
// This is primarily used to create a single functional option that can be used
// repeatedly across multiple queries.
func ComposeOptions(funcOpts ...FuncOption) FuncOption {
	return func(opt *options) error {
		for _, f := range funcOpts {
			if err := f(opt); err != nil {
				return err
			}
		}
		return nil
	}
}

// Fields is a functional option for setting the included fields in the results from a query.
func Fields(fields ...string) FuncOption {
	return func(opt *options) error {
		if len(fields) <= 0 {
			return ErrMissingInput
		}

		f := strings.Join(fields, ",")
		f = removeWhitespace(f)
		opt.Filters["fields"] = f

		return nil
	}
}

// Exclude is a functional option for setting the exluded fields in the results from a query.
func Exclude(fields ...string) FuncOption {
	return func(opt *options) error {
		if len(fields) <= 0 {
			return ErrMissingInput
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
		opt.Filters["sort"] = field + " " + order
		return nil
	}
}

// Search is a functional option for searching for a value.
func Search(term string) FuncOption {
	return func(opt *options) error {
		opt.Filters["search"] = term
		return nil
	}
}

// removeWhitespace returns the provided string with all of the whitespace removed.
// This includes spaces, tabs, newlines, returns, and form feeds.
func removeWhitespace(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, "")
}
