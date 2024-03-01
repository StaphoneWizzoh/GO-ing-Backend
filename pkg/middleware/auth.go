package middleware

import (
	"net/http"

	"github.com/StaphoneWizzoh/Go_Auth/pkg/utils"
)


func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const bearerSchema = "Bearer "

		// Getting authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Checking if authorization header is valid
		if len(authHeader) < len(bearerSchema) || authHeader[:len(bearerSchema)] != bearerSchema {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Getting token
		token := authHeader[len(bearerSchema):]
		claims, err := utils.ParseToken(token, true)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Setting up user ID in context
		ctx := utils.SetUserIdInContext(r.Context(), claims.UserID)
		
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}