package auth

import (
	"errors"
	"fmt"
	"net/http"
)

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
