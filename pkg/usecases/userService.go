package usecases

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/StaphoneWizzoh/Go_Auth/internal/database"
	"github.com/StaphoneWizzoh/Go_Auth/pkg/model"
	"github.com/StaphoneWizzoh/Go_Auth/pkg/repository"
	"github.com/StaphoneWizzoh/Go_Auth/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(
	ctx context.Context,
	email string,
	password string,
	firstName string,
	lastName string,
	userRole string,
) (model.User, error) {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	// generate username
	username, err := generateUsername(ctx, s.userRepo, firstName, lastName)
	if err != nil {
		return model.User{}, err
	}

	// create user
	newUser := model.UserRegister{
		ID:             uuid.New(),
		Username:       username,
		Email:          email,
		HashedPassword: string(hashedPassword),
		FirstName:      strings.ToLower(firstName),
		LastName:       strings.ToLower(lastName),
		UserRole:       userRole,
	}

	return s.userRepo.CreateUser(ctx, newUser)
}

// sanitizeUsername sanitizes the given username
func sanitizeUsername(username string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9_.-]+")
	return reg.ReplaceAllString(username, "")
}

// generateUsername generates a username from the given first and last name
func generateUsername(ctx context.Context, userRepo repository.UserRepository, firstName string, lastName string) (string, error) {
	// generate username from first and last name
	baseUsername := strings.ToLower(sanitizeUsername(firstName + lastName))
	username := baseUsername
	const maxUsernameLength = 20

	// trim username to max length
	if len(baseUsername) > maxUsernameLength {
		baseUsername = baseUsername[:maxUsernameLength]
		username = baseUsername
	}

	for suffix := 1; ; suffix++ {
		// check if username is available
		count, err := userRepo.CountAllUsersByUsername(ctx, username)
		if err != nil {
			return "", err
		}
		if count == 0 {
			return username, nil
		}

		// append suffix to username
		suffixStr := strconv.Itoa(suffix)
		cutOffLength := maxUsernameLength - len(suffixStr)
		if cutOffLength < len(baseUsername) {
			username = baseUsername[:cutOffLength]
		}

		username += suffixStr
	}
}

// LoginUser logs in a user
func (s *UserService) LoginUser(ctx context.Context, email string, password string) (model.LoginResponse, error) {
	// get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return model.LoginResponse{}, err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return model.LoginResponse{}, err
	}

	// generate access token and refresh token
	accessToken, refreshToken, expireTime, err := utils.GenerateTokens(user.ID, user.Username, user.Email, user.UserRole)
	if err != nil {
		return model.LoginResponse{}, err
	}

	// save refresh token
	_, err = s.userRepo.StoreRefreshToken(ctx, user.ID, refreshToken, expireTime)
	if err != nil {
		return model.LoginResponse{}, err
	}

	// update last login
	err = s.userRepo.UpdateUserLastLogin(ctx, user.ID)
	if err != nil {
		log.Printf("Failed to update last login for user with id %s: %v", user.ID.String(), err)
	}

	return model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken refreshes a user's access token
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (model.LoginResponse, error) {
	// generate access token
	newAccessToken, err := utils.RefreshToken(refreshToken)
	if err != nil {
		return model.LoginResponse{}, err
	}

	return model.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
	}, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(
	ctx context.Context,
	email string,
	firstName string,
	lastName string,
	phoneNumber string,
	gender string,
	dateOfBirth string,
) (model.User, error) {
	// get user
	userId := ctx.Value("userId").(uuid.UUID)
	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return model.User{}, err
	}

	// update user
	if email != "" {
		user.Email = email
	}

	if firstName != "" {
		user.FirstName = firstName
	}

	if lastName != "" {
		user.LastName = lastName
	}

	if phoneNumber != "" {
		phoneNumberValue := sql.NullString{String: phoneNumber, Valid: phoneNumber != ""}
		user.PhoneNumber = phoneNumberValue
	}

	if gender != "" {
		genderValue := sql.NullString{String: gender, Valid: gender != ""}
		user.Gender = genderValue
	}

	if dateOfBirth != "" {
		dateOfBirthDate, err := time.Parse("02-01-2006", dateOfBirth)
		if err != nil {
			dateOfBirthDate = time.Time{}
		}
		dateOfBirthValue := sql.NullTime{Time: dateOfBirthDate, Valid: dateOfBirthDate != time.Time{}}
		user.DateOfBirth = dateOfBirthValue
	}

	return s.userRepo.UpdateUser(ctx, user)
}

// UpdateProfilePicture updates a user's profile picture
func (s *UserService) UpdateProfilePicture(ctx context.Context, profilePicture string) (model.User, error) {
	// get user
	userId := ctx.Value("userId").(uuid.UUID)
	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return model.User{}, err
	}

	// update user
	profilePictureValue := sql.NullString{String: profilePicture, Valid: profilePicture != ""}
	user.ProfilePicture = profilePictureValue

	return s.userRepo.UpdateUserProfilePicture(ctx, user)
}

// SendResetPasswordEmail sends a reset password email to the user
func (s *UserService) SendResetPasswordEmail(ctx context.Context, email string) error {
	// get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	// check if valid user
	if user.ID == uuid.Nil {
		return errors.New("invalid user")
	}

	// send reset password email
	return utils.SendResetPasswordEmail(user.ID, user.Email)
}

// ResetPassword resets a user's password
func (s *UserService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	// verify reset password token
	userId, err := utils.VerifyResetPasswordToken(token)
	if err != nil {
		return err
	}

	// validate new password
	//err = utils.ValidatePassword(newPassword)
	//if err != nil {
	//	return err
	//}

	// hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// update user password
	return s.userRepo.UpdateUserPassword(ctx, userId, string(hashedPassword))
}

// VerifyResetPasswordToken verifies a reset password token
func (s *UserService) VerifyResetPasswordToken(ctx context.Context, token string) (uuid.UUID, error) {
	// verify reset password token
	passwordToken, err := utils.VerifyResetPasswordToken(token)
	return passwordToken, err
}

func (s *UserService) PromoteUserToAdmin(ctx context.Context, userId uuid.UUID)(model.User, error){
	return s.userRepo.PromoteUserToAdmin(ctx, userId)
}

func (s *UserService) PromoteUserToSuperAdmin(ctx context.Context, userId uuid.UUID)(model.User, error){
	return s.userRepo.PromoteUserToSuperAdmin(ctx, userId)
}

func (s *UserService) DemoteSuperAdminToAdmin(ctx context.Context, userId uuid.UUID)(model.User, error){
	return s.userRepo.DemoteSuperAdminToAdmin(ctx, userId)
}

func (s *UserService) DemoteSuperAdminToUser(ctx context.Context, userId uuid.UUID)(model.User, error){
	return s.userRepo.DemoteSuperAdminToUser(ctx, userId)
}

func (s *UserService) DemoteAdminToUser(ctx context.Context, userId uuid.UUID)(model.User, error){
	return s.userRepo.DemoteAdminToUser(ctx, userId)
}

func (s *UserService) SuspendUser(ctx context.Context, userId uuid.UUID)(model.User, error){
	return s.userRepo.SuspendUser(ctx, userId)
}

func (s *UserService) RecoverUser(ctx context.Context, userId uuid.UUID)(model.User, error){
	return s.userRepo.RecoverUser(ctx, userId)
}

func (s *UserService) DeleteUser(ctx context.Context, userId uuid.UUID) error{
	return s.userRepo.DeleteUser(ctx, userId)
}


func (s *UserService) GetUserByEmail(ctx context.Context, email string)(model.User, error){
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetAllUsers(ctx context.Context, limit, offset int32)([]database.User, error){
	return s.userRepo.GetAllUsers(ctx, limit, offset)
}

func (s *UserService) GetAdminUsers(ctx context.Context, limit, offset int32)([]database.User, error){
	return s.userRepo.GetAdminUsers(ctx, limit, offset)
}

func (s *UserService) GetSuperAdminUsers(ctx context.Context, limit, offset int32)([]database.User, error){
	return s.userRepo.GetSuperAdminUsers(ctx, limit, offset)
}

func (s *UserService) GetActiveUsers(ctx context.Context, limit, offset int32)([]database.User, error){
	return s.userRepo.GetActiveUsers(ctx, limit, offset)
}

func (s *UserService) GetInactiveUsers(ctx context.Context, limit, offset int32)([]database.User, error){
	return s.userRepo.GetInactiveUsers(ctx, limit, offset)
}

func (s *UserService) GetSuspendedUsers(ctx context.Context, limit, offset int32)([]database.User, error){
	return s.userRepo.GetSuspendedUsers(ctx, limit, offset)
}

func (s *UserService) GetDeletedUsers(ctx context.Context, limit, offset int32)([]database.User, error){
	return s.userRepo.GetDeletedUsers(ctx, limit, offset)
}


