package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"html/template"
	"net/http"

	"github.com/StaphoneWizzoh/Go_Auth/pkg/usecases"
)

type UserHandler struct {
	userService *usecases.UserService
}

func NewUserHandler(userService *usecases.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// User accessible handlers

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// params
	var params struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		UserRole  string `json:"user_role"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// create user
	user, err := h.userService.CreateUser(r.Context(), params.Email,
		params.Password, params.FirstName, params.LastName, params.UserRole)
	if err != nil {
		// User email already exist in the database
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint \"users_email_key\""){
			RespondWithError(w, http.StatusBadRequest, "User email already exists")
			return
		}

		// User username(firstname+lastname) already exists in the database
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint \"users_username_key\""){
			RespondWithError(w, http.StatusBadRequest, "User firstname and lastname combination already exists")
			return
		}

		// Other errors
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	// respond with user
	RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// params
	var params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// login user
	user, err := h.userService.LoginUser(r.Context(), params.Email, params.Password)
	if err != nil {
		// Check for specific error cases and return corresponding status codes
		// For non-existing user
		if strings.Contains(err.Error(), "sql: no rows in result set"){
			RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		// for incorrect passwords
		if strings.Contains(err.Error(), "crypto/bcrypt: hashedPassword is not the hash of the given password"){
			RespondWithError(w, http.StatusUnauthorized, "Incorrect password")
			return
		}

		// Other errors
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to login user: %v", err))
		return
	}

	// respond with user
	RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// params
	var params struct {
		RefreshToken string `json:"refresh_token"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// refresh token
	user, err := h.userService.RefreshToken(r.Context(), params.RefreshToken)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to refresh token: %v", err))
		return
	}

	// respond with user
	RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// params
	var params struct {
		Email       string `json:"email"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		DateOfBirth string `json:"date_of_birth"`
		Gender      string `json:"gender"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// update user
	user, err := h.userService.UpdateUser(r.Context(), params.Email, params.FirstName, params.LastName,
		params.PhoneNumber, params.Gender, params.DateOfBirth)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update user: %v", err))
		return
	}

	// respond with user
	RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	// params
	var params struct {
		ProfilePicture string `json:"profile_picture"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// update user
	user, err := h.userService.UpdateProfilePicture(r.Context(), params.ProfilePicture)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update profile picture: %v", err))
		return
	}

	// respond with user
	RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	// params
	var params struct {
		Email string `json:"email"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// send reset password email
	if err := h.userService.SendResetPasswordEmail(r.Context(), params.Email); err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to send reset password email: %v", err))
		return
	}

	// respond with success message
	RespondWithSuccess(w, http.StatusOK, "Reset password email sent successfully")
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Serve the reset password form with the token embedded as a hidden input
		token := r.URL.Query().Get("token")
		if token == "" {
			RespondWithError(w, http.StatusBadRequest, "Token is required")
			return
		}

		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("pkg/templates/reset-password.html"))
		data := map[string]interface{}{
			"Token": token,
		}
		if err := tmpl.Execute(w, data); err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to open reset password page: %v", err))
			return
		}
	} else if r.Method == http.MethodPost {
		// Handle form submission
		err := r.ParseForm()
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid form data")
			return
		}
		token := r.FormValue("token")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// Validate the inputs
		if token == "" {
			RespondWithError(w, http.StatusBadRequest, "Token is required")
			return
		}
		if password != confirmPassword {
			RespondWithError(w, http.StatusBadRequest, "Passwords do not match")
			return
		}

		// Reset password logic...
		if err := h.userService.ResetPassword(r.Context(), token, password); err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to reset password: %v", err))
			return
		}

		// Respond with success message
		RespondWithSuccess(w, http.StatusOK, "Password reset successfully")
	}
}

// Admin accessible handlers

func (h *UserHandler) PromoteUserToAdmin(w http.ResponseWriter, r *http.Request){
		// params
		var params struct {
			Email     string `json:"email"`
		}
	
		// decode request body
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
			return
		}
	
		// get user id
		user, err := h.userService.GetUserByEmail(r.Context(), params.Email)

		// check if the if the user is already an admin
		if user.UserRole == "admin"{
			RespondWithError(w, http.StatusBadRequest,"User already an admin")
			return
		}

		// Error handling
		if err != nil {
			// For non-existing user
			if strings.Contains(err.Error(), "sql: no rows in result set"){
				RespondWithError(w, http.StatusNotFound, "User not found")
				return
			}

			// Other errors
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to promote user: %v", err))
			return
		}

		// promote user
		updatedUser, err := h.userService.PromoteUserToAdmin(r.Context(), user.ID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to promote user to admin status : %v", err))
			return
		}
	
		// respond with promoted user
		RespondWithJSON(w, http.StatusOK, updatedUser)
}

func (h *UserHandler) PromoteUserToSuperAdmin(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Email     string `json:"email"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// get user id
	user, err := h.userService.GetUserByEmail(r.Context(), params.Email)

	// check if the if the user is already an admin
	if user.UserRole == "superadmin"{
		RespondWithError(w, http.StatusBadRequest,"User already a super admin")
		return
	}

	if err != nil {
		// For non-existing user
		if strings.Contains(err.Error(), "sql: no rows in result set"){
			RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		// Other errors
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to promote user: %v", err))
		return
	}

	// promote user
	updatedUser, err := h.userService.PromoteUserToSuperAdmin(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to promote user to super admin status : %v", err))
		return
	}

	// respond with promoted user
	RespondWithJSON(w, http.StatusOK, updatedUser)
}

func (h *UserHandler) DemoteSuperAdminToAdmin(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Email     string `json:"email"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// get superadmin id
	superAdmin, err := h.userService.GetUserByEmail(r.Context(), params.Email)

	
	if err != nil {
		// For non-existing user
		if strings.Contains(err.Error(), "sql: no rows in result set"){
			RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		// Other errors
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to promote user: %v", err))
		return
	}

	// Check if user is a superadmin
	if superAdmin.UserRole != "superadmin"{
		RespondWithError(w, http.StatusBadRequest,"User is not a super administrator")
		return
	}

	// demote superadmin to admin
	admin, err := h.userService.DemoteSuperAdminToAdmin(r.Context(), superAdmin.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to demote superadmin to admin status : %v", err))
		return
	}

	// respond with promoted user
	RespondWithJSON(w, http.StatusOK, admin)
}

func (h *UserHandler) DemoteSuperAdminToUser(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Email     string `json:"email"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// get superadmin id
	superAdmin, err := h.userService.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		// For non-existing user
		if strings.Contains(err.Error(), "sql: no rows in result set"){
			RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		// Other errors
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to promote user: %v", err))
		return
	}

	// Check if user is a superadmin
	if superAdmin.UserRole != "superadmin"{
		RespondWithError(w, http.StatusBadRequest,"User is not a super administrator")
		return
	}

	// demote superadmin to user
	user, err := h.userService.DemoteSuperAdminToUser(r.Context(), superAdmin.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to demote superadmin to user status : %v", err))
		return
	}

	// respond with promoted user
	RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DemoteAdminToUser(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Email     string `json:"email"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// get admin id
	admin, err := h.userService.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		// For non-existing user
		if strings.Contains(err.Error(), "sql: no rows in result set"){
			RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		// Other errors
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to promote user: %v", err))
		return
	}

	// Check if user is a superadmin
	if admin.UserRole != "admin"{
		RespondWithError(w, http.StatusBadRequest,"User is not an administrator")
		return
	}

	// demote admin to user
	user, err := h.userService.DemoteAdminToUser(r.Context(), admin.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to demote admin to user status : %v", err))
		return
	}

	// respond with promoted user
	RespondWithJSON(w, http.StatusOK, user)
}


func (h *UserHandler) SuspendUser(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Email     string `json:"email"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// get user id
	user, err := h.userService.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch user details: %v", err))
		return
	}

	// Check if user is suspended
	if user.AccountStatus == "suspended"{
		RespondWithError(w, http.StatusBadRequest,"User is already suspended")
		return
	}

	// suspend user
	suspendedUser, err := h.userService.SuspendUser(r.Context(), user.ID)

	if err != nil {
		// For non-existing user
		if strings.Contains(err.Error(), "sql: no rows in result set"){
			RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		// Other errors
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to suspend user: %v", err))
		return
	}

	

	// respond with suspended user
	RespondWithJSON(w, http.StatusOK, suspendedUser)
}

func (h *UserHandler) RecoverUser(w http.ResponseWriter, r *http.Request){
		// params
		var params struct {
			Email     string `json:"email"`
		}
	
		// decode request body
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
			return
		}
	
		// get user id
		user, err := h.userService.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch user details: %v", err))
			return
		}
		
		// Check if user is recovered
		if user.AccountStatus == "active"{
			RespondWithError(w, http.StatusBadRequest,"User account is already recovered")
			return
		}
	
		// recover user
		recoveredUser, err := h.userService.RecoverUser(r.Context(), user.ID)
		if err != nil {
			// For non-existing user
			if strings.Contains(err.Error(), "sql: no rows in result set"){
				RespondWithError(w, http.StatusNotFound, "User not found")
				return
			}
	
			// Other errors
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to recover user: %v", err))
			return
		}

	
		// respond with recovered user
		RespondWithJSON(w, http.StatusOK, recoveredUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Email     string `json:"email"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	// get user id
	user, err := h.userService.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch user details: %v", err))
		return
	}

	// Check if user account is deadctivated
	if user.AccountStatus == "deleted"{
		RespondWithError(w, http.StatusBadRequest,"User account is already deleted.")
		return
	}

	// deactivate user account
	err = h.userService.DeleteUser(r.Context(), user.ID)
	if err != nil {
		// For non-existing user
		if strings.Contains(err.Error(), "sql: no rows in result set"){
			RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		// Other errors
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete user: %v", err))
		return
	}

	

	// respond with status
	RespondWithSuccess(w, http.StatusOK, "Successfully deactivated the user's account")
}


func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request){

	// params
	var params struct {
		Limit	    int32		`json:"limit"`
		Offset		int32		`json:"offset"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	users, err := h.userService.GetAllUsers(r.Context(), params.Limit, params.Offset)
	if err != nil{
		RespondWithError(w, http.StatusInternalServerError, "Error fetching all users from the database.")
		return
	}
	if len(users)==0 {
		RespondWithSuccess(w, http.StatusNotFound, "There are no users in the database")
		return
	}

	RespondWithJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetAdminUsers(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Limit	    int32		`json:"limit"`
		Offset		int32		`json:"offset"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	admins, err := h.userService.GetAdminUsers(r.Context(), params.Limit, params.Offset)
	if err != nil{
		RespondWithError(w, http.StatusInternalServerError, "Error fetching all administrators from the database.")
		return
	}
	if len(admins)==0 {
		RespondWithSuccess(w, http.StatusNotFound, "There are no administrators in the database")
		return
	}

	RespondWithJSON(w, http.StatusOK, admins)
}

func (h *UserHandler) GetSuperAdminUsers(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Limit	    int32		`json:"limit"`
		Offset		int32		`json:"offset"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	superAdmins, err := h.userService.GetSuperAdminUsers(r.Context(), params.Limit, params.Offset)
	if err != nil{
		RespondWithError(w, http.StatusInternalServerError, "Error fetching all super administrators from the database.")
		return
	}
	if len(superAdmins)==0 {
		RespondWithSuccess(w, http.StatusNotFound, "There are no super administrators in the database")
		return
	}

	RespondWithJSON(w, http.StatusOK, superAdmins)
}

func (h *UserHandler) GetActiveUsers(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Limit	    int32		`json:"limit"`
		Offset		int32		`json:"offset"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	activeUsers, err := h.userService.GetActiveUsers(r.Context(), params.Limit, params.Offset)
	if err != nil{
		RespondWithError(w, http.StatusInternalServerError, "Error fetching all active users from the database.")
		return
	}
	if len(activeUsers)==0 {
		RespondWithSuccess(w, http.StatusNotFound, "There are no active users in the database")
		return
	}

	RespondWithJSON(w, http.StatusOK, activeUsers)
}

func (h *UserHandler) GetInactiveUsers(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Limit	    int32		`json:"limit"`
		Offset		int32		`json:"offset"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	inactiveUsers, err := h.userService.GetInactiveUsers(r.Context(), params.Limit, params.Offset)
	if err != nil{
		RespondWithError(w, http.StatusInternalServerError, "Error fetching inactive users from the database.")
		return
	}
	if len(inactiveUsers)==0 {
		log.Println("Inactive users:", inactiveUsers)
		RespondWithSuccess(w, http.StatusNotFound, "There are no inactive users in the database")
		return
	}
	
	RespondWithJSON(w, http.StatusOK, inactiveUsers)
}

func (h *UserHandler) GetSuspendedUsers(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Limit	    int32		`json:"limit"`
		Offset		int32		`json:"offset"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	suspendedUsers, err := h.userService.GetSuspendedUsers(r.Context(), params.Limit, params.Offset)
	if err != nil{
		RespondWithError(w, http.StatusInternalServerError, "Error fetching suspended users from the database.")
		return
	}
	if len(suspendedUsers)==0 {
		RespondWithSuccess(w, http.StatusNotFound, "There are no suspended users in the database")
		return
	}

	RespondWithJSON(w, http.StatusOK, suspendedUsers)
}

func (h *UserHandler) GetDeletedUsers(w http.ResponseWriter, r *http.Request){
	// params
	var params struct {
		Limit	    int32		`json:"limit"`
		Offset		int32		`json:"offset"`
	}

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	disabledUsers, err := h.userService.GetDeletedUsers(r.Context(), params.Limit, params.Offset)
	if err != nil{
		RespondWithError(w, http.StatusInternalServerError, "Error fetching inactive users from the database.")
		return
	}
	if len(disabledUsers)==0 {
		RespondWithSuccess(w, http.StatusNotFound, "There are no disabled users in the database")
		return
	}

	RespondWithJSON(w, http.StatusOK, disabledUsers)
}