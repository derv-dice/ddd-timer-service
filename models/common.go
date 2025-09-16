package models

import "errors"

var (
	ErrorNotInitialized = errors.New("object is not initialized")

	ErrorUserNotFound = errors.New("user not found")
)

const OnlyDateLayout = "02.01.2006"

const tmplUserString = `
ID: %d
Дата начала службы: %s
Дата окончания службы: %s`
