package dbx

import (
	"errors"
	"reflect"
	"testing"
)

type SampleStruct struct {
	Name  *string
	Age   *int
	Email *string
}

type AnotherStruct struct {
	ID   *int
	Code *string
}

func TestApplyPatch(t *testing.T) {
	tests := []struct {
		name       string
		patch      any
		target     any
		expectErr  error
		finalState any
	}{
		{
			name: "valid patch with changes",
			patch: &SampleStruct{
				Name:  ptr("John"),
				Age:   ptr(30),
				Email: nil,
			},
			target: &SampleStruct{
				Name:  ptr("Doe"),
				Age:   ptr(25),
				Email: ptr("doe@example.com"),
			},
			expectErr: nil,
			finalState: &SampleStruct{
				Name:  ptr("John"),
				Age:   ptr(30),
				Email: ptr("doe@example.com"),
			},
		},
		{
			name: "valid patch with no overlaps",
			patch: &SampleStruct{
				Name:  nil,
				Age:   nil,
				Email: nil,
			},
			target: &SampleStruct{
				Name:  ptr("Doe"),
				Age:   ptr(25),
				Email: ptr("doe@example.com"),
			},
			expectErr: nil,
			finalState: &SampleStruct{
				Name:  ptr("Doe"),
				Age:   ptr(25),
				Email: ptr("doe@example.com"),
			},
		},
		{
			name: "invalid patch type",
			patch: &AnotherStruct{
				ID:   ptr(123),
				Code: ptr("xyz"),
			},
			target: &SampleStruct{
				Name:  ptr("Doe"),
				Age:   ptr(25),
				Email: ptr("doe@example.com"),
			},
			expectErr: errors.New("patch and target must be of the same type"),
			finalState: &SampleStruct{
				Name:  ptr("Doe"),
				Age:   ptr(25),
				Email: ptr("doe@example.com"),
			},
		},
		{
			name: "patch not pointer",
			patch: SampleStruct{
				Name:  ptr("John"),
				Age:   ptr(30),
				Email: nil,
			},
			target: &SampleStruct{
				Name:  ptr("Doe"),
				Age:   ptr(25),
				Email: ptr("doe@example.com"),
			},
			expectErr: errors.New("patch must be a pointer"),
			finalState: &SampleStruct{
				Name:  ptr("Doe"),
				Age:   ptr(25),
				Email: ptr("doe@example.com"),
			},
		},
		{
			name: "target not pointer",
			patch: &SampleStruct{
				Name:  ptr("John"),
				Age:   ptr(30),
				Email: nil,
			},
			target: SampleStruct{
				Name:  ptr("Doe"),
				Age:   ptr(25),
				Email: ptr("doe@example.com"),
			},
			expectErr: errors.New("target must be a pointer"),
			finalState: SampleStruct{
				Name:  ptr("Doe"),
				Age:   ptr(25),
				Email: ptr("doe@example.com"),
			},
		},
		{
			name: "patch and target both not pointer",
			patch: SampleStruct{
				Name:  ptr("John"),
				Age:   ptr(30),
				Email: nil,
			},
			target: SampleStruct{
				Name:  ptr("Doe"),
				Age:   ptr(25),
				Email: ptr("doe@example.com"),
			},
			expectErr: errors.New("patch must be a pointer"),
			finalState: SampleStruct{
				Name:  ptr("Doe"),
				Age:   ptr(25),
				Email: ptr("doe@example.com"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ApplyPatch(tc.patch, tc.target)
			if tc.expectErr != nil {
				if err == nil || err.Error() != tc.expectErr.Error() {
					t.Fatalf("expected error %v, got %v", tc.expectErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if !reflect.DeepEqual(tc.target, tc.finalState) {
					t.Fatalf("expected target %v, got %v", tc.finalState, tc.target)
				}
			}
		})
	}
}

func ptr[T any](value T) *T {
	return &value
}
