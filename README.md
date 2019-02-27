# Apicalypse

[![GoDoc](https://godoc.org/github.com/Henry-Sarabia/apicalypse?status.svg)](https://godoc.org/github.com/Henry-Sarabia/apicalypse) [![Build Status](https://travis-ci.com/Henry-Sarabia/apicalypse.svg?branch=master)](https://travis-ci.com/Henry-Sarabia/apicalypse) [![Coverage Status](https://coveralls.io/repos/github/Henry-Sarabia/apicalypse/badge.svg?branch=master)](https://coveralls.io/github/Henry-Sarabia/apicalypse?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/henry-sarabia/apicalypse)](https://goreportcard.com/report/github.com/henry-sarabia/apicalypse)

Create custom [Apicalypse](https://apicalypse.io/) compliant queries with the `apicalypse`
[Go](https://golang.org/) package. With the `apicalypse` package, you can create queries with
any of the supported filters including: Fields, Exclude, Where, Limit, Offset, Sort, and Search. 
The package will be kept up to date with Apicalypse - if more filters are introduced, they will
be supported here as well. To see the complete list of filters and their respective syntax please
visit the Apicalypse syntax page [here](https://apicalypse.io/syntax/).

If you would like to lend a hand with the package, please feel free to submit a pull request.
Any and all contributions are welcome and appreciated!

## Installation

If you do not have Go installed yet, you can find installation instructions 
[here](https://golang.org/doc/install).

To pull the most recent version of `apicalypse`, use `go get`.

```
go get github.com/Henry-Sarabia/apicalypse
```

Then import the package into your project as you normally would.

```go
import "github.com/Henry-Sarabia/apicalypse"
```

## Usage

### Creating A New Query

Creating a new Apicalypse query is simple - use the provided `Query()` function along with whatever
filters you need. In our example, we need a query that will limit the number of results to 25
and retrieve an object's "name" and "age" fields.

```go
qry, err := apicalypse.Query(Limit(25), Fields("name", "age"))
if err != nil {
	// handle error
}
```

With that, we have a query ready to be sent to your favorite Apicalypse-powered API. This is sent
primarily in the body of a GET request. Of course, this package makes that easy.

### Creating A New Request

The `apicalypse` package also provides a convenient `NewRequest()` function to get you up and 
running. In this example, we need a request fulfulling the same requirements as the previous 
example.

```go
req, err := apicalypse.NewRequest(
	"GET", 
	"https://myapi.com/actors", 
	Limit(25), 
	Fields("name", "age"),
	)
if err != nil {
	// handle error
}
```

Just like that, we have a request configured with the necessary filters. If we had any headers 
that need their own configuration, this would be the time to do it. If not, the request can be sent
off straight away.

You may provide any method and URL to the function but be aware that the Apicalypse specifications
recommend a simple GET request for the majority of cases. POST and PUT requests may be used
under certain circumstances.

For more information, please visit the Apicalypse implementation
page [here](https://apicalypse.io/implementation/).

### Functional Options

The `apicalypse` package uses functional options to apply the different query filters to
a request. More specifically, functional options are first-order functions that are passed 
to the `NewRequest()` function. It's much simpler than it sounds and the package makes it
easy!

Let's walk through a few different functional option examples to demonstrate their
ease of use.

To specify what fields to return from an API query, pass `Fields()` to the `NewRequest()`
function.
```go
req, err := apicalypse.NewRequest("GET", "https://myapi.com/actors", Fields("name", "movies", "age"))
```
With that, our new request is configured to only fetch the "name", "movies", and "age" fields
from our friends at the totally, absolutely real myapi.com.

To specify the limit of results we want returned from an API query, pass `Limit()` to the `NewRequest()`
function.
```go
req, err := apicalypse.NewRequest("GET", "https://myapi.com/actors", Limit(15))
```
It's no more difficult than that - we're ready to fetch up to 15 results!

Although this next example may initially look more complicated, it is only a matter of adhering
to the proper Apicalypse syntax. Again, that can be found [here](https://apicalypse.io/syntax/).

To specify a filter for our results from an API query, pass `Where()` to the `NewRequest()`
function.
```go
req, err := apicalypse.NewRequest("GET", "https://myapi.com/actors", Where("age > 50 & movies != null"))
```
Our new request is now configured to filter the results so only the results which have an
age above 50 and a non-null movies field are returned.

The remaining functional options are no more complicated than the examples presented here.
Moreover, they are further described in the [documentation](https://godoc.org/github.com/Henry-Sarabia/apicalypse#FuncOption).

### Functional Option Composition

More often than not you will need to set multiple options for an API query.
Fortunately, this functionality is supported through variadic functions and
functional option composition.

First, the `NewRequest()` and `Query()` functions are variadic which means you can pass in any 
number of functional options.
```go
req, err := apicalypse.NewRequest(
	"GET",
	"https://myapi.com/actors",
	Fields("name", "movies", "age"),
	Where("age > 50 & movies != null"),
	Limit(15),
	)
```
Leveraging the variadic nature of our primary functions, we have quickly and simply
created a request with several different filters automatically applied to it.

This request is now configured to return results with the fields "name", "movies", and "age".
In addition, the results will be filtered so only results with age above 50 and a non-null 
movies field are returned. Finally, only up to 15 results will be returned.

Second, the `apicalypse` package provides a `ComposeOptions()` function which takes any number
of functional options and composes them into a single, custom made, ready-to-use functional option.
```go
myOpt := apicalypse.ComposeOptions(
	Fields("name", "movies", "age"),
        Where("age > 50 & movies != null"),
        Limit(15),
	)
```
This call to `ComposeOptions()` creates a single functional option that performs the same
filters we've been using up until now. The major difference is that now we only need to provide
this new single functional option to any new requests that require those specific filters.
```go
req, err := apicalypse.NewRequest("GET", "https://myapi.com/actors", myOpt)
```
Of course, you can still pass in additional functional options if need be.

For example, now that we've gathered our 15 results, perhaps we want the next 15. It's as simple
as passing in one more functional option.
```go
req, err := apicalypse.NewRequest("GET", "https://myapi.com/actors", myOpt, Offset(15))
```
Our new request is ready to return the next 15 results!

Functional option composition reduces duplicate code and helps keep your code
DRY. You can even compose newly composed functional options for even more
finely grained control over similar queries.

## Examples

The repository contains a few examples that demonstrate how one could use the `apicalypse`
package. The examples can be found in their respective test files or be read conveniently on the
GoDoc reference found [here](https://godoc.org/github.com/Henry-Sarabia/apicalypse#pkg-examples).

If you have used the `apicalypse` package for a project and would like to have it featured
here as a reference for new users, please submit an issue or pull request and I'll be sure to
add it. Thank you!

## Contributions

If you would like to contribute to this project, please adhere to the following
guidelines.

* Submit an issue describing the problem.
* Fork the repo and add your contribution.
* Add appropriate tests.
* Run go fmt, go vet, and golint.
* Prefer idiomatic Go over non-idiomatic code.
* Follow the basic Go conventions found [here](https://github.com/golang/go/wiki/CodeReviewComments).
* If in doubt, try to match your code to the current codebase.
* Create a pull request with a description of your changes.

I'll review pull requests as they come in and merge them if everything checks out.

Again, any contribution is greatly appreciated!

## Special Thanks

* The [IGDB](https://github.com/igdb) team who developed Apicalypse with an accompanying [Node client](https://github.com/igdb/node-apicalypse)
* The helpful community over on the IGDB [Discord server](https://discord.gg/pXn8Jh9) 
* The insightful Dave Cheney for his [article](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
on functional options