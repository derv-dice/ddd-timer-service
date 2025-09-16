package service

import (
	"context"
)

func (s *Service) CheckUserHasServiceDates(_ context.Context, userID int64) bool {
	// Если нет в кэше, значит и в БД тоже не может быть.
	u := s.usersCache.Get(userID)
	if u == nil {
		return false
	}

	if !u.ServeFrom.IsZero() && !u.ServeTo.IsZero() {
		return true
	}

	return false
}
