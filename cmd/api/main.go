package main

import (
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mloughton/blog-agg/internal/feeds"
	"github.com/mloughton/blog-agg/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	httpServer, internalServer, err := server.NewServer()
	if err != nil {
		panic(err)
	}
	go feeds.StartScraping(internalServer.DB, 10, time.Duration(60*time.Second))
	httpServer.ListenAndServe()
}
