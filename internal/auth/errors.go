package auth

import "errors"

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrDefaultRoleNotFound = errors.New("default user role not found")
	ErrInvalidCredentials  = errors.New("invalid email or password")
)