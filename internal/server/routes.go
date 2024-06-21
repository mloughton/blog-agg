package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mloughton/blog-agg/internal/auth"
	"github.com/mloughton/blog-agg/internal/database"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", s.GetHealthHandler)
	mux.HandleFunc("GET /v1/err", s.GetErrorHandler)

	mux.HandleFunc("POST /v1/users", s.PostUserHandler)
	mux.HandleFunc("GET /v1/users", s.GetUsersHandler)
	return mux
}

func (s *Server) GetHealthHandler(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (s *Server) GetErrorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (s *Server) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	newUUID, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	userParams := database.CreateUserParams{
		ID:        newUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
	}
	ctx := context.Background()
	user, err := s.DB.CreateUser(ctx, userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (s *Server) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAuthHeader(r, "ApiKey")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := context.Background()
	user, err := s.DB.GetUserByAPIKey(ctx, apiKey)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found")
		return
	}
	respondWithJSON(w, http.StatusOK, user)

}
