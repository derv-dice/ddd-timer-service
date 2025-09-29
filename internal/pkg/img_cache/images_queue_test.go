package img_cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	imgType1 = iota + 1
	imgType2
)

func TestImagesCache(t *testing.T) {
	c := NewCache(100)

	img1 := make([]byte, 10)
	img2 := make([]byte, 50)
	img3 := make([]byte, 150)
	img4 := make([]byte, 100)
	img5 := make([]byte, 101)

	id1 := int64(100001)
	id2 := int64(100002)
	id3 := int64(100003)

	c.Set(id1, imgType1, img1)

	assert.Equal(t, 10, c.sizeBytes)
	c.Set(id2, imgType1, img2)
	assert.Equal(t, 60, c.sizeBytes)
	c.Set(id3, imgType1, img3)
	assert.Equal(t, 60, c.sizeBytes)
	c.Set(id3, imgType2, img4)
	assert.Equal(t, 100, c.sizeBytes)

	c.Clear()

	// Большая картинка не влезет в кеш
	c.Set(id2, imgType2, img5)
	assert.Equal(t, 0, c.sizeBytes)

	// А маленькая да
	c.Set(id3, imgType2, img4)
	assert.Equal(t, 100, c.sizeBytes)

	c.Set(id2, imgType1, img2)
	assert.Equal(t, 50, c.sizeBytes)
}
