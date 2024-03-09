package sqlc

import (
	"context"
	"database/sql"

	"log"
	"time"

	"github.com/StaphoneWizzoh/Go_Auth/internal/database"
	"github.com/StaphoneWizzoh/Go_Auth/pkg/model"
	"github.com/google/uuid"
)

type SQLUserRepository struct {
	DB *database.Queries
}

func NewSQLUserRepository(db *database.Queries) *SQLUserRepository {
	return &SQLUserRepository{
		DB: db,
	}
}

// CreateUser creates a new user
func (r *SQLUserRepository) CreateUser(ctx context.Context, user model.UserRegister) (model.User, error) {
	// insert user into database
	createdUser, err := r.DB.CreateUser(ctx, database.CreateUserParams{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		UserRole:       user.UserRole,
	})
	if err != nil {
		return model.User{}, err
	}

	// return created user
	return model.User{
		ID:             createdUser.ID,
		Username:       createdUser.Username,
		Email:          createdUser.Email,
		CreatedAt:      createdUser.CreatedAt,
		LastLogin:      createdUser.LastLogin,
		AccountStatus:  createdUser.AccountStatus,
		UserRole:       createdUser.UserRole,
		ProfilePicture: createdUser.ProfilePicture,
		TwoFactorAuth:  createdUser.TwoFactorAuth,
	}, nil
}

// CountAllUsersByUsername returns the number of users with the given username
func (r *SQLUserRepository) CountAllUsersByUsername(ctx context.Context, username string) (int64, error) {
	return r.DB.CountAllUsersByUsername(ctx, username)
}

// GetUserByEmail returns the user with the given email
func (r *SQLUserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	// get user from database
	user, err := r.DB.FindUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}

	// return user
	return model.User{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    user.PhoneNumber,
		DateOfBirth:    user.DateOfBirth,
		Gender:         user.Gender,
		CreatedAt:      user.CreatedAt,
		LastLogin:      user.LastLogin,
		AccountStatus:  user.AccountStatus,
		UserRole:       user.UserRole,
		ProfilePicture: user.ProfilePicture,
		TwoFactorAuth:  user.TwoFactorAuth,
	}, nil
}

// StoreRefreshToken stores the refresh token in the database
func (r *SQLUserRepository) StoreRefreshToken(ctx context.Context, userId uuid.UUID, refreshToken string, expiresAt time.Time) (model.RefreshToken, error) {
	log.Printf("Storing refresh token for user with id %s", userId.String())
	// insert refresh token into database
	createdRefreshToken, err := r.DB.StoreRefreshToken(ctx, database.StoreRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    userId,
		Token:     refreshToken,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	})
	if err != nil {
		log.Printf("Error storing refresh token for user with id %s: %s", userId.String(), err.Error())
		return model.RefreshToken{}, err
	}

	// return refresh token
	log.Printf("Successfully stored refresh token for user with id %s", userId.String())
	return model.RefreshToken{
		ID:        createdRefreshToken.ID,
		UserID:    createdRefreshToken.UserID,
		Token:     createdRefreshToken.Token,
		CreatedAt: createdRefreshToken.CreatedAt,
		ExpiresAt: createdRefreshToken.ExpiresAt,
		RevokedAt: createdRefreshToken.RevokedAt,
	}, nil
}

// UpdateUserLastLogin updates the last login of the user
func (r *SQLUserRepository) UpdateUserLastLogin(ctx context.Context, userId uuid.UUID) error {
	log.Printf("Updating last login of user with id %s", userId.String())

	err := r.DB.UpdateUserLastLogin(ctx, database.UpdateUserLastLoginParams{
		ID:        userId,
		LastLogin: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		log.Printf("Error updating last login of user with id %s: %s", userId.String(), err.Error())
	}
	return err
}

// GetUserById returns the user with the given id
func (r *SQLUserRepository) GetUserById(ctx context.Context, userId uuid.UUID) (model.User, error) {
	// get user from database
	user, err := r.DB.FindUserByID(ctx, userId)
	if err != nil {
		return model.User{}, err
	}

	// return user
	return model.User{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    user.PhoneNumber,
		DateOfBirth:    user.DateOfBirth,
		Gender:         user.Gender,
		CreatedAt:      user.CreatedAt,
		LastLogin:      user.LastLogin,
		AccountStatus:  user.AccountStatus,
		UserRole:       user.UserRole,
		ProfilePicture: user.ProfilePicture,
		TwoFactorAuth:  user.TwoFactorAuth,
	}, nil
}

// UpdateUser updates a user
func (r *SQLUserRepository) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	log.Printf("Updating user with id %s", user.ID.String())

	// update user
	updatedUser, err := r.DB.UpdateUser(ctx, database.UpdateUserParams{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		DateOfBirth: user.DateOfBirth,
		Gender:      user.Gender,
	})
	if err != nil {
		log.Printf("Error updating user with id %s: %s", user.ID.String(), err.Error())
		return model.User{}, err
	}

	// return updated user
	return model.User{
		ID:             updatedUser.ID,
		Username:       updatedUser.Username,
		Email:          updatedUser.Email,
		FirstName:      updatedUser.FirstName,
		LastName:       updatedUser.LastName,
		PhoneNumber:    updatedUser.PhoneNumber,
		DateOfBirth:    updatedUser.DateOfBirth,
		Gender:         updatedUser.Gender,
		CreatedAt:      updatedUser.CreatedAt,
		LastLogin:      updatedUser.LastLogin,
		AccountStatus:  updatedUser.AccountStatus,
		UserRole:       updatedUser.UserRole,
		ProfilePicture: updatedUser.ProfilePicture,
		TwoFactorAuth:  updatedUser.TwoFactorAuth,
	}, nil
}

// UpdateUserProfilePicture updates the profile picture of a user
func (r *SQLUserRepository) UpdateUserProfilePicture(ctx context.Context, user model.User) (model.User, error) {
	log.Printf("Updating profile picture of user with id %s", user.ID.String())

	// update user
	updatedUser, err := r.DB.UpdateUserProfilePicture(ctx, database.UpdateUserProfilePictureParams{
		ID:             user.ID,
		ProfilePicture: user.ProfilePicture,
	})
	if err != nil {
		log.Printf("Error updating profile picture of user with id %s: %s", user.ID.String(), err.Error())
		return model.User{}, err
	}

	// return updated user
	return model.User{
		ID:             updatedUser.ID,
		Username:       updatedUser.Username,
		Email:          updatedUser.Email,
		ProfilePicture: updatedUser.ProfilePicture,
	}, nil
}

// UpdateUserPassword updates the password of a user
func (r *SQLUserRepository) UpdateUserPassword(ctx context.Context, userId uuid.UUID, newPassword string) error {
	log.Printf("Updating password of user with id %s", userId.String())

	// update user
	_, err := r.DB.UpdateUserPassword(ctx, database.UpdateUserPasswordParams{
		ID:             userId,
		HashedPassword: newPassword,
	})
	if err != nil {
		log.Printf("Error updating password of user with id %s: %s", userId.String(), err.Error())
	}
	return err
}

// PromoteUserToAdmin promotes a regular user to an admin
func (r *SQLUserRepository) PromoteUserToAdmin(ctx context.Context, userId uuid.UUID) (model.User, error){
	log.Printf("Updating the role of user with id %s to admin status", userId.String())

	// update user role
	user, err := r.DB.PromoteUserToAdmin(ctx, userId)
	if err != nil {
		log.Printf("Error updating the role of user with id %s to admin status: %s", userId.String(), err.Error())
		return model.User{}, err
	}
	// Return updated user
	return model.User{
		Username:       user.Username,
		Email:          user.Email,
		LastLogin:      user.LastLogin,
		AccountStatus:  user.AccountStatus,
		UserRole:       user.UserRole,
	}, nil
}

// PromoteUserToSuperAdmin promotes a regular user to a super admin
func (r *SQLUserRepository) PromoteUserToSuperAdmin(ctx context.Context, userId uuid.UUID) (model.User, error){
	log.Printf("Updating the role of user with id %s to super admin status", userId.String())

	// update user
	user, err := r.DB.PromoteUserToSuperAdmin(ctx, userId)
	if err != nil {
		log.Printf("Error updating the role of user with id %s to super admin status: %s", userId.String(), err.Error())
		return model.User{}, err
	}
	return model.User{
		Username:       user.Username,
		Email:          user.Email,
		LastLogin:      user.LastLogin,
		AccountStatus:  user.AccountStatus,
		UserRole:       user.UserRole,
	}, nil
}

// DemoteSuperAdminToAdmin demotes a superadmin to an admin
func (r *SQLUserRepository) DemoteSuperAdminToAdmin(ctx context.Context, userId uuid.UUID)(model.User, error){
	log.Printf("Updating the role of super admin with id %s to admin status", userId.String())

	// update user
	admin, err := r.DB.DemoteSuperAdminToAdmin(ctx, userId)
	if err != nil{
		log.Printf("Error updating the role of super admin with id %s to admin status: %s", userId.String(), err.Error())
		return model.User{}, err
	}
	return model.User{
		Username:       admin.Username,
		Email:          admin.Email,
		LastLogin:      admin.LastLogin,
		AccountStatus:  admin.AccountStatus,
		UserRole:       admin.UserRole,
	}, nil
}

// DemoteSuperAdminToUser demotes a superadmin to a regular user
func (r *SQLUserRepository) DemoteSuperAdminToUser(ctx context.Context, userId uuid.UUID)(model.User, error){
	log.Printf("Updating the role of super admin with id %s to user status", userId.String())

	// update user
	user, err := r.DB.DemoteSuperAdminToUser(ctx, userId)
	if err != nil{
		log.Printf("Error updating the role of super admin with id %s to user status: %s", userId.String(), err.Error())
		return model.User{}, err
	}
	return model.User{
		Username:       user.Username,
		Email:          user.Email,
		LastLogin:      user.LastLogin,
		AccountStatus:  user.AccountStatus,
		UserRole:       user.UserRole,
	}, nil
}

// DemoteAdminToUser demotes an admin to a regular user
func (r *SQLUserRepository) DemoteAdminToUser(ctx context.Context, userId uuid.UUID)(model.User, error){
	log.Printf("Updating the role of super admin with id %s to user status", userId.String())

	// update user
	user, err := r.DB.DemoteAdminToUser(ctx, userId)
	if err != nil{
		log.Printf("Error updating the role of admin with id %s to user status: %s", userId.String(), err.Error())
		return model.User{}, err
	}
	return model.User{
		Username:       user.Username,
		Email:          user.Email,
		LastLogin:      user.LastLogin,
		AccountStatus:  user.AccountStatus,
		UserRole:       user.UserRole,
	}, nil
}

// SuspendUser suspendds an active user account
func (r *SQLUserRepository) SuspendUser(ctx context.Context, userId uuid.UUID) (model.User, error){
	log.Printf("Suspending user account with id %s :", userId.String())

	// suspending user
	user, err := r.DB.SuspendUser(ctx, userId)
	if err != nil{
		log.Printf("Error suspending user account with id %s: %s", userId.String(), err.Error())
		return model.User{}, err
	}

	return model.User{
		Username:       user.Username,
		Email:          user.Email,
		LastLogin:      user.LastLogin,
		AccountStatus:  user.AccountStatus,
		UserRole:       user.UserRole,
	}, nil
}

// RecoverUser returns a a suspended user to an active user
func (r *SQLUserRepository) RecoverUser(ctx context.Context, userId uuid.UUID) (model.User, error){
	log.Printf("Recovering user with id %s to an active one:", userId.String())
	// update user
	user, err := r.DB.RecoverUser(ctx, userId)
	if err != nil {
		log.Printf("Error recovering user with id %s to an active one: %s", userId.String(), err.Error())
		return model.User{}, err
	}
	return model.User{
		Username:       user.Username,
		Email:          user.Email,
		LastLogin:      user.LastLogin,
		AccountStatus:  user.AccountStatus,
		UserRole:       user.UserRole,
	}, nil
}

func (r *SQLUserRepository) DeleteUser(ctx context.Context, userId uuid.UUID) error{
	log.Printf("Changing the account status of user with id %s to inactive:", userId.String())

	// deleting the user
	err := r.DB.DeleteUser(ctx, userId)
	if err != nil{
		log.Printf("Error deleting user with id %s : %s", userId.String(), err.Error())
		return err
	}
	return nil
}


// GetAllUsers returns a list of all accounts ever registered
func (r *SQLUserRepository) GetAllUsers(ctx context.Context, limit, offset int32) ([]database.User, error){
	// Fetching all users from the database
	users, err := r.DB.GetAllUsers(ctx, database.GetAllUsersParams{
		Limit: limit,
		Offset: offset,
	})
	if err != nil{
		log.Printf("Error fetching all users from database: %s", err)
		return []database.User{}, err
	}
	
	return users, nil
}

// GetAdminUsers returns a list of admin users
func (r *SQLUserRepository) GetAdminUsers(ctx context.Context, limit, offset int32) ([]database.User, error){
	// Festching all administrators from the database
	admins, err := r.DB.GetAdminUsers(ctx, database.GetAdminUsersParams{
		Limit: limit,
		Offset: offset,
	})

	if err != nil{
		log.Printf("Error fetching administrators from database %s", err)
		return []database.User{}, err
	}

	return admins, nil
}

// GetSuperAdminUsers returns a list of super admin users
func (r *SQLUserRepository) GetSuperAdminUsers(ctx context.Context, limit, offset int32) ([]database.User, error){
	// Festching all administrators from the database
	superAdmins, err := r.DB.GetSuperAdminUsers(ctx, database.GetSuperAdminUsersParams{
		Limit: limit,
		Offset: offset,
	})

	if err != nil{
		log.Printf("Error fetching super administrators from database %s", err)
		return []database.User{}, err
	}

	return superAdmins, nil
}

// GetActiveUsers returns a list of active users
func (r *SQLUserRepository) GetActiveUsers(ctx context.Context, limit, offset int32) ([]database.User, error){
	// Fetching active users from the database
	activeUsers, err := r.DB.GetActiveUsers(ctx, database.GetActiveUsersParams{
		Offset: offset,
		Limit: limit,
	})

	if err != nil{
		log.Printf("Error fetching active users from the database %s", err)
		return []database.User{}, err
	}

	return activeUsers, nil
}

// GetInactiveUsers returns a list of inactive users
func (r *SQLUserRepository) GetInactiveUsers(ctx context.Context, limit, offset int32) ([]database.User, error){
	// Fetching inactive users from the database
	inactiveUsers, err := r.DB.GetInactiveUsers(ctx, database.GetInactiveUsersParams{
		Limit: limit,
		Offset: offset,
	})

	if err != nil{
		log.Printf("Error fetching inactive users from the database: %s", err)
		return []database.User{}, err
	}

	return inactiveUsers, nil
}

// GetSuspendedUsers returns a list of suspended users
func (r *SQLUserRepository) GetSuspendedUsers(ctx context.Context, limit, offset int32) ([]database.User, error){
	// Fetching suspended users from the database
	suspendedUsers, err := r.DB.GetSuspendedUsers(ctx, database.GetSuspendedUsersParams{
		Limit: limit,
		Offset: offset,
	})

	if err != nil{
		log.Printf("Error fetching suspended users from the database: %s", err)
		return []database.User{}, err
	}

	return suspendedUsers, nil
}

// GetDeletedUsers returns a list of users with account status disabled
func (r *SQLUserRepository) GetDeletedUsers(ctx context.Context, limit, offset int32) ([]database.User, error){
	// Fetching disabled users from the database
	disabledUsers, err := r.DB.GetDeletedUsers(ctx, database.GetDeletedUsersParams{
		Limit: limit,
		Offset: offset,
	})

	if err != nil{
		log.Printf("Error fetching disabled users from the database: %s", err)
		return []database.User{}, err
	}

	return disabledUsers, nil
}