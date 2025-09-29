package users_cache

import "ddd-timer-service/models"

type UsersCache interface {
	Set(id int64, value *models.User)
	Get(id int64) *models.User
	Remove(id int64)
}
