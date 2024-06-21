package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mloughton/blog-agg/internal/database"
)

type Server struct {
	DB *database.Queries
}

func NewServer() (*http.Server, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("couldn't find port in .env")
	}

	dbURL := os.Getenv("DATABASE_CONN")
	if dbURL == "" {
		return nil, errors.New("couldn't find database connection string in .env")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	NewServer := &Server{
		DB: database.New(db),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", port),
		Handler: NewServer.RegisterRoutes(),
	}
	return server, nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	responseJSON, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responseJSON)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type response struct {
		Error string `json:"error"`
	}
	responseBody := response{
		Error: msg,
	}
	respondWithJSON(w, code, responseBody)
}
