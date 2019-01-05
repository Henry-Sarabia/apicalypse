package apicalypse

import (
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		url         string
		opts        []FuncOption
		wantRequest *http.Request
		wantErr     error
	}{
		{"GET method, non-empty url, zero options", "GET", "http://fake.com/", nil, httptest.NewRequest("GET", "http://fake.com/", nil), nil},
		{"GET method, non-empty url, single option", "GET", "http://fake.com/", []FuncOption{Limit(15)}, httptest.NewRequest("GET", "http://fake.com/", strings.NewReader("limit 15; ")), nil},
		{"GET method, non-empty url, error option", "GET", "http://fake.com/", []FuncOption{Limit(-99)}, httptest.NewRequest("GET", "http://fake.com/", nil), ErrNegativeInput},
		{"GET method, empty url, zero options", "GET", "", nil, nil, ErrBlankArgument},
		{"GET method, empty url, single option", "GET", "", []FuncOption{Limit(15)}, nil, ErrBlankArgument},
		{"GET method, empty url, error option", "GET", "", []FuncOption{Limit(-99)}, nil, ErrBlankArgument},
		{"POST method, non-empty url, zero options", "POST", "http://fake.com/", nil, httptest.NewRequest("POST", "http://fake.com/", nil), nil},
		{"POST method, non-empty url, single option", "POST", "http://fake.com/", []FuncOption{Limit(15)}, httptest.NewRequest("POST", "http://fake.com/", strings.NewReader("limit 15; ")), nil},
		{"POST method, non-empty url, error option", "POST", "http://fake.com/", []FuncOption{Limit(-99)}, httptest.NewRequest("POST", "http://fake.com/", nil), ErrNegativeInput},
		{"POST method, empty url, zero options", "POST", "", nil, nil, ErrBlankArgument},
		{"POST method, empty url, single option", "POST", "", []FuncOption{Limit(15)}, nil, ErrBlankArgument},
		{"POST method, empty url, error option", "POST", "", []FuncOption{Limit(-99)}, nil, ErrBlankArgument},
		{"Empty method, non-empty url, zero options", "", "http://fake.com/", nil, httptest.NewRequest("GET", "http://fake.com/", nil), nil},
		{"Empty method, non-empty url, single option", "", "http://fake.com/", []FuncOption{Limit(15)}, httptest.NewRequest("GET", "http://fake.com/", strings.NewReader("limit 15; ")), nil},
		{"Empty method, non-empty url, error option", "", "http://fake.com/", []FuncOption{Limit(-99)}, httptest.NewRequest("GET", "http://fake.com/", nil), ErrNegativeInput},
		{"Empty method, empty url, zero options", "", "", nil, nil, ErrBlankArgument},
		{"Empty method, empty url, single option", "", "", []FuncOption{Limit(15)}, nil, ErrBlankArgument},
		{"Empty method, empty url, error option", "", "", []FuncOption{Limit(-99)}, nil, ErrBlankArgument},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := NewRequest(test.method, test.url, test.opts...)
			if !reflect.DeepEqual(errors.Cause(err), test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if test.wantErr != nil {
				return
			}

			if req.Method != test.wantRequest.Method {
				t.Errorf("got: <%v>, want: <%v>", req.Method, test.wantRequest.Method)
			}

			if req.URL.String() != test.wantRequest.URL.String() {
				t.Errorf("got: <%v>, want: <%v>", req.URL.String(), test.wantRequest.URL.String())
			}

			if !reflect.DeepEqual(req.Body, test.wantRequest.Body) {
				t.Errorf("got: <%v>, want: <%v>", req.Body, test.wantRequest.Body)
			}

		})
	}
}
