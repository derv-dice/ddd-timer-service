package service

import (
	"context"
	"ddd-timer-service/internal/pkg/tracelog"
	"ddd-timer-service/models"
	"fmt"
	"time"
)

func (s *Service) SaveUser(ctx context.Context, tgID int64, from, to time.Time) error {
	tl, ctx := tracelog.Begin(ctx, "service.SaveUser")
	defer tl.End()

	tl.AddAttributes(tracelog.String("user_id", fmt.Sprint(tgID)))

	if tgID == 0 {
		return fmt.Errorf("tgID is empty")
	}

	u := &models.User{
		ID:        tgID,
		ServeFrom: from,
		ServeTo:   to,
	}

	if err := u.Validate(); err != nil {
		tl.Error(fmt.Errorf("validation user failed: %v", err))
		return err
	}

	if err := s.repo.SaveUser(ctx, u); err != nil {
		tl.AddError(fmt.Errorf("save user failed: %v", err))
		return err
	}

	tl.Info("save user success")
	return nil
}
