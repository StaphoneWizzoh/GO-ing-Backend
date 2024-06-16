package middleware

import (
	"log"
	"net/http"

	"github.com/StaphoneWizzoh/Go_Auth/pkg/utils"
	"github.com/jub0bs/cors"
)

// CreateCORSMiddleware is a function that creates a CORS middleware with specified configurations
func CreateCORSMiddleware() (*cors.Middleware, error) {
    // Create CORS middleware with specified configurations
    corsMw, err := cors.NewMiddleware(cors.Config{
        Origins:        []string{"http://localhost:5173"},
        Methods:        []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
        RequestHeaders: []string{"Authorization", "Content-Type"},
        ResponseHeaders: []string{"X-Response-Time"},
        MaxAgeInSeconds: 600,
        Credentialed:    true,
    })
    if err != nil {
        return nil, err
    }

    corsMw.SetDebug(true) // turn debug mode on (optional)

    return corsMw, nil
}

// Admin middleware
func CorsAdmin(next http.Handler) http.Handler {
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
        if claims.Role != "admin" && claims.Role != "superadmin" {
            log.Println("Middleware error: The claim isn't an administrator one: ", err)
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        // Setting up user role in context
        ctx := utils.SetUserRoleInContext(r.Context(), claims.Role)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Auth middleware
func CorsAuth(next http.Handler) http.Handler {
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
