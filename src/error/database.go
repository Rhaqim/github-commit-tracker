package error

import "errors"

var (
	ErrDatabaseConnection = errors.New("database connection error")
)
