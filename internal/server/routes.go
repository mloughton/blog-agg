package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mloughton/blog-agg/internal/database"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", s.GetHealthHandler)
	mux.HandleFunc("GET /v1/err", s.GetErrorHandler)

	mux.HandleFunc("POST /v1/users", s.PostUserHandler)
	mux.HandleFunc("GET /v1/users", s.middlewareAuth(s.GetUsersHandler))

	mux.HandleFunc("POST /v1/feeds", s.middlewareAuth(s.PostFeedsHandler))
	mux.HandleFunc("GET /v1/feeds", s.GetFeedsHandler)

	mux.HandleFunc("POST /v1/feed_follows", s.middlewareAuth(s.PostFeedFollowsHandler))
	mux.HandleFunc("GET /v1/feed_follows", s.middlewareAuth(s.GetFeedFollowsHandler))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", s.DeleteFeedFollowsHandler)
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

func (s *Server) GetUsersHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	respondWithJSON(w, http.StatusOK, u)
}

func (s *Server) PostFeedsHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	var req struct {
		Name string `json="name"`
		URL  string `json="url"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	newFeedUUID, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	feedParams := database.CreateFeedParams{
		ID:        newFeedUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
		Url:       req.URL,
		UserID:    u.ID,
	}
	ctx := context.Background()
	feed, err := s.DB.CreateFeed(ctx, feedParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	newFollowUUID, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	followParams := database.CreateFeedFollowParams{
		ID:        newFollowUUID,
		FeedID:    feed.ID,
		UserID:    u.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	feedFollow, err := s.DB.CreateFeedFollow(ctx, followParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	res := struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}{
		Feed:       feed,
		FeedFollow: feedFollow,
	}
	respondWithJSON(w, http.StatusCreated, res)
}

func (s *Server) GetFeedsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	feeds, err := s.DB.GetFeedsAll(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	}
	respondWithJSON(w, http.StatusOK, feeds)
}

func (s *Server) PostFeedFollowsHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	var req struct {
		FeedID uuid.UUID `json:"feed_id"`
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
	params := database.CreateFeedFollowParams{
		ID:        newUUID,
		FeedID:    req.FeedID,
		UserID:    u.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ctx := context.Background()
	feed, err := s.DB.CreateFeedFollow(ctx, params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	}
	respondWithJSON(w, http.StatusCreated, feed)
}

func (s *Server) GetFeedFollowsHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	ctx := context.Background()
	feedFollows, err := s.DB.GetFeedFollowsAll(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	}
	respondWithJSON(w, http.StatusOK, feedFollows)
}

func (s *Server) DeleteFeedFollowsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("feedFollowID")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Bad Request")
	}
	ctx := context.Background()
	err := s.DB.DeleteFeedFollow(ctx, uuid.UUID([]byte(id)))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusNoContent)
}
