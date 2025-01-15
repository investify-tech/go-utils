package must

import (
	"errors"
	"os"
	"testing"
)

func TestVoid(t *testing.T) {
	tests := []struct {
		name          string
		possibleError error
		expectPanic   bool
	}{
		{"nil error", nil, false},
		{"non-nil error", errors.New("test error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.expectPanic {
					t.Errorf("expected panic: %v, got: %v", tt.expectPanic, r != nil)
				}
			}()
			Void(tt.possibleError)
		})
	}
}

func TestAnyType(t *testing.T) {
	tests := []struct {
		name          string
		result        any
		possibleError error
		expectPanic   bool
	}{
		{"valid int", 42, nil, false},
		{"valid string", "hello", nil, false},
		{"error present", nil, errors.New("test error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.expectPanic {
					t.Errorf("expected panic: %v, got: %v", tt.expectPanic, r != nil)
				}
			}()
			result := AnyType(tt.result, tt.possibleError)
			if !tt.expectPanic && result != tt.result {
				t.Errorf("expected result: %v, got: %v", tt.result, result)
			}
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name          string
		result        any
		possibleError error
		expectPanic   bool
	}{
		{"valid int", 10, nil, false},
		{"valid string", "world", nil, false},
		{"error present", nil, errors.New("test error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.expectPanic {
					t.Errorf("expected panic: %v, got: %v", tt.expectPanic, r != nil)
				}
			}()
			result := Any(tt.result, tt.possibleError)
			if !tt.expectPanic && result != tt.result {
				t.Errorf("expected result: %v, got: %v", tt.result, result)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name          string
		result        string
		possibleError error
		expectPanic   bool
	}{
		{"valid string", "go", nil, false},
		{"error present", "", errors.New("test error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.expectPanic {
					t.Errorf("expected panic: %v, got: %v", tt.expectPanic, r != nil)
				}
			}()
			result := String(tt.result, tt.possibleError)
			if !tt.expectPanic && result != tt.result {
				t.Errorf("expected result: %v, got: %v", tt.result, result)
			}
		})
	}
}

func TestAnySlice(t *testing.T) {
	tests := []struct {
		name          string
		result        []int
		possibleError error
		expectPanic   bool
	}{
		{"valid slice", []int{1, 2, 3}, nil, false},
		{"error present", nil, errors.New("test error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.expectPanic {
					t.Errorf("expected panic: %v, got: %v", tt.expectPanic, r != nil)
				}
			}()
			result := AnySlice(tt.result, tt.possibleError)
			if !tt.expectPanic && !equalSlices(result, tt.result) {
				t.Errorf("expected result: %v, got: %v", tt.result, result)
			}
		})
	}
}

func equalSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestEnvVarValue(t *testing.T) {
	tests := []struct {
		name        string
		envVarName  string
		envVarValue string
		setEnv      bool
		expectPanic bool
	}{
		{"env var set", "TEST_VAR", "test_value", true, false},
		{"env var not set", "MISSING_VAR", "", false, true},
		{"env var empty", "EMPTY_VAR", "", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				_ = os.Setenv(tt.envVarName, tt.envVarValue)
			} else {
				_ = os.Unsetenv(tt.envVarName)
			}
			defer func() {
				if r := recover(); (r != nil) != tt.expectPanic {
					t.Errorf("expected panic: %v, got: %v", tt.expectPanic, r != nil)
				}
			}()
			result := EnvVarValue(tt.envVarName)
			if !tt.expectPanic && result != tt.envVarValue {
				t.Errorf("expected result: %v, got: %v", tt.envVarValue, result)
			}
		})
	}
}
