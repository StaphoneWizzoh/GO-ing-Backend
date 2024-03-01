package main

import (
	"log"
	"net/http"

	"github.com/StaphoneWizzoh/Go_Auth/internal/database"
	"github.com/StaphoneWizzoh/Go_Auth/pkg/config"
	"github.com/StaphoneWizzoh/Go_Auth/pkg/handlers"
	"github.com/StaphoneWizzoh/Go_Auth/pkg/middleware"
	"github.com/StaphoneWizzoh/Go_Auth/pkg/repository/sqlc"
	"github.com/StaphoneWizzoh/Go_Auth/pkg/usecases"
	"github.com/gorilla/mux"
)

func main(){
	cfg := config.LoadConfig()

	// Database connection initialization
	conn, err := config.NewDatabaseConnection(cfg.DbURL)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
	}

	db := database.New(conn)

	// Repository initializations
	userRepo := sqlc.NewSQLUserRepository(db)

	// Services initializations
	userService := usecases.NewUserService(userRepo)

	// Handlers initializations
	userHandler := handlers.NewUserHandler(userService)
	

	// Setting up routes
	r := mux.NewRouter()
	r.Use(middleware.CORS)
	getUserRouter(r, userHandler)	

	// Starting the server
	log.Printf("Server listening on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}