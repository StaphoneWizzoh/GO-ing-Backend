package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID      `json:"id"`
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	HashedPassword string         `json:"hashed_password"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	PhoneNumber    sql.NullString `json:"phone_number"`
	DateOfBirth    sql.NullTime   `json:"date_of_birth"`
	Gender         sql.NullString `json:"gender"`
	CreatedAt      time.Time      `json:"created_at"`
	LastLogin      sql.NullTime   `json:"last_login"`
	AccountStatus  string         `json:"account_status"`
	UserRole       string         `json:"user_role"`
	ProfilePicture sql.NullString `json:"profile_picture"`
	TwoFactorAuth  bool           `json:"two_factor_auth"`
}

type UserRegister struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	UserRole       string    `json:"user_role"`
}
