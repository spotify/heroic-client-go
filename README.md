spotify/heroic-client-go
========================

Go client for talking to a [heroic](https://github.com/spotify/heroic) cluster.

[![BuildStatus Widget]][BuildStatus Result]
[![GoReport Widget]][GoReport Status]
[![GoDocWidget]][GoDocReference]

[BuildStatus Result]: https://travis-ci.org/spotify/heroic-client-go
[BuildStatus Widget]: https://travis-ci.org/spotify/heroic-client-go.svg?branch=master

[GoReport Status]: https://goreportcard.com/report/github.com/spotify/heroic-client-go
[GoReport Widget]: https://goreportcard.com/badge/github.com/spotify/heroic-client-go

[GoDocWidget]: https://godoc.org/github.com/spotify/heroic-client-go?status.svg
[GoDocReference]:https://godoc.org/github.com/spotify/heroic-client-go 

## Usage

```go
import "github.com/spotify/heroic-client-go"
```

Construct a new Heroic client and get the status of a cluster:

```go
u, _ := url.Parse("http://heroic.spotify.net/")
c := heroic.NewClient(u, nil, nil)
ctx := context.Background()
status, _ := c.Status(ctx)
fmt.Println(status.Service.Name) // "The Heroic Time Series Database"
```

## Roadmap

This library is being developed for an internal application at Spotify, so API 
methods will likely be implemented in the order that they are needed by that 
application.


## Code of Conduct

This project adheres to the [Open Code of Conduct][code-of-conduct]. By 
participating, you are expected to honor this code.


[code-of-conduct]: https://github.com/spotify/code-of-conduct/blob/master/code-of-conduct.md
