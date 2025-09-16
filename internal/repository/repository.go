package repository

import (
	"context"
	"ddd-timer-service/models"
)

type Repository interface {
	SaveUser(ctx context.Context, user *models.User) error
	LoadUser(ctx context.Context, userID int64) (*models.User, error)
	DeleteUser(ctx context.Context, userID int64) error

	LoadAllUsers(ctx context.Context) ([]*models.User, error)
}
