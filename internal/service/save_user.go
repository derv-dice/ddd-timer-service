package service

import (
	"context"
	"ddd-timer-service/internal/pkg/tracelog"
	"ddd-timer-service/models"
	"fmt"
)

func (s *Service) SaveUser(ctx context.Context, u *models.User) error {
	tl, ctx := tracelog.Begin(ctx, "service.SaveUser")
	defer tl.End()

	if err := u.Validate(); err != nil {
		tl.Error(fmt.Errorf("validation user failed: %v", err))
		return err
	}

	tl.AddAttributes(tracelog.String("user_id", fmt.Sprint(u.ID)))

	if err := s.repo.SaveUser(ctx, u); err != nil {
		tl.AddError(fmt.Errorf("save user failed: %v", err))
		return err
	}

	tl.Info("save user success")
	return nil
}
