package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mloughton/blog-agg/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	server, err := server.NewServer()
	if err != nil {
		panic(err)
	}
	server.ListenAndServe()
}
