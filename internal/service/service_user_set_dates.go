package service

import (
	"context"
	"ddd-timer-service/models"
	"fmt"
	"time"
)

func (s *Service) SetUserDatesFromStringMessage(ctx context.Context, userID int64, dateFrom, dateTo string) error {
	from, err := time.Parse(models.OnlyDateLayout, dateFrom)
	if err != nil {
		return fmt.Errorf("некорректная дата: '%s'", dateFrom)
	}

	to, err := time.Parse(models.OnlyDateLayout, dateTo)
	if err != nil {
		return fmt.Errorf("некорректная дата: '%s'", dateTo)
	}

	return s.SaveUser(ctx, userID, from, to)
}
