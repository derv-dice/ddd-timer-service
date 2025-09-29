package cells_drawer

import (
	"ddd-timer-service/internal/pkg/stats_counter"
	"ddd-timer-service/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCellsDrawer_NewCellsImagePNG(t *testing.T) {
	createTmpDir()
	defer removeTmpDir()

	d := NewCellsDrawer()

	user := &models.User{
		ServeFrom: testDate1,
		ServeTo:   testDate1.AddDate(1, 0, 0),
	}

	now := testDate1.AddDate(0, 5, 11)

	stats, _ := stats_counter.NewStats(user, now)

	data, err := d.PNG(*stats)
	assert.NoError(t, err)

	err = os.WriteFile(pathInTmpDir("cells.png"), data, os.ModePerm)
	assert.NoError(t, err)
}
