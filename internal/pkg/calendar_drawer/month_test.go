package calendar_drawer

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMonth_PNG(t *testing.T) {
	createTmpDir()
	defer removeTmpDir()

	m := NewMonth(2025, 11)

	imgBytes, _, err := m.PNG()
	assert.NoError(t, err)

	err = os.WriteFile(pathInTmpDir("month.png"), imgBytes, os.ModePerm)
	assert.NoError(t, err)

	img2Bytes, _, err := m.progressMask(testDate1, testDate2, time.Now())
	assert.NoError(t, err)

	err = os.WriteFile(pathInTmpDir("monthProgressMask.png"), img2Bytes, os.ModePerm)
	assert.NoError(t, err)
}

func TestMonth_StringsOpacityMask(t *testing.T) {
	m := NewMonth(2025, 11)

	arr := m.stringsOpacityMask(testDate1, testDate2, time.Now())
	for _, row := range arr {
		fmt.Println(row)
	}
}
