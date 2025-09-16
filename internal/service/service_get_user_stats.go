package service

import (
	"context"
	"ddd-timer-service/internal/pkg/stats_counter"
	"fmt"
	"time"
)

func (s *Service) GetUserStats(_ context.Context, userID int64) (*stats_counter.Stats, error) {
	var err error

	user := s.usersCache.Get(userID)
	if user == nil {
		return nil, fmt.Errorf("user not exist")
	}

	stats, err := stats_counter.NewStats(user, time.Now())
	if err != nil {
		return nil, err
	}

	return stats, nil
}
