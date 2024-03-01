package utils

import (
	"context"

	"github.com/google/uuid"
)

func SetUserIdInContext(ctx context.Context, userID uuid.UUID) context.Context{
	return context.WithValue(ctx, "userId", userID)
}

func SetUserRoleInContext(ctx context.Context,  userRole string) context.Context{
	return context.WithValue(ctx, "userRole", userRole)
}