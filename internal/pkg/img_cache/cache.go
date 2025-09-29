package img_cache

import "sync"

type ImgT uint8

type Cache struct {
	sizeBytes    int
	maxSizeBytes int

	pq *imgQueue

	cache map[int64]map[ImgT][]byte // map[userID]->map[ImageType]->[]byte
	mx    sync.RWMutex
}

func NewCache(maxCacheSize int) *Cache {
	c := &Cache{
		pq:           new(imgQueue),
		maxSizeBytes: maxCacheSize,
		cache:        make(map[int64]map[ImgT][]byte),
	}

	return c
}

func (i *Cache) Set(userID int64, imgType ImgT, bytes []byte) {
	i.mx.Lock()
	defer i.mx.Unlock()

	if _, ok := i.cache[userID]; !ok {
		i.cache[userID] = make(map[ImgT][]byte)
	}

	imgSize := len(bytes)

	// Если картинка слишком большая и не влезает в кеш, то ничего не делаем
	if imgSize > i.maxSizeBytes {
		return
	}

	// Удаляем старые картинки из кеша до тех пор, пока не освободится место под новую
	for {
		if i.sizeBytes+imgSize > i.maxSizeBytes {
			i.removeLast()
			continue
		}

		break
	}

	i.cache[userID][imgType] = bytes
	i.pq.push(i.pq.encodeKey(userID, imgType))
	i.sizeBytes += imgSize
}

func (i *Cache) Get(userID int64, imgType ImgT) []byte {
	i.mx.RLock()
	defer i.mx.RUnlock()

	if _, ok := i.cache[userID]; !ok {
		return nil
	}

	if _, ok := i.cache[userID][imgType]; !ok {
		return nil
	}

	return i.cache[userID][imgType]
}

func (i *Cache) Clear() {
	i.mx.Lock()
	defer i.mx.Unlock()

	i.cache = make(map[int64]map[ImgT][]byte)
	i.sizeBytes = 0
	i.pq.clear()
}

// TODO методы для сбора метрик Cache

// removeLast - не использовать без i.mx.Lock
func (i *Cache) removeLast() {
	lastKey := i.pq.pop()
	if lastKey == "" {
		return
	}

	userID, imgType, _ := i.pq.decodeKey(lastKey)
	imgSize := len(i.cache[userID][imgType])

	delete(i.cache[userID], imgType)
	i.sizeBytes -= imgSize
}
