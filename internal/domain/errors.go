package domain

import "errors"

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrNotFound        = errors.New("not found")
	ErrConflict        = errors.New("conflict")
	ErrAlreadyExists   = errors.New("already exists")
	ErrAlreadyReserved = errors.New("item already reserved")
)
