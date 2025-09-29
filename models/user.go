package models

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID int64

	ServeFrom time.Time
	ServeTo   time.Time
	BirthDate time.Time
}

func (u *User) Validate() error {
	if u.ServeFrom.IsZero() {
		return errors.New("не указана дата начала службы")
	}

	if u.ServeTo.IsZero() {
		return errors.New("не указана дата окончания службы")
	}

	if u.ServeFrom.Round(time.Hour*24) == u.ServeTo.Round(time.Hour*24) {
		return errors.New("дата начала службы не может быть равна дате окончания службы")
	}

	if u.ServeFrom.After(u.ServeTo) {
		return errors.New("дата начала службы не может быть больше даты окончания службы")
	}

	if u.ServeTo.Sub(u.ServeFrom).Hours()/24/365 > maxServicePeriodYears {
		return errors.New("период службы не может быть больше 82 лет")
	}

	return nil
}

func (u *User) String() string {
	return fmt.Sprintf(tmplUserString, u.ID, u.ServeFrom.Format(OnlyDateLayout), u.ServeTo.Format(OnlyDateLayout))
}
