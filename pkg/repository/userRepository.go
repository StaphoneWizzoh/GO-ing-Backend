package repository

import (
	"context"

	"time"

	"github.com/StaphoneWizzoh/Go_Auth/internal/database"
	"github.com/StaphoneWizzoh/Go_Auth/pkg/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	// create
	CreateUser(ctx context.Context, user model.UserRegister) (model.User, error)

	// update
	StoreRefreshToken(ctx context.Context, userId uuid.UUID, refreshToken string, expiresAt time.Time) (model.RefreshToken, error)
	UpdateUserLastLogin(ctx context.Context, userId uuid.UUID) error
	UpdateUser(ctx context.Context, user model.User) (model.User, error)
	UpdateUserProfilePicture(ctx context.Context, user model.User) (model.User, error)
	UpdateUserPassword(ctx context.Context, userId uuid.UUID, newPassword string) error

	PromoteUserToAdmin(ctx context.Context, userId uuid.UUID) (model.User, error)
	PromoteUserToSuperAdmin(ctx context.Context, userId uuid.UUID) (model.User, error)
	DemoteSuperAdminToAdmin(ctx context.Context, userId uuid.UUID)(model.User, error)
	DemoteSuperAdminToUser(ctx context.Context, userId uuid.UUID)(model.User, error)
	DemoteAdminToUser(ctx context.Context, userId uuid.UUID)(model.User, error)
	SuspendUser(ctx context.Context, userId uuid.UUID) (model.User, error)
	RecoverUser(ctx context.Context, userId uuid.UUID) (model.User, error)

	// delete
	DeleteUser(ctx context.Context, userId uuid.UUID) error

	// get
	CountAllUsersByUsername(ctx context.Context, username string) (int64, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)

	GetAllUsers(ctx context.Context, limit, offset int32) ([]database.User, error)
	GetAdminUsers(ctx context.Context, limit, offset int32) ([]database.User, error)
	GetSuperAdminUsers(ctx context.Context, limit, offset int32) ([]database.User, error)
	GetActiveUsers(ctx context.Context, limit, offset int32) ([]database.User, error)
	GetInactiveUsers(ctx context.Context, limit, offset int32) ([]database.User, error)
	GetSuspendedUsers(ctx context.Context, limit, offset int32) ([]database.User, error)
	GetDeletedUsers(ctx context.Context, limit, offset int32) ([]database.User, error)
	
}
