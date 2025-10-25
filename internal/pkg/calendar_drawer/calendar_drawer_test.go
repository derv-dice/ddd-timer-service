package calendar_drawer

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CalendarDrawerCaching(t *testing.T) {
	cd := NewCalendarDrawer(context.TODO(), 2, 50)

	t1 := time.Now()
	imgBytes, err := cd.BySeasonsWithProgressPNG(testDate1, testDate1.AddDate(0, 0, 1), time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	d1 := time.Since(t1)
	t.Logf("no cache dur: %s", d1)

	err = os.WriteFile(pathInTmpDir("calendar.png"), imgBytes, os.ModePerm)
	assert.NoError(t, err)
	t2 := time.Now()
	imgBytes2, err := cd.BySeasonsWithProgressPNG(testDate1, testDate1.AddDate(0, 0, 1), time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	d2 := time.Since(t2)
	t.Logf("with cache dur: %s", d2)

	cd.cache.Clear()

	t3 := time.Now()
	imgBytes3, err := cd.BySeasonsWithProgressPNG(testDate1, testDate1.AddDate(0, 0, 1), time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	d3 := time.Since(t3)
	t.Logf("after clear dur: %s", d3)

	assert.Greater(t, d1, d2)
	assert.Greater(t, d3, d2)
	assert.Equal(t, imgBytes2, imgBytes)
	assert.Equal(t, imgBytes3, imgBytes)
}
