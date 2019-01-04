package apicalypse

import (
	"reflect"
	"testing"
)

func TestFields(t *testing.T) {
	tests := []struct {
		name      string
		fields    []string
		expFields string
		expErr    error
	}{
		{"Single non-empty field", []string{"name"}, "name", nil},
		{"Multiple non-empty fields", []string{"name", "popularity", "rating"}, "name,popularity,rating", nil},
		{"Empty fields slice", []string{}, "", ErrMissingInput},
		{"Single empty field", []string{"  "}, "", ErrBlankField},
		{"Multiple empty fields", []string{"", " ", "", ""}, "", ErrBlankField},
		{"Mixed empty and non-empty fields", []string{"", "id", "  ", "url"}, "", ErrBlankField},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Fields(test.fields...)(opt)

			if !reflect.DeepEqual(err, test.expErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.expErr)
			}

			if opt.Filters["fields"] != test.expFields {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["fields"], test.expFields)
			}
		})
	}
}

func TestExclude(t *testing.T) {
	tests := []struct {
		name      string
		fields    []string
		expFields string
		expErr    error
	}{
		{"Single non-empty field", []string{"name"}, "name", nil},
		{"Multiple non-empty fields", []string{"name", "popularity", "rating"}, "name,popularity,rating", nil},
		{"Empty fields slice", []string{}, "", ErrMissingInput},
		{"Single empty field", []string{"  "}, "", ErrBlankField},
		{"Multiple empty fields", []string{"", " ", "", ""}, "", ErrBlankField},
		{"Mixed empty and non-empty fields", []string{"", "id", "  ", "url"}, "", ErrBlankField},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Exclude(test.fields...)(opt)

			if !reflect.DeepEqual(err, test.expErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.expErr)
			}

			if opt.Filters["exclude"] != test.expFields {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["exclude"], test.expFields)
			}
		})
	}
}

func TestWhere(t *testing.T) {
	tests := []struct {
		name       string
		filters    []string
		expFilters string
		expErr     error
	}{
		{"Single non-empty filter", []string{"b.count >= 14"}, "b.count >= 14", nil},
		{"Multiple non-empty filters", []string{"b.count >= 14", "a != n"}, "b.count >= 14 & a != n", nil},
		{"Empty filters slice", []string{}, "", ErrMissingInput},
		{"Single empty filter", []string{" "}, "", ErrBlankField},
		{"Multiple empty filters", []string{"", " ", "  "}, "", ErrBlankField},
		{"Mixed empty and non-empty filters", []string{"b.count >= 14", "", "a != n", " "}, "", ErrBlankField},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Where(test.filters...)(opt)

			if !reflect.DeepEqual(err, test.expErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.expErr)
			}

			if opt.Filters["where"] != test.expFilters {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["where"], test.expFilters)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		expLimit string
		expErr   error
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

			if !reflect.DeepEqual(err, test.expErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.expErr)
			}

			if opt.Filters["limit"] != test.expLimit {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["limit"], test.expLimit)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	tests := []struct {
		name      string
		offset    int
		expOffset string
		expErr    error
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

			if !reflect.DeepEqual(err, test.expErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.expErr)
			}

			if opt.Filters["offset"] != test.expOffset {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["offset"], test.expOffset)
			}
		})
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		order   string
		expSort string
		expErr  error
	}{
		{"Non-empty field and non-empty order", "b.count", "desc", "b.count desc", nil},
		{"Non-empty field and empty order", "b.count", " ", "", ErrBlankField},
		{"Empty field and non-empty order", "", "desc", "", ErrBlankField},
		{"Empty field and empty order", "", "", "", ErrBlankField},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Sort(test.field, test.order)(opt)

			if !reflect.DeepEqual(err, test.expErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.expErr)
			}

			if opt.Filters["sort"] != test.expSort {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["sort"], test.expSort)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name    string
		term    string
		expTerm string
		expErr  error
	}{
		{"Non-empty term", "halo", "halo", nil},
		{"Empty term", "", "", ErrBlankField},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opt, err := newOptions()
			if err != nil {
				t.Fatal(err)
			}

			err = Search(test.term)(opt)

			if !reflect.DeepEqual(err, test.expErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.expErr)
			}

			if opt.Filters["search"] != test.expTerm {
				t.Errorf("got: <%v>, want: <%v>", opt.Filters["search"], test.expTerm)
			}
		})
	}
}

func TestRemoveWhitespace(t *testing.T) {
	tests := []struct {
		name string
		s    string
		exp  string
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
			s := removeWhitespace(test.s)

			if s != test.exp {
				t.Errorf("got: <%v>, want: <%v>", s, test.exp)
			}
		})
	}
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name string
		s    string
		exp  bool
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
			b := isBlank(test.s)

			if b != test.exp {
				t.Errorf("got: <%v>, want: <%v>", b, test.exp)
			}
		})
	}
}
