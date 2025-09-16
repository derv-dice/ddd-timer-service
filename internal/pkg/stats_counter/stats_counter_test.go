package stats_counter

import (
	"ddd-timer-service/models"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	d1 := time.Now()
	d2 := d1.Add(time.Hour * 24 * 364)

	now := d1.Add(time.Hour * 24 * 309)

	u := &models.User{
		ID:        10,
		ServeFrom: d1,
		ServeTo:   d2,
	}

	stats, err := NewStats(u, now)
	assert.NoError(t, err)

	assert.Equal(t, math.Round(stats.PassedHours()*100)/100, 7416.00)
	assert.Equal(t, math.Round(stats.PassedDays()*100)/100, 309.00)
	assert.Equal(t, math.Round(stats.PassedWeeks()*100)/100, 44.14)
	assert.Equal(t, math.Round(stats.PassedPercents()*100)/100, 84.89)
	assert.Equal(t, math.Round(stats.LeftHours()*100)/100, 1320.00)
	assert.Equal(t, math.Round(stats.LeftDays()*100)/100, 55.00)
	assert.Equal(t, math.Round(stats.LeftWeeks()*100)/100, 7.86)
	assert.Equal(t, math.Round(stats.LeftPercents()*100)/100, 15.11)
}
