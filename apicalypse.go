package apicalypse

import (
	"bytes"
	"github.com/pkg/errors"
	"net/http"
)

// NewRequest returns a request configured for the provided url and with the provided body.
func NewRequest(url string, body *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create request for url at '%s'", url)
	}

	return req, nil
}
