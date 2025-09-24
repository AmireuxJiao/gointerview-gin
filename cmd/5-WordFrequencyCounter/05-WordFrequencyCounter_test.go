package main

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestCountWordFrequency(t *testing.T) {
	testCases := []struct {
		input    string
		expected map[string]int
	}{
		{
			input: "The quick brown fox jumps over the lazy dog.",
			expected: map[string]int{
				"the":   2,
				"quick": 1,
				"brown": 1,
				"fox":   1,
				"jumps": 1,
				"over":  1,
				"lazy":  1,
				"dog":   1,
			},
		},
		{
			input: "Hello, hello! How are you doing today? Today is a great day.",
			expected: map[string]int{
				"hello": 2,
				"how":   1,
				"are":   1,
				"you":   1,
				"doing": 1,
				"today": 2,
				"is":    1,
				"a":     1,
				"great": 1,
				"day":   1,
			},
		},
		{
			input: "Go, go, go! Let's learn Go programming.",
			expected: map[string]int{
				"go":          4,
				"lets":        1,
				"learn":       1,
				"programming": 1,
			},
		},
		{
			input: "  Spaces,   tabs,\t\tand\nnew-lines are ignored!  ",
			expected: map[string]int{
				"spaces":  1,
				"tabs":    1,
				"and":     1,
				"new":     1,
				"lines":   1,
				"are":     1,
				"ignored": 1,
			},
		},
		{
			input: "Numbers123 and456 mixed789 content",
			expected: map[string]int{
				"numbers123": 1,
				"and456":     1,
				"mixed789":   1,
				"content":    1,
			},
		},
		{
			input:    "",
			expected: map[string]int{},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test Case #%d", i+1), func(t *testing.T) {
			result := CountWordFrequency(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v but got %v", formatMap(tc.expected), formatMap(result))
			}
		})
	}
}

// Helper function to format maps for better error messages
func formatMap(m map[string]int) string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := "{\n"
	for _, k := range keys {
		result += fmt.Sprintf("  %q: %d,\n", k, m[k])
	}
	result += "}"
	return result
}
