package service

import (
	"context"
	"ddd-timer-service/models"
	"fmt"
)

func (s *Service) GetUser(_ context.Context, userID int64) (*models.User, error) {
	user := s.usersCache.Get(userID)
	if user == nil {
		return nil, fmt.Errorf("user not exist")
	}

	return user, nil
}
