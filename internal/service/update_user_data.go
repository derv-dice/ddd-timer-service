package service

import (
	"context"
	"ddd-timer-service/internal/pkg/tracelog"
)

func (s *Service) UpdateUserData(ctx context.Context, userID int64, username, name string) {
	tl, ctx := tracelog.Begin(ctx, "service/UpdateUserData")
	defer tl.End()

	if name == "" && username == "" {
		return
	}

	u, err := s.repo.LoadUser(ctx, userID)
	if err != nil {
		tl.AddError(err)
		return
	}

	if username != "" && u.Username != username {
		u.Username = username
	}

	if name != "" && u.Name != name {
		u.Name = name
	}

	if err = s.repo.SaveUser(ctx, u); err != nil {
		tl.AddError(err)
		return
	}
}
