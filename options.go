package apicalypse

import (
	"bytes"
	"net/url"
	"strings"
)

// options contains the optional filters for a custom API query.
// The options type is not accessed directly but instead mutated
// using the functional options that return a FuncOption.
type options struct {
	Filters map[string]string
}

// newOptions returns a new options object mutated by the provided FuncOption arguments.
// If no FuncOption's are provided, an empty options object is returned.
func newOptions(funcOpts ...FuncOption) (*options, error) {
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

// encode returns the options' filters as a URL encoded string.
func (o *options) encode() string {
	if len(o.Filters) <= 0 {
		return ""
	}
	b := strings.Builder{}
	for k, v := range o.Filters {
		b.WriteString(k + " " + v + "; ")
	}

	return url.PathEscape(b.String())
}

// buffer returns the options' filters as a *bytes.Buffer.
func (o *options) buffer() *bytes.Buffer {
	if len(o.Filters) <= 0 {
		return nil
	}

	b := &bytes.Buffer{}
	for k, v := range o.Filters {
		b.WriteString(k + " " + v + "; ")
	}

	return b
}
