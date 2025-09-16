package repository

import (
	"context"
	"ddd-timer-service/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestImplSQLiteRepository(t *testing.T) {
	repo, err := NewSQLiteRepository("./testdata/test_db.sqlite", true)
	if err != nil {
		t.Fatalf("NewSQLiteRepository: %v", err)
	}

	id := int64(100)
	sf := time.Now().Add(-10 * time.Hour).Round(time.Second)
	st := time.Now().Add(10 * time.Hour).Round(time.Second)
	tb := st.Add(-10 * time.Hour).Round(time.Second)
	tb2 := st.Add(-20 * time.Hour).Round(time.Second)

	user := &models.User{
		ID:        id,
		ServeFrom: sf,
		ServeTo:   st,
		BirthDate: tb,
	}

	t.Run("SaveUser", func(t *testing.T) {
		err = repo.SaveUser(context.Background(), user)
		assert.NoError(t, err)

		user2, err := repo.LoadUser(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, user, user2)
	})

	t.Run("Edit and save existing user", func(t *testing.T) {
		user.BirthDate = tb2
		err = repo.SaveUser(context.Background(), user)
		assert.NoError(t, err)

		user3, err := repo.LoadUser(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, tb2, user3.BirthDate)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		err = repo.DeleteUser(context.Background(), id)
		assert.NoError(t, err)

		user4, err := repo.LoadUser(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, user4)
	})

	t.Run("LoadAllUsers", func(t *testing.T) {
		user5 := &models.User{
			ID:        id + 1,
			ServeFrom: sf,
			ServeTo:   st,
			BirthDate: tb,
		}

		err = repo.SaveUser(context.Background(), user5)
		assert.NoError(t, err)

		user6 := &models.User{
			ID:        id + 2,
			ServeFrom: sf,
			ServeTo:   st,
			BirthDate: tb,
		}
		err = repo.SaveUser(context.Background(), user6)
		assert.NoError(t, err)

		var users []*models.User
		users, err = repo.LoadAllUsers(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 2, len(users))

		err = repo.DeleteUser(context.Background(), user5.ID)
		assert.NoError(t, err)

		err = repo.DeleteUser(context.Background(), user6.ID)
		assert.NoError(t, err)
	})
}
