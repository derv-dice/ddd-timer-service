package service

import (
	"context"
	"ddd-timer-service/internal/pkg/cells_drawer"
	"fmt"
)

func (s *Service) GenerateCellsPNG(_ context.Context, userID int64) ([]byte, error) {
	stats, err := s.GetUserStats(context.Background(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %v", err)
	}

	img, err := cells_drawer.NewCellsDrawer().NewCellsImagePNG(*stats)
	if err != nil {
		return nil, fmt.Errorf("failed to generate cells png: %v", err)
	}

	return img, nil
}
