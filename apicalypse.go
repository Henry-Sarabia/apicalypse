package apicalypse

import (
	"github.com/pkg/errors"
	"net/http"
)

// BodyRequest returns a request configured for the provided url using the provided method.
// The provided query options are written to the body of the request.
// The default method is GET.
func BodyRequest(method string, url string, options ...FuncOption) (*http.Request, error) {
	if isBlank(url) {
		return nil, ErrBlankArgument
	}

	opt, err := newOptions(options...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create new options")
	}

	req, err := http.NewRequest(method, url, opt.reader())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create request with method '%s' for url '%s'", method, url)
	}

	return req, nil
}
