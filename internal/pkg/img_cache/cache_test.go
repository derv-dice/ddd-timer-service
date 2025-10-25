package img_cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImagesQueue(t *testing.T) {
	q := new(imgQueue)

	key1 := "10_123123123"
	key2 := "11_121287683123"
	key3 := "9_12312315876"

	key := "123123123"
	imgType := ImgT(10)

	encodedKey := q.encodeKey(key, imgType)

	assert.Equal(t, key1, encodedKey)

	dKey, dImgType, err := q.decodeKey(encodedKey)
	assert.NoError(t, err)
	assert.Equal(t, dKey, key)
	assert.Equal(t, imgType, dImgType)

	q.push(key1)
	q.push(key2)
	q.push(key3)

	assert.Equal(t, 3, q.size)

	pKey := q.pop()
	assert.Equal(t, key1, pKey)
	assert.Equal(t, 2, q.size)

	pKey = q.pop()
	assert.Equal(t, key2, pKey)
	assert.Equal(t, 1, q.size)

	pKey = q.pop()
	assert.Equal(t, key3, pKey)
	assert.Equal(t, 0, q.size)

	pKey = q.pop()
	assert.Equal(t, "", pKey)
	assert.Equal(t, 0, q.size)
}
