package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Тестовые константы
var (
	validBirthDate = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	validServeFrom = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	validServeTo   = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	tooLongServeTo = time.Date(2120, 1, 1, 0, 0, 0, 0, time.UTC) // 100 лет спустя
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		wantErr  bool
		errorMsg string
	}{
		{
			name: "валидный пользователь",
			user: User{
				ID:        1,
				ServeFrom: validServeFrom,
				ServeTo:   validServeTo,
				BirthDate: validBirthDate,
			},
			wantErr: false,
		},
		{
			name: "не указана дата начала службы",
			user: User{
				ID:        1,
				ServeFrom: time.Time{}, // zero value
				ServeTo:   validServeTo,
				BirthDate: validBirthDate,
			},
			wantErr:  true,
			errorMsg: "не указана дата начала службы",
		},
		{
			name: "не указана дата окончания службы",
			user: User{
				ID:        1,
				ServeFrom: validServeFrom,
				ServeTo:   time.Time{}, // zero value
				BirthDate: validBirthDate,
			},
			wantErr:  true,
			errorMsg: "не указана дата окончания службы",
		},
		{
			name: "дата начала равна дате окончания",
			user: User{
				ID:        1,
				ServeFrom: validServeFrom,
				ServeTo:   validServeFrom,
				BirthDate: validBirthDate,
			},
			wantErr:  true,
			errorMsg: "дата начала службы не может быть равна дате окончания службы",
		},
		{
			name: "дата начала после даты окончания",
			user: User{
				ID:        1,
				ServeFrom: validServeTo,
				ServeTo:   validServeFrom,
				BirthDate: validBirthDate,
			},
			wantErr:  true,
			errorMsg: "дата начала службы не может быть больше даты окончания службы",
		},
		{
			name: "период службы больше максимального",
			user: User{
				ID:        1,
				ServeFrom: validServeFrom,
				ServeTo:   tooLongServeTo,
				BirthDate: validBirthDate,
			},
			wantErr:  true,
			errorMsg: "период службы не может быть больше 82 лет",
		},
		{
			name: "период службы чуть больше 82 лет",
			user: User{
				ID:        1,
				ServeFrom: validServeFrom,
				ServeTo:   validServeFrom.AddDate(maxServicePeriodYears, 0, 1), // +1 день
				BirthDate: validBirthDate,
			},
			wantErr:  true,
			errorMsg: "период службы не может быть больше 82 лет",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Equal(t, tt.errorMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
