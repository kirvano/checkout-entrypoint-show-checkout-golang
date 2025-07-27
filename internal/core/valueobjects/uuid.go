package valueobjects

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

var (
	ErrInvalidUUID = errors.New("invalid UUID format")
	uuidRegex      = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

type UUID struct {
	value string
}

func NewUUID(value string) (*UUID, error) {
	if !IsValidUUID(value) {
		return nil, ErrInvalidUUID
	}
	return &UUID{value: value}, nil
}

func NewRandomUUID() *UUID {
	return &UUID{value: uuid.New().String()}
}

func (u *UUID) String() string {
	return u.value
}

func (u *UUID) Value() string {
	return u.value
}

func IsValidUUID(value string) bool {
	return uuidRegex.MatchString(value)
}
