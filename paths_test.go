package fieldmask

import (
	"reflect"
	"testing"
)

func Test_removeEmptyPaths(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no empty strings",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "some empty strings",
			input:    []string{"a", "", "b", " ", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "all empty strings",
			input:    []string{"", " ", "   "},
			expected: []string{},
		},
		{
			name:     "empty input slice",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeEmptyPaths(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func Test_removeDuplicatePaths(t *testing.T) {
	tests := []struct {
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
			input:    []string{"a", "b", "b", "a", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "all duplicates",
			input:    []string{"a", "a", "a", "a"},
			expected: []string{"a"},
		},
		{
			name:     "empty input slice",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeDuplicatePaths(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
