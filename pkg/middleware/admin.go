package middleware

import (
	"log"
	"net/http"

	"github.com/StaphoneWizzoh/Go_Auth/pkg/utils"
)

func Admin(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const bearerSchema = "Bearer "

		// Getting authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Middleware error: Authorization header blank")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Checking if authorization header is valid
		if len(authHeader) < len(bearerSchema) || authHeader[:len(bearerSchema)] != bearerSchema {
			log.Println("Middleware error: Authorization header is invalid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Getting token
		token := authHeader[len(bearerSchema):]
		claims, err := utils.ParseToken(token, true)
		if err != nil {
			log.Println("Middleware error: Error in parsing token:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Checking the user role
		if claims.Role != "admin" && claims.Role != "superadmin"{
			log.Println("Middleware error: The claim isn't an administrator one: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Setting up user role in context
		ctx := utils.SetUserRoleInContext(r.Context(), claims.Role)
		
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}