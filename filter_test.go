package apicalypse

import (
	"bytes"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"testing"
)

func TestNewFilters(t *testing.T) {
	tests := []struct {
		name        string
		funcOpts    []FuncOption
		wantFilters map[string]string
		wantErr     error
	}{
		{"Empty option", []FuncOption{}, map[string]string{}, nil},
		{"Single option", []FuncOption{Limit(15)}, map[string]string{"limit": "15"}, nil},
		{"Multiple options", []FuncOption{Limit(15), Offset(10), Fields("name", "rating")}, map[string]string{"limit": "15", "offset": "10", "fields": "name,rating"}, nil},
		{"Single error option", []FuncOption{Limit(-99)}, nil, ErrNegativeInput},
		{"Multiple error options", []FuncOption{Fields(), Exclude(), Where()}, nil, ErrMissingInput},
		{"Mixed options", []FuncOption{Limit(10), Offset(-99)}, nil, ErrNegativeInput},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filters, err := newFilters(test.funcOpts...)
			if !reflect.DeepEqual(errors.Cause(err), test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if !reflect.DeepEqual(filters, test.wantFilters) {
				t.Errorf("got: <%v>, want: <%v>", filters, test.wantFilters)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name    string
		filters map[string]string
		wants   []string
	}{
		{"Zero filters", map[string]string{}, nil},
		{"Single filter", map[string]string{"limit": "15"}, []string{"limit 15; "}},
		{"Multiple filters", map[string]string{"limit": "15", "fields": "id,name,rating"}, []string{"limit 15; ", "fields id,name,rating; "}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := toString(test.filters)

			for _, want := range test.wants {
				if !strings.Contains(got, want) {
					t.Errorf("got: <%v>, want: <%v>", got, want)
				}
			}
		})
	}
}

func TestToReader(t *testing.T) {
	tests := []struct {
		name    string
		filters map[string]string
		wants   []string
	}{
		{"Zero filters", map[string]string{}, nil},
		{"Single filter", map[string]string{"limit": "15"}, []string{"limit 15; "}},
		{"Multiple filters", map[string]string{"limit": "15", "fields": "id,name,rating"}, []string{"limit 15; ", "fields id,name,rating; "}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			buf.ReadFrom(toReader(test.filters))
			got := buf.String()

			for _, want := range test.wants {
				if !strings.Contains(got, want) {
					t.Errorf("got: <%v>, want: <%v>", got, want)
				}
			}
		})
	}
}
