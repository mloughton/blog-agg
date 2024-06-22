package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mloughton/blog-agg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (h authedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, s *Server) {
	user, err := s.GetAuthUser(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	h(w, r, user)
}

func (s *Server) middlewareAuth(h authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r, s)
	}
}

func GetAuthHeader(r *http.Request, keyName string) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < (len(keyName) + 1) {
		return "", errors.New("malformed authorization header")
	}
	if authHeader[:len(keyName)] != keyName {
		return "", fmt.Errorf("%s not found in authorization header", keyName)
	}

	return authHeader[len(keyName)+1:], nil
}

func (s *Server) GetAuthUser(r *http.Request) (database.User, error) {
	apiKey, err := GetAuthHeader(r, "ApiKey")
	if err != nil {
		return database.User{}, err
	}
	ctx := context.Background()
	user, err := s.DB.GetUserByAPIKey(ctx, apiKey)
	if err != nil {
		return database.User{}, err
	}
	return user, nil
}
