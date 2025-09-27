package service

import (
	"context"
	"ddd-timer-service/internal/pkg/calendar_drawer"
	"time"
)

func (s *Service) GenerateCalendarPNG(_ context.Context, userID int64, withProgressMarks bool) ([]byte, error) {
	u, err := s.GetUser(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	var pngBytes []byte

	if withProgressMarks {
		pngBytes, _, err = calendar_drawer.NewCalendarDrawer().BySeasonsWithProgressPNG(u.ServeFrom, u.ServeTo, time.Now())
	} else {
		pngBytes, _, err = calendar_drawer.NewCalendarDrawer().BySeasonsPNG(u.ServeFrom, u.ServeTo)
	}

	if err != nil {
		return nil, err
	}

	return pngBytes, nil
}
