package time_provider

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"time"
)

type StandardTimeProvider struct{}

// NewStandardTimeProvider initializes a new StandardTimeProvider instance.
// It returns a pointer to the StandardTimeProvider struct.
func NewStandardTimeProvider() common.TimeProvider {
	return &StandardTimeProvider{}
}

func (s StandardTimeProvider) Now() time.Time {
	return time.Now()
}

func (s StandardTimeProvider) Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

func (s StandardTimeProvider) Format(t time.Time, layout string) string {
	return t.Format(layout)
}
