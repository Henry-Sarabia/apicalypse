package apicalypse

import (
	"github.com/Henry-Sarabia/whitespace"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

// Query processes the provided functional options into an Apicalypse compliant
// format and returns it as a string. The string is ready to be written into the
// body of an HTTP Request.
func Query(opts ...Option) (string, error) {
	for _, opt := range opts {
		if opt == nil {
			return "", errors.New("a provided option is nil")
		}
	}
	filters, err := newFilters(opts...)
	if err != nil {
		return "", errors.Wrap(err, "cannot create new filter map")
	}

	return toString(filters), nil
}

// NewRequest returns a request configured for the provided url using the provided method.
// The provided query options are written to the body of the request. The default method is GET.
func NewRequest(method string, url string, opts ...Option) (*http.Request, error) {
	if whitespace.IsBlank(url) {
		return nil, ErrBlankArgument
	}

	q, err := Query(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create a query")
	}

	req, err := http.NewRequest(method, url, strings.NewReader(q))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create request with method '%s' for url '%s'", method, url)
	}

	return req, nil
}
