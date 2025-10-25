package cells_drawer

import (
	"ddd-timer-service/internal/pkg/stats_counter"
	"ddd-timer-service/models"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCellsDrawer_NewCellsImagePNG(t *testing.T) {
	createTmpDir()
	defer removeTmpDir()

	d := NewCellsDrawer()

	user := &models.User{
		ServeFrom: testDate1,
		ServeTo:   testDate1.AddDate(1, 0, 0).AddDate(0, 0, -1),
	}

	now := testDate1.AddDate(0, 5, 11)

	stats, _ := stats_counter.NewStats(user, now)

	data, err := d.PNG(*stats)
	assert.NoError(t, err)

	err = os.WriteFile(pathInTmpDir("cells.png"), data, os.ModePerm)
	assert.NoError(t, err)
}
