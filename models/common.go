package models

import "errors"

var (
	ErrorNotInitialized = errors.New("object is not initialized")

	ErrorUserNotFound = errors.New("user not found")
)

const OnlyDateLayout = "02.01.2006"

const (
	tmplUserID   = "\nID: %d"
	tmplDateFrom = "\nДата начала службы: %s"
	tmplDateTo   = "\nДата окончания службы: %s"
	tmplPhone    = "\nНомер телефона: %s"
	tmplUsername = "\nИмя: %s"
)

const maxServicePeriodYears = 100 - 18
