package user

import "errors"

var (

	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidUserID      = errors.New("invalid user id")
	 ErrAdminCannotDeleteSelf = errors.New("admin cannot delete own account")
)