package time_provider

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"testing"
	"time"
)

func TestStandardTimeProvider_Now(t *testing.T) {
	provider := StandardTimeProvider{}

	before := time.Now()
	now := provider.Now()
	after := time.Now()

	if now.Before(before) || now.After(after) {
		t.Errorf("Now() = %v, expected time to be between %v and %v", now, before, after)
	}
}

func TestStandardTimeProvider_Parse(t *testing.T) {
	provider := StandardTimeProvider{}

	tests := []struct {
		name    string
		layout  string
		value   string
		want    time.Time
		wantErr bool
	}{
		{"Valid ISO8601", time.RFC3339, "2025-07-27T12:00:00Z", time.Date(2025, 7, 27, 12, 0, 0, 0, time.UTC), false},
		{"Invalid Date", time.RFC3339, "invalid-date", time.Time{}, true},
		{"Empty Value", time.RFC3339, "", time.Time{}, true},
		{"Empty Layout", "", "2025-07-27T12:00:00Z", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := provider.Parse(tt.layout, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestNewStandardTimeProvider(t *testing.T) {
	provider := NewStandardTimeProvider()

	// Verify the returned type implements the TimeProvider interface
	_, ok := provider.(common.TimeProvider)
	if !ok {
		t.Errorf("NewStandardTimeProvider() = %T, expected to implement common.TimeProvider", provider)
	}

	// Validate the behavior of the Now method
	before := time.Now()
	now := provider.Now()
	after := time.Now()

	if now.Before(before) || now.After(after) {
		t.Errorf("Now() = %v, expected time to be between %v and %v", now, before, after)
	}
}
func TestStandardTimeProvider_Format(t *testing.T) {
	provider := StandardTimeProvider{}

	tests := []struct {
		name   string
		time   time.Time
		layout string
		want   string
	}{
		{"ISO8601 Format", time.Date(2025, 7, 27, 12, 0, 0, 0, time.UTC), time.RFC3339, "2025-07-27T12:00:00Z"},
		{"ANSIC Format", time.Date(2025, 7, 27, 12, 0, 0, 0, time.UTC), time.ANSIC, "Sun Jul 27 12:00:00 2025"},
		{"Empty Layout", time.Date(2025, 7, 27, 12, 0, 0, 0, time.UTC), "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := provider.Format(tt.time, tt.layout)
			if got != tt.want {
				t.Errorf("Format() got = %v, want %v", got, tt.want)
			}
		})
	}
}
