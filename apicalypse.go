package apicalypse

import (
	"github.com/pkg/errors"
	"net/http"
)

// NewRequest returns a request configured for the provided url using the provided method.
// The provided query options are written to the body of the request. The default method is GET.
func NewRequest(method string, url string, options ...FuncOption) (*http.Request, error) {
	if isBlank(url) {
		return nil, ErrBlankArgument
	}

	filters, err := newFilters(options...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create new options")
	}

	req, err := http.NewRequest(method, url, toReader(filters))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create request with method '%s' for url '%s'", method, url)
	}

	return req, nil
}
