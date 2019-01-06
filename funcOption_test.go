package apicalypse

import (
	"fmt"
	"net/http"
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

func ExampleComposeOptions() {
	// Composing FuncOptions to filter out unpopular results
	composedOpts := ComposeOptions(
		Fields("title", "username", "game", "likes", "content"),
		Where("likes > 10"),
		Where("views >= 200"),
		Limit(25),
	)

	// Using composed FuncOptions
	req1, err := NewRequest("GET", "https://some-internet-game-database-api/games/", composedOpts)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Reusing composed FuncOptions
	req2, err := NewRequest("GET", "https://some-internet-game-database-api/games/", composedOpts, Offset(25))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Retrieves first set of 15 popular games
	http.DefaultClient.Do(req1)
	// Retrieves second set of 15 popular games
	http.DefaultClient.Do(req2)
}

func ExampleFields() {
	// Retrieve games with name field
	req, _ := NewRequest("GET", "https://some-internet-game-database-api/games/", Fields("name"))

	// Retrieve games with genres field
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Fields("genres"))

	// Retrieve games with both name and genres field
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Fields("name", "genres"))

	// Retrieve games with a parent game
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Fields("parent_game"))

	// Retrieve games with a parent game name
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Fields("parent_game.name"))

	// Retrieve games with any number of fields
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Fields("name", "genres", "popularity", "rating"))

	// Retrieve games with all available fields
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Fields("*"))

	// Execute latest request
	http.DefaultClient.Do(req)
}

func ExampleExclude() {
	// Retrieve games without name field
	req, _ := NewRequest("GET", "https://some-internet-game-database-api/games/", Exclude("name"))

	// Retrieve games without genres field
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Exclude("genres"))

	// Retrieve games without name or genres field
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Exclude("name", "genres"))

	// Retrieve games without parent game
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Exclude("parent_game"))

	// Retrieve games without parent game name
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Exclude("parent_game.name"))

	// Retrieve games without any number of fields
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Exclude("name", "genres", "popularity", "rating"))

	// Execute latest request
	http.DefaultClient.Do(req)
}

func ExampleWhere() {
	// Retrieve games with a rating equal to 50
	req, _ := NewRequest("GET", "https://some-internet-game-database-api/games/", Where("rating = 50"))

	// Retrieve games with a rating not equal to 50
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Where("rating != 50"))

	// Retrieve games with a rating greater than 50
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Where("rating > 50"))

	// Retrieve games with a rating greater than or equal to 50
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Where("rating >= 50"))

	// Retrieve games with a rating less than 50
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Where("rating < 50"))

	// Retrieve games with a rating less than or equal to 50
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Where("rating <= 50"))

	// Retrieve games with all of the following genres: Roleplaying, Adventure, and MMO
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Where("genres = [Roleplaying, Adventure, MMO]"))

	// Retrieve games without all of the following genres: genres Roleplaying, Adventure, and MMO
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Where("genres != [Roleplaying, Adventure, MMO]"))

	// Retrieve games with at least one of the following genres: Roleplaying, Adventure, or MMO
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Where("genres = (Roleplaying, Adventure, MMO)"))

	// Retrieve games without any of the following genres: Roleplaying, Adventure, or MMO
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Where("genres != (Roleplaying, Adventure, MMO)"))

	// Retrieve games with exclusively the following genres: Roleplaying, Adventure, or MMO
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Where("genres = {Roleplaying, Adventure, MMO}"))

	// Retrieve games with a rating greater than 50 AND with at least one of the following genres: Roleplaying, Adventure, or MMO
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Where("rating > 50", "genres = (Roleplaying, Adventure, MMO)"))

	// Execute latest request
	http.DefaultClient.Do(req)
}

func ExampleLimit() {
	// Retrieve up to 1 result
	req, _ := NewRequest("GET", "https://some-internet-game-database-api/games/", Limit(1))

	// Retrieve up to 25 results
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Limit(25))

	// Retrieve up to 50 results
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Limit(50))

	// Execute latest request
	http.DefaultClient.Do(req)
}

func ExampleOffset() {
	// Retrieve the first batch of 10 results
	req, _ := NewRequest("GET", "https://some-internet-game-database-api/games/", Limit(10))

	// Retrieve the second batch of 10 results
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Limit(10), Offset(10))

	// Retrieve the third batch of 10 results
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Limit(10), Offset(20))

	// Retrieve the fourth batch of 10 results
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Limit(10), Offset(30))

	http.DefaultClient.Do(req)
}

func ExampleSort() {
	// Retrieve the most popular games
	req, _ := NewRequest("GET", "https://some-internet-game-database-api/games/", Sort("popularity", "desc"))

	// Retrieve the least popular games
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Sort("popularity", "asc"))

	// Retrieve the earliest released games by their first release date
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Sort("first_release_date", "asc"))

	// Retrieve games with the latest release date
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/games/", Sort("first_release_date", "desc"))

	http.DefaultClient.Do(req)
}

func ExampleSearch() {
	// Search for results with the name "Halo"
	req, _ := NewRequest("GET", "https://some-internet-game-database-api/search/", Search("Halo"))

	// Search for results with the name "Zelda"
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/search/", Search("Zelda"))

	http.DefaultClient.Do(req)
}
