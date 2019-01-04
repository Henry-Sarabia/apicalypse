package apicalypse

import (
	"reflect"
	"testing"
)

func TestFields(t *testing.T) {
	tests := []struct {
		name       string
		fields     []string
		wantFields string
		wantErr    error
	}{
		{"Single non-empty field", []string{"name"}, "name", nil},
		{"Multiple non-empty fields", []string{"name", "popularity", "rating"}, "name,popularity,rating", nil},
		{"Empty fields slice", []string{}, "", ErrMissingInput},
		{"Single empty field", []string{"  "}, "", ErrBlankArgument},
		{"Multiple empty fields", []string{"", " ", "", ""}, "", ErrBlankArgument},
		{"Mixed empty and non-empty fields", []string{"", "id", "  ", "url"}, "", ErrBlankArgument},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Fields(test.fields...)(opt)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if opt.Filters["fields"] != test.wantFields {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["fields"], test.wantFields)
			}
		})
	}
}

func TestExclude(t *testing.T) {
	tests := []struct {
		name       string
		fields     []string
		wantFields string
		wantErr    error
	}{
		{"Single non-empty field", []string{"name"}, "name", nil},
		{"Multiple non-empty fields", []string{"name", "popularity", "rating"}, "name,popularity,rating", nil},
		{"Empty fields slice", []string{}, "", ErrMissingInput},
		{"Single empty field", []string{"  "}, "", ErrBlankArgument},
		{"Multiple empty fields", []string{"", " ", "", ""}, "", ErrBlankArgument},
		{"Mixed empty and non-empty fields", []string{"", "id", "  ", "url"}, "", ErrBlankArgument},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Exclude(test.fields...)(opt)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if opt.Filters["exclude"] != test.wantFields {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["exclude"], test.wantFields)
			}
		})
	}
}

func TestWhere(t *testing.T) {
	tests := []struct {
		name        string
		filters     []string
		wantFilters string
		wantErr     error
	}{
		{"Single non-empty filter", []string{"b.count >= 14"}, "b.count >= 14", nil},
		{"Multiple non-empty filters", []string{"b.count >= 14", "a != n"}, "b.count >= 14 & a != n", nil},
		{"Empty filters slice", []string{}, "", ErrMissingInput},
		{"Single empty filter", []string{" "}, "", ErrBlankArgument},
		{"Multiple empty filters", []string{"", " ", "  "}, "", ErrBlankArgument},
		{"Mixed empty and non-empty filters", []string{"b.count >= 14", "", "a != n", " "}, "", ErrBlankArgument},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Where(test.filters...)(opt)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if opt.Filters["where"] != test.wantFilters {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["where"], test.wantFilters)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	tests := []struct {
		name      string
		limit     int
		wantLimit string
		wantErr   error
	}{
		{"Positive limit", 25, "25", nil},
		{"Zero limit", 0, "0", nil},
		{"Negative limit", -25, "", ErrNegativeInput},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Limit(test.limit)(opt)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if opt.Filters["limit"] != test.wantLimit {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["limit"], test.wantLimit)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	tests := []struct {
		name       string
		offset     int
		wantOffset string
		wantErr    error
	}{
		{"Positive offset", 10, "10", nil},
		{"Zero offset", 0, "0", nil},
		{"Negative offset", -10, "", ErrNegativeInput},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Offset(test.offset)(opt)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if opt.Filters["offset"] != test.wantOffset {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["offset"], test.wantOffset)
			}
		})
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		order    string
		wantSort string
		wantErr  error
	}{
		{"Non-empty field and non-empty order", "b.count", "desc", "b.count desc", nil},
		{"Non-empty field and empty order", "b.count", " ", "", ErrBlankArgument},
		{"Empty field and non-empty order", "", "desc", "", ErrBlankArgument},
		{"Empty field and empty order", "", "", "", ErrBlankArgument},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Sort(test.field, test.order)(opt)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if opt.Filters["sort"] != test.wantSort {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["sort"], test.wantSort)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name     string
		term     string
		wantTerm string
		wantErr  error
	}{
		{"Non-empty term", "halo", "halo", nil},
		{"Empty term", "", "", ErrBlankArgument},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Search(test.term)(opt)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if opt.Filters["search"] != test.wantTerm {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["search"], test.wantTerm)
			}
		})
	}
}

func TestRemoveWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Empty", "", ""},
		{"Space", " ", ""},
		{"Tab", "\t", ""},
		{"Newline", "\n", ""},
		{"Return", "\r", ""},
		{"Mixed", " \t\n\r", ""},
		{"Non-blank with space", "one two", "onetwo"},
		{"Non-blank with tab", "one \t two", "onetwo"},
		{"Non-blank with newline", "one\ntwo", "onetwo"},
		{"Non-blank with return", "one\rtwo", "onetwo"},
		{"Non-blank with mixed", "one \t\n\r two", "onetwo"},
		{"Non-blank with surrounding space", " onetwo ", "onetwo"},
		{"Non-blank with surrounding tab", "\t one\ttwo \t", "onetwo"},
		{"Non-blank with surrounding newline", "\none\ntwo\n", "onetwo"},
		{"Non-blank with surrounding return", "\rone\rtwo\r", "onetwo"},
		{"Non-blank with surrounding mixed", "\t\n\r onetwo \t\n\r", "onetwo"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := removeWhitespace(test.input)

			if s != test.want {
				t.Errorf("got: <%v>, want: <%v>", s, test.want)
			}
		})
	}
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty", "", true},
		{"Space", " ", true},
		{"Tab", "\t", true},
		{"Newline", "\n", true},
		{"Return", "\r", true},
		{"Mixed", " \t\n\r", true},
		{"Non-blank with space", "one two", false},
		{"Non-blank with tab", "one \t two", false},
		{"Non-blank with newline", "one\ntwo", false},
		{"Non-blank with return", "one\rtwo", false},
		{"Non-blank with mixed", "one \t\n\r two", false},
		{"Non-blank with surrounding space", " onetwo ", false},
		{"Non-blank with surrounding tab", "\t one\ttwo \t", false},
		{"Non-blank with surrounding newline", "\none\ntwo\n", false},
		{"Non-blank with surrounding return", "\rone\rtwo\r", false},
		{"Non-blank with surrounding mixed", "\t\n\r onetwo \t\n\r", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := isBlank(test.input)

			if b != test.want {
				t.Errorf("got: <%v>, want: <%v>", b, test.want)
			}
		})
	}
}
