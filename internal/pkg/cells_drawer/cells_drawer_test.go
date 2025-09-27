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
	createTmpDir()
	defer removeTmpDir()

	d := NewCellsDrawer()

	user := &models.User{
		ServeFrom: testDate1,
		ServeTo:   testDate2,
	}

	now := time.Date(2025, time.September, 20, 0, 0, 0, 0, time.Local)

	stats, _ := stats_counter.NewStats(user, now)

	data, err := d.NewCellsImagePNG(*stats)
	assert.NoError(t, err)

	err = os.WriteFile(pathInTmpDir("test.png"), data, os.ModePerm)
	assert.NoError(t, err)
}
