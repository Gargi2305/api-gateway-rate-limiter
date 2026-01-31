package middleware

import (
	"context"
	"net/http"
	"strings"

	"api-gateway/internal/metrics"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			metrics.GatewayErrorsTotal.
				WithLabelValues(r.URL.Path, "401").
				Inc()

			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Expect format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			metrics.GatewayErrorsTotal.
				WithLabelValues(r.URL.Path, "401").
				Inc()

			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse & validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure token method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte("my-secret-key"), nil
		})

		if err != nil || !token.Valid {
			metrics.GatewayErrorsTotal.
				WithLabelValues(r.URL.Path, "401").
				Inc()

			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			metrics.GatewayErrorsTotal.
				WithLabelValues(r.URL.Path, "401").
				Inc()

			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			metrics.GatewayErrorsTotal.
				WithLabelValues(r.URL.Path, "401").
				Inc()

			http.Error(w, "invalid user id in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		r = r.WithContext(ctx)

		// Token is valid → allow request
		next.ServeHTTP(w, r)
	})
}
