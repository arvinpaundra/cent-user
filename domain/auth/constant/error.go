package constant

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrEmailAlreadyTaken    = errors.New("email has already been taken")
	ErrWrongEmailOrPassword = errors.New("wrong email or password")
)
