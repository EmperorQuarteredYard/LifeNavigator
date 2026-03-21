package repository

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("record not found")
	ErrDuplicate        = errors.New("duplicate key violation")
	ErrInvalidInput     = errors.New("invalid input")
	ErrRecordExist      = errors.New("record already exists")
	ErrInviteCodeUsed   = errors.New("invite code already used")
	ErrUnexpected       = errors.New("unexpected error in repository")
	ErrPermissionDenied = errors.New("permission denied")
	ErrConcurrentUpdate = errors.New("concurrent update")
)
