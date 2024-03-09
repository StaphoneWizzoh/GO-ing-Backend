package main

import (
	"net/http"

	"github.com/StaphoneWizzoh/Go_Auth/pkg/handlers"
	"github.com/StaphoneWizzoh/Go_Auth/pkg/middleware"
	"github.com/gorilla/mux"
)

func getUserRouter(r *mux.Router, userHandler *handlers.UserHandler) {
	resetPasswordRouter := r.PathPrefix("/reset-password").Subrouter()
	resetPasswordRouter.HandleFunc("", userHandler.ResetPassword).Methods(http.MethodGet)
	resetPasswordRouter.HandleFunc("", userHandler.ResetPassword).Methods(http.MethodPost)

	userRouter := r.PathPrefix("/api/users").Subrouter()
	userRouter.HandleFunc("/register", userHandler.RegisterUser).Methods(http.MethodPost)
	userRouter.HandleFunc("/login", userHandler.LoginUser).Methods(http.MethodPost)
	userRouter.HandleFunc("/refresh", userHandler.RefreshToken).Methods(http.MethodPost)

	// Protected Routes

	// Authenticated user routes
	protectedUserRouter := userRouter.PathPrefix("").Subrouter()
	protectedUserRouter.Use(middleware.Auth)
	protectedUserRouter.HandleFunc("/update", userHandler.UpdateUser).Methods(http.MethodPut)
	protectedUserRouter.HandleFunc("/update-profile-picture", userHandler.UpdateProfilePicture).Methods(http.MethodPut)
	protectedUserRouter.HandleFunc("/reset-password", userHandler.RequestPasswordReset).Methods(http.MethodPut)

	// Authenticated admin routes
	protectedAdminRouter := r.PathPrefix("/api/admin").Subrouter()
	protectedAdminRouter.Use(middleware.Auth)
	protectedAdminRouter.Use(middleware.Admin)
	protectedAdminRouter.HandleFunc("/promote-admin", userHandler.PromoteUserToAdmin).Methods(http.MethodPut)
	protectedAdminRouter.HandleFunc("/promote-super-admin", userHandler.PromoteUserToSuperAdmin).Methods(http.MethodPut)
	protectedAdminRouter.HandleFunc("/demote-super-admin-to-admin", userHandler.DemoteSuperAdminToAdmin).Methods(http.MethodPut)
	protectedAdminRouter.HandleFunc("/demote-super-admin-to-user", userHandler.DemoteSuperAdminToUser).Methods(http.MethodPut)
	protectedAdminRouter.HandleFunc("/demote-admin-to-user", userHandler.DemoteAdminToUser).Methods(http.MethodPut)
	protectedAdminRouter.HandleFunc("/suspend-user", userHandler.SuspendUser).Methods(http.MethodPut)
	protectedAdminRouter.HandleFunc("/recover-user", userHandler.RecoverUser).Methods(http.MethodPut)
	protectedAdminRouter.HandleFunc("/delete-user", userHandler.DeleteUser).Methods(http.MethodDelete)
	protectedAdminRouter.HandleFunc("/all-users", userHandler.GetAllUsers).Methods(http.MethodPost)
	protectedAdminRouter.HandleFunc("/active-users", userHandler.GetActiveUsers).Methods(http.MethodPost)
	protectedAdminRouter.HandleFunc("/admin-users", userHandler.GetAdminUsers).Methods(http.MethodPost)
	protectedAdminRouter.HandleFunc("/super-admin-users", userHandler.GetSuperAdminUsers).Methods(http.MethodPost)
	protectedAdminRouter.HandleFunc("/deleted-users", userHandler.GetDeletedUsers).Methods(http.MethodPost)
	protectedAdminRouter.HandleFunc("/inactive-users", userHandler.GetInactiveUsers).Methods(http.MethodPost)
	protectedAdminRouter.HandleFunc("/suspended-users", userHandler.GetSuspendedUsers).Methods(http.MethodPost)
}