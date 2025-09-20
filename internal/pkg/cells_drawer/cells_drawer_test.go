package cells_drawer

import (
	"ddd-timer-service/internal/pkg/stats_counter"
	"ddd-timer-service/models"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCellsDrawer_NewCellsImagePNG(t *testing.T) {

	d := NewCellsDrawer()

	user := &models.User{
		ServeFrom: time.Date(2024, time.November, 11, 0, 0, 0, 0, time.Local),
		ServeTo:   time.Date(2025, time.November, 10, 0, 0, 0, 0, time.Local),
	}

	now := time.Date(2025, time.September, 20, 0, 0, 0, 0, time.Local)

	stats, _ := stats_counter.NewStats(user, now)

	data, err := d.NewCellsImagePNG(*stats)
	assert.NoError(t, err)

	err = os.WriteFile("test.png", data, os.ModePerm)
	assert.NoError(t, err)
}
