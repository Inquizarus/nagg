package httptools

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func ExtractBearerTokenFromRequest(r *http.Request) (string, error) {

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 {
		return "", errors.New("invalid Authorization header format")
	}

	if parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization type: %s", parts[0])
	}

	return parts[1], nil
}
