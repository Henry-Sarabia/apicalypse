package apicalypse

import (
	"bytes"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
)

// ErrFilterOverlap occurs when a single-use option is set multiple times in a single query.
var ErrFilterOverlap = errors.New("filter already set")

// Query contains the filters for a custom API query.
type Query struct {
	Values url.Values
}

// NewRequest returns a request configured for the provided url and with the provided body.
func NewRequest(url string, body *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create request for url at '%s'", url)
	}

	return req, nil
}

// Limit is a funcitonal option for setting the limit filter on the provided query.
func Limit(n int) func(*Query) error {
	return func(q *Query) error {
		if q.Values.Get("limit") != "" {
			return ErrFilterOverlap
		}

		q.Values.Set("limit", strconv.Itoa(n))
		return nil
	}
}
