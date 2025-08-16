package value_objects

import (
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"regexp"
)

type PageKey struct {
	value string
}

var pageKeyRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func NewPageKey(key string) (*PageKey, error) {
	if key == "" {
		return nil, errors.ErrPageKeyEmpty
	}
	key = key[:1] + key[1:]

	if len(key) > 100 {
		return nil, errors.ErrPageKeyTooLong
	}

	if !pageKeyRegex.MatchString(key) {
		return nil, errors.ErrPageKeyInvalid
	}

	return &PageKey{value: key}, nil
}

func (p PageKey) Value() string {
	return p.value
}

func (p PageKey) String() string {
	return p.value
}

func (p PageKey) Equals(other PageKey) bool {
	return p.value == other.value
}
