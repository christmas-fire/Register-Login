package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/christmas-fire/register-login/internal/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := jwt.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", claims["username"])
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
