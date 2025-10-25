package service

import (
	"context"
	"ddd-timer-service/internal/pkg/stats_counter"
	"time"
)

func (s *Service) GetUserStats(ctx context.Context, userID int64) (*stats_counter.Stats, error) {
	var err error

	user, err := s.repo.LoadUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats, err := stats_counter.NewStats(user, time.Now())
	if err != nil {
		return nil, err
	}

	return stats, nil
}
