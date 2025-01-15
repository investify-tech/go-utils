package str_test

import (
	"github.com/investify-tech/go-utils/str"
	"reflect"
	"testing"
)

func TestSortMapOfStrings(test *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]string
		expected map[string]string
	}{
		{
			name:     "Test case for an empty map",
			input:    map[string]string{},
			expected: map[string]string{},
		},
		{
			name:     "Test case for a map with one element",
			input:    map[string]string{"one": "1"},
			expected: map[string]string{"one": "1"},
		},
		{
			name:     "Test case for a map with multiple elements",
			input:    map[string]string{"one": "1", "two": "2", "three": "3"},
			expected: map[string]string{"one": "1", "three": "3", "two": "2"},
		},
		{name: "Test case for a map with duplicate values",
			input:    map[string]string{"one": "1", "two": "1", "three": "1"},
			expected: map[string]string{"one": "1", "three": "1", "two": "1"},
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.name, func(t *testing.T) {
			original := testCase.input
			expected := testCase.expected
			str.SortMapOfStrings(original)
			if !reflect.DeepEqual(original, expected) {
				t.Errorf("Expected %v, got %v", expected, original)
			}
		})
	}
}

func TestDeduplicateSliceOfStrings(test *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "some duplicates",
			input:    []string{"a", "d", "c", "d", "a"},
			expected: []string{"a", "d", "c"},
		},
		{
			name:     "all duplicates",
			input:    []string{"a", "a", "a", "a"},
			expected: []string{"a"},
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.name, func(t *testing.T) {
			result := str.DeduplicateSliceOfStrings(testCase.input)
			if !reflect.DeepEqual(result, testCase.expected) {
				t.Errorf("Expected %v but got %v", testCase.expected, result)
			}
		})
	}
}

func TestHashAndTrim(test *testing.T) {
	testCases := []struct {
		name        string
		inputStr    string
		inputLength int
		expected    string
	}{
		{
			name:        "simple test",
			inputStr:    "investify-gitlab-ob-bucket-2024",
			inputLength: 7,
			expected:    "5fa7a81",
		},
		{
			name:        "empty string",
			inputStr:    "",
			inputLength: 5,
			expected:    "d41d8",
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.name, func(t *testing.T) {
			result := str.HashAndTrim(testCase.inputStr, testCase.inputLength)
			if result != testCase.expected {
				t.Errorf("Expected %v but got %v", testCase.expected, result)
			}
		})
	}
}
