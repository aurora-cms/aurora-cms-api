package common

import "time"

type TimeProvider interface {
	Now() time.Time
	Parse(layout, value string) (time.Time, error)
	Format(t time.Time, layout string) string
}
