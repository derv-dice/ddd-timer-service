package service

import (
	"context"
)

func (s *Service) CheckUserHasServiceDates(ctx context.Context, userID int64) bool {
	u, err := s.repo.LoadUser(ctx, userID)
	if err != nil {
		return false
	}

	if !u.ServeFrom.IsZero() && !u.ServeTo.IsZero() {
		return true
	}

	return false
}
