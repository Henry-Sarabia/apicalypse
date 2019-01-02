package apicalypse

import (
	"bytes"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var ErrMissingInput = errors.New("missing input parameters")
var ErrNegativeInput = errors.New("input cannot be a negative number")

// Query contains the filters for a custom API query.
type Query struct {
	Filters map[string]string
}

// NewRequest returns a request configured for the provided url and with the provided body.
func NewRequest(url string, body *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create request for url at '%s'", url)
	}

	return req, nil
}

// Fields is a functional option for setting the included fields in the results from a query.
func Fields(fields ...string) func(*Query) error {
	return func(q *Query) error {
		if len(fields) <= 0 {
			return ErrMissingInput
		}

		f := strings.Join(fields, ",")
		f = removeWhitespace(f)
		q.Filters["fields"] = f

		return nil
	}
}

// Exclude is a functional option for setting the exluded fields in the results from a query.
func Exclude(fields ...string) func(*Query) error {
	return func(q *Query) error {
		if len(fields) <= 0 {
			return ErrMissingInput
		}

		f := strings.Join(fields, ",")
		f = removeWhitespace(f)
		q.Filters["exclude"] = f

		return nil
	}
}

// Where is a functional option for setting a custom data filter similar to SQL.
// If multiple filters are provided, they are AND'd together.
// For the full list of filters and more information, visit: https://apicalypse.io/syntax/
func Where(filters ...string) func(*Query) error {
	return func(q *Query) error {
		if len(filters) <= 0 {
			return ErrMissingInput
		}

		f := strings.Join(filters, " & ")
		q.Filters["where"] = f

		return nil
	}
}

// Limit is a functional option for setting the number of items to return from a query.
// This usually has a maximum limit.
func Limit(n int) func(*Query) error {
	return func(q *Query) error {
		if n < 0 {
			return ErrNegativeInput
		}
		q.Filters["limit"] = strconv.Itoa(n)

		return nil
	}
}

// Offset is a functional option for setting the index to start returning results from a query.
func Offset(n int) func(*Query) error {
	return func(q *Query) error {
		if n < 0 {
			return ErrNegativeInput
		}
		q.Filters["offset"] = strconv.Itoa(n)

		return nil
	}
}

// Sort is a functional option for sorting the results of a query by a certain field's
// values and the use of "asc" or "desc" to sort by ascending or descending order.
func Sort(field, order string) func(*Query) error {
	return func(q *Query) error {
		q.Filters["sort"] = field + " " + order
		return nil
	}
}

// Search is a functional option for searching for a value.
func Search(term string) func(*Query) error {
	return func(q *Query) error {
		q.Filters["search"] = term
		return nil
	}
}

// removeWhitespace returns the provided string with all of the whitespace removed.
// This includes spaces, tabs, newlines, returns, and form feeds.
func removeWhitespace(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, "")
}
