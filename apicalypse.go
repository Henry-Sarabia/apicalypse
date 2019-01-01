package apicalypse

import (
	"bytes"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

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

// Fields
func Fields(fields ...string) func(*Query) error {
	return func(q *Query) error {
		f := strings.Join(fields, ",")
		f = removeWhitespace(f)
		q.Filters["fields"] = f
		return nil
	}
}

//Exclude

//Where

// Limit is a functional option for setting the limit filter on the provided query.
func Limit(n int) func(*Query) error {
	return func(q *Query) error {
		q.Filters["limit"] = strconv.Itoa(n)
		return nil
	}
}

// Offset
func Offset(n int) func(*Query) error {
	return func(q *Query) error {
		q.Filters["offset"] = strconv.Itoa(n)
		return nil
	}
}

// Sort

//Search

// removeWhitespace returns the provided string with all of the whitespace removed.
// This includes spaces, tabs, newlines, returns, and form feeds.
func removeWhitespace(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, "")
}
