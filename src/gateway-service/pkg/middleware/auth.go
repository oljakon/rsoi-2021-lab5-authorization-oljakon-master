package middleware

import (
	"net/http"
	"strings"

	"rsoi2/src/gateway-service/pkg/jwt"
)

type Auth struct {
	JwtTokenVerifier jwt.JwtTokenVerifier
}

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/authorization" || r.URL.Path == "/api/v1/callback" {
			next.ServeHTTP(w, r)
		}
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ok, err := a.JwtTokenVerifier.ValidateToken(r.Context(), bearerToken[1])
		if !ok && err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
