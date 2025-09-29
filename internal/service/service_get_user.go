package service

import (
	"context"
	"ddd-timer-service/models"
)

func (s *Service) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.repo.LoadUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
