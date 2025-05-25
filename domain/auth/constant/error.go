package constant

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrSessionNotFound      = errors.New("session not found")
	ErrTokenInvalid         = errors.New("token invalid")
	ErrEmailAlreadyTaken    = errors.New("email has already been taken")
	ErrWrongEmailOrPassword = errors.New("wrong email or password")
)
