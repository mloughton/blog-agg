# Blog Aggregator

A simple implementation of a http api using golang and postgreSQL.

Completed via guided project on [Boot.dev](https://www.boot.dev)

## Implementation

* Go standard library [http package](https://pkg.go.dev/net/http) used for server and routes
* [Goose](https://pkg.go.dev/github.com/pressly/goose/v3) migration tool for performing database migrations
* [SQLc](https://pkg.go.dev/github.com/kyleconroy/sqlc) used for compiling SQL queries into go code

## Installation

```bash
go get github.com/mloughton/blog-agg
```