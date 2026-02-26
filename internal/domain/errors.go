package domain

import "errors"

var (
	ErrPhotoNotFound   = errors.New("photo not found")
	ErrInvalidPhoto    = errors.New("invalid photo data")
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidSession  = errors.New("invalid session data")
	ErrDeviceNotFound  = errors.New("device not found")
	ErrInvalidDevice   = errors.New("invalid device data")
)
