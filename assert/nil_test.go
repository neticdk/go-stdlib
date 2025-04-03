package assert_test

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestNil(t *testing.T) {
	tests := []struct {
		name     string
		got      any
		wantPass bool
	}{
		{"nil interface", nil, true},
		{"nil pointer", (*int)(nil), true},
		{"nil slice", ([]int)(nil), true},
		{"empty slice", []int{}, false},
		{"nil map", (map[string]int)(nil), true},
		{"empty map", map[string]int{}, false},
		{"zero int", 0, false},
		{"non-nil int", 1, false},
		{"zero string", "", false},
		{"non-nil string", "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			pass := assert.Nil(mockT, tt.got)
			if pass != tt.wantPass {
				t.Errorf("Nil() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}

			if mockT.Failed() != !tt.wantPass {
				t.Errorf("Nil() Failed() result = %v, want %v", mockT.Failed(), !tt.wantPass)
			}

			if !tt.wantPass && len(mockT.errorMessages) == 0 {
				t.Errorf("Nil() expected error message, but got none")
			}
			if tt.wantPass && len(mockT.errorMessages) > 0 {
				t.Errorf("Nil() expected no error message, but got: %v", mockT.errorMessages)
			}
		})
	}
}

func TestNotNil(t *testing.T) {
	tests := []struct {
		name     string
		got      any
		wantPass bool
	}{
		{"nil interface", nil, false},
		{"nil pointer", (*int)(nil), false},
		{"nil slice", ([]int)(nil), false},
		{"empty slice", []int{}, true},
		{"nil map", (map[string]int)(nil), false},
		{"empty map", map[string]int{}, true},
		{"zero int", 0, true},
		{"non-nil int", 1, true},
		{"zero string", "", true},
		{"non-nil string", "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			pass := assert.NotNil(mockT, tt.got)
			if pass != tt.wantPass {
				t.Errorf("NotNil() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}

			if mockT.Failed() != !tt.wantPass {
				t.Errorf("NotNil() Failed() result = %v, want %v", mockT.Failed(), !tt.wantPass)
			}

			if !tt.wantPass && len(mockT.errorMessages) == 0 {
				t.Errorf("NotNil() expected error message, but got none")
			}
			if tt.wantPass && len(mockT.errorMessages) > 0 {
				t.Errorf("NotNil() expected no error message, but got: %v", mockT.errorMessages)
			}
		})
	}
}
