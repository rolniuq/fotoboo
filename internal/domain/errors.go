package domain

import "errors"

var (
	ErrPhotoNotFound = errors.New("photo not found")
	ErrInvalidPhoto  = errors.New("invalid photo data")
)
