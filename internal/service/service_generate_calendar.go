package service

import (
	"context"
	"time"
)

func (s *Service) GenerateCalendarPNG(_ context.Context, userID int64, withProgressMarks bool) ([]byte, error) {
	u, err := s.GetUser(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	var pngBytes []byte

	if withProgressMarks {
		pngBytes, err = s.calendarDrawer.BySeasonsWithProgressPNG(userID, u.ServeFrom, u.ServeTo, time.Now(), false)
	} else {
		pngBytes, err = s.calendarDrawer.BySeasonsPNG(userID, u.ServeFrom, u.ServeTo, false)
	}

	if err != nil {
		return nil, err
	}

	return pngBytes, nil
}
