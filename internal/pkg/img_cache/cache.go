package img_cache

import "sync"

type ImgT uint8

type Cache struct {
	sizeBytes    int
	maxSizeBytes int

	pq *imgQueue

	cache map[string]map[ImgT][]byte // map[key]->map[ImageType]->[]byte
	mx    sync.RWMutex
}

func NewCache(maxCacheSize int) *Cache {
	c := &Cache{
		pq:           new(imgQueue),
		maxSizeBytes: maxCacheSize,
		cache:        make(map[string]map[ImgT][]byte),
	}

	return c
}

func (i *Cache) Set(key string, imgType ImgT, bytes []byte) {
	i.mx.Lock()
	defer i.mx.Unlock()

	if _, ok := i.cache[key]; !ok {
		i.cache[key] = make(map[ImgT][]byte)
	}

	imgSizeBytes := len(bytes)

	// Если картинка слишком большая и не влезает в кеш, значит не будем напрягаться впихивать невпихуемое
	if imgSizeBytes > i.maxSizeBytes {
		return
	}

	// Удаляем старые картинки из кеша до тех пор, пока не освободится место под новую
	for {
		if i.sizeBytes+imgSizeBytes > i.maxSizeBytes {
			i.removeLast()
			continue
		}

		break
	}

	i.cache[key][imgType] = bytes
	i.pq.push(i.pq.encodeKey(key, imgType))
	i.sizeBytes += imgSizeBytes
}

func (i *Cache) Get(key string, imgType ImgT) []byte {
	i.mx.RLock()
	defer i.mx.RUnlock()

	if _, ok := i.cache[key]; !ok {
		return nil
	}

	if _, ok := i.cache[key][imgType]; !ok {
		return nil
	}

	return i.cache[key][imgType]
}

func (i *Cache) Clear() {
	i.mx.Lock()
	defer i.mx.Unlock()

	i.cache = make(map[string]map[ImgT][]byte)
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

	key, imgType, _ := i.pq.decodeKey(lastKey)
	imgSize := len(i.cache[key][imgType])

	delete(i.cache[key], imgType)
	i.sizeBytes -= imgSize
}
