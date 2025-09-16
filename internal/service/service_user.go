package service

import (
	"context"
	"ddd-timer-service/models"
	"fmt"
	"time"
)

func (s *Service) SaveUser(ctx context.Context, tgID int64, from, to time.Time) error {
	if tgID == 0 {
		return fmt.Errorf("tgID is empty")
	}

	u := &models.User{
		ID:        tgID,
		ServeFrom: from,
		ServeTo:   to,
	}

	if err := u.Validate(); err != nil {
		return err
	}

	if err := s.repo.SaveUser(ctx, u); err != nil {
		return err
	}

	s.usersCache.Set(tgID, u)
	s.Logger().Info().Int64("user_id", u.ID).Str("user", fmt.Sprintf("%+v", *u)).Msgf("user saved")

	return nil
}
