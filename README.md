# Apicalypse

[![GoDoc](https://godoc.org/github.com/Henry-Sarabia/apicalypse?status.svg)](https://godoc.org/github.com/Henry-Sarabia/apicalypse) [![Build Status](https://travis-ci.com/Henry-Sarabia/apicalypse.svg?branch=master)](https://travis-ci.com/Henry-Sarabia/apicalypse) [![Coverage Status](https://coveralls.io/repos/github/Henry-Sarabia/apicalypse/badge.svg?branch=master)](https://coveralls.io/github/Henry-Sarabia/apicalypse?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/henry-sarabia/apicalypse)](https://goreportcard.com/report/github.com/henry-sarabia/apicalypse)

Create custom [Apicalypse](https://apicalypse.io/) compliant queries with the `apicalypse`
[Go](https://golang.org/) package. With the `apicalypse` package, you can create queries with
any of the supported filters including: Fields, Exclude, Where, Limit, Offset, Sort, and Search. 
The package will be kept up to date with Apicalypse - if more filters are introduced, they will
be supported here as well. To see the complete list of filters and their respective syntax please
visit the Apicalypse syntax page [here](https://apicalypse.io/syntax/).

If you would like to lend a hand with the package, please feel free to submit a pull request.
Any and all contributions are greatly appreciated.

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

Now you're ready to Go.

## Usage

### Creating A New Request

Creating a new request is as straightforward as you would imagine - simple use the 
`NewRequest()` function as follows.

```go
req, err := apicalypse.NewRequest("GET", "https://myapi.com/actors")
```

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
req, err := apicalypse.NewRequest("GET", "https://myapi.com/actors, Limit(15))
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
