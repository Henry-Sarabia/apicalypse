package apicalypse

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"testing"
)

func TestComposeOptions(t *testing.T) {
	tests := []struct {
		name     string
		funcOpts []Option
	}{
		{"Zero options", []Option{}},
		{"Single option", []Option{Limit(15)}},
		{"Multiple options", []Option{Limit(15), Fields("name", "rating")}},
		{"Single error option", []Option{Limit(-99)}},
		{"Multiple error options", []Option{Limit(-99), Fields()}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			comp := ComposeOptions(test.funcOpts...)

			wantFilters, wantErr := newFilters(test.funcOpts...)
			gotFilters, gotErr := newFilters(comp)
			if !reflect.DeepEqual(errors.Cause(gotErr), errors.Cause(wantErr)) {
				t.Errorf("got: <%v>, want: <%v>", errors.Cause(gotErr), errors.Cause(wantErr))
			}
			if !reflect.DeepEqual(gotFilters, wantFilters) {
				t.Errorf("got: <%v>, want: <%v>", gotFilters, wantFilters)
			}
		})
	}
}

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
			filters, err := newFilters()
			if err != nil {
				t.Fatal(err)
			}

			err = Fields(test.fields...)(filters)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if filters["fields"] != test.wantFields {
				t.Errorf("got: <%v>, want: <%v>", filters["fields"], test.wantFields)
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
			filters, err := newFilters()
			if err != nil {
				t.Fatal(err)
			}

			err = Exclude(test.fields...)(filters)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if filters["exclude"] != test.wantFields {
				t.Errorf("got: <%v>, want: <%v>", filters["exclude"], test.wantFields)
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
			filters, err := newFilters()
			if err != nil {
				t.Fatal(err)
			}

			err = Where(test.filters...)(filters)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if filters["where"] != test.wantFilters {
				t.Errorf("got: <%v>, want: <%v>", filters["where"], test.wantFilters)
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
			filters, err := newFilters()
			if err != nil {
				t.Fatal(err)
			}

			err = Limit(test.limit)(filters)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if filters["limit"] != test.wantLimit {
				t.Errorf("got: <%v>, want: <%v>", filters["limit"], test.wantLimit)
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
			filters, err := newFilters()
			if err != nil {
				t.Fatal(err)
			}

			err = Offset(test.offset)(filters)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if filters["offset"] != test.wantOffset {
				t.Errorf("got: <%v>, want: <%v>", filters["offset"], test.wantOffset)
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
			filters, err := newFilters()
			if err != nil {
				t.Fatal(err)
			}

			err = Sort(test.field, test.order)(filters)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if filters["sort"] != test.wantSort {
				t.Errorf("got: <%v>, want: <%v>", filters["sort"], test.wantSort)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name    string
		column  string
		term    string
		want    string
		wantErr error
	}{
		{"Non-empty column and non-empty term", "name", "halo", `name "halo"`, nil},
		{"Empty column and non-empty term", "", "halo", `"halo"`, nil},
		{"Non-empty column and empty term", "name", "", "", ErrBlankArgument},
		{"Empty column and empty term", "", "", "", ErrBlankArgument},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filters, err := newFilters()
			if err != nil {
				t.Fatal(err)
			}

			err = Search(test.column, test.term)(filters)

			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("got: <%v>, want: <%v>", err, test.wantErr)
			}

			if filters["search"] != test.want {
				t.Errorf("got: <%v>, want: <%v>", filters["search"], test.want)
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
	req, _ := NewRequest("GET", "https://some-internet-game-database-api/search/", Search("", "Halo"))

	// Search for results with the name "Zelda"
	req, _ = NewRequest("GET", "https://some-internet-game-database-api/search/", Search("", "Zelda"))

	http.DefaultClient.Do(req)
}
