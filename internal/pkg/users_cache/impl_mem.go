package users_cache

import (
	"ddd-timer-service/models"
	"sync"
)

type implUsersCacheMem struct {
	mx    sync.RWMutex
	cache map[int64]*models.User
}

func NewImplUsersCacheMem() UsersCache {
	return &implUsersCacheMem{cache: make(map[int64]*models.User)}
}

func (i *implUsersCacheMem) Set(userID int64, u *models.User) {
	i.mx.Lock()
	defer i.mx.Unlock()
	i.cache[userID] = u
}

func (i *implUsersCacheMem) Get(userID int64) *models.User {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.cache[userID]
}

func (i *implUsersCacheMem) Remove(userID int64) {
	i.mx.Lock()
	defer i.mx.Unlock()
	delete(i.cache, userID)
}
