package repository

import "errors"

var (
	ErrNotFound       = errors.New("record not found")
	ErrDuplicate      = errors.New("duplicate key violation")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUserNameExists = errors.New("username already exists")
	ErrInviteCodeUsed = errors.New("invite code already used")
	ErrUnexpected     = errors.New("unexpected error in repository")
)
