package apicalypse

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestNewOptions(t *testing.T) {
	tests := []struct {
		name     string
		funcOpts []FuncOption
		wantOpts *options
		wantErr  error
	}{
		{"Empty option", []FuncOption{}, &options{Filters: map[string]string{}}, nil},
		{"Single option", []FuncOption{Limit(15)}, &options{map[string]string{"limit": "15"}}, nil},
		{"Multiple options", []FuncOption{Limit(15), Offset(10), Fields("name", "rating")}, &options{map[string]string{"limit": "15", "offset": "10", "fields": "name,rating"}}, nil},
		{"Single error option", []FuncOption{Limit(-99)}, nil, ErrNegativeInput},
		{"Multiple error options", []FuncOption{Fields(), Exclude(), Where()}, nil, ErrMissingInput},
		{"Mixed options", []FuncOption{Limit(10), Offset(-99)}, nil, ErrNegativeInput},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions(test.funcOpts...)
			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if !reflect.DeepEqual(opt, test.wantOpts) {
				t.Errorf("got: <%v>, want: <%v>", opt, test.wantOpts)
			}
		})
	}
}

func TestComposeOptions(t *testing.T) {
	tests := []struct {
		name     string
		funcOpts []FuncOption
	}{
		{"Zero options", []FuncOption{}},
		{"Single option", []FuncOption{Limit(15)}},
		{"Multiple options", []FuncOption{Limit(15), Fields("name", "rating")}},
		{"Single error option", []FuncOption{Limit(-99)}},
		{"Multiple error options", []FuncOption{Limit(-99), Fields()}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			comp := ComposeOptions(test.funcOpts...)

			wantOpt, wantErr := newOptions(test.funcOpts...)
			gotOpt, gotErr := newOptions(comp)
			if !reflect.DeepEqual(gotErr, wantErr) {
				t.Errorf("got: <%v>, want: <%v>", gotErr, wantErr)
			}
			if !reflect.DeepEqual(gotOpt, wantOpt) {
				t.Errorf("got: <%v>, want: <%v>", gotOpt, wantOpt)
			}
		})
	}
}

func TestOptions_String(t *testing.T) {
	tests := []struct {
		name  string
		opt   *options
		wants []string
	}{
		{"Zero filters", &options{Filters: map[string]string{}}, nil},
		{"Single filter", &options{Filters: map[string]string{"limit": "15"}}, []string{"limit 15; "}},
		{"Multiple filters", &options{Filters: map[string]string{"limit": "15", "fields": "id,name,rating"}}, []string{"limit 15; ", "fields id,name,rating; "}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.opt.string()

			for _, want := range test.wants {
				if !strings.Contains(got, want) {
					t.Errorf("got: <%v>, want: <%v>", got, want)
				}
			}
		})
	}
}

func TestOptions_Reader(t *testing.T) {
	tests := []struct {
		name  string
		opt   *options
		wants []string
	}{
		{"Zero filters", &options{Filters: map[string]string{}}, nil},
		{"Single filter", &options{Filters: map[string]string{"limit": "15"}}, []string{"limit 15; "}},
		{"Multiple filters", &options{Filters: map[string]string{"limit": "15", "fields": "id,name,rating"}}, []string{"limit 15; ", "fields id,name,rating; "}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			buf.ReadFrom(test.opt.reader())
			got := buf.String()

			for _, want := range test.wants {
				if !strings.Contains(got, want) {
					t.Errorf("got: <%v>, want: <%v>", got, want)
				}
			}
		})
	}
}
