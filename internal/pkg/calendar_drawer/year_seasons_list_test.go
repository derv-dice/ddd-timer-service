package calendar_drawer

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestYearSeasonsSet_PNG(t *testing.T) {
	createTmpDir()
	defer removeTmpDir()

	now := time.Now()

	//seasons := NewCalendar(testDate1, testDate2).Seasons()
	seasons := NewCalendar(now, now.AddDate(2, 0, 0)).Seasons()

	imgBytes, _, err := seasons.PNG()
	assert.NoError(t, err)

	err = os.WriteFile(pathInTmpDir("seasons.png"), imgBytes, os.ModePerm)
	assert.NoError(t, err)

	img2Bytes, _, err := seasons.PNGWithProgressMask(testDate1, testDate2, time.Now())
	assert.NoError(t, err)

	err = os.WriteFile(pathInTmpDir("seasonsMasked.png"), img2Bytes, os.ModePerm)
	assert.NoError(t, err)
}
