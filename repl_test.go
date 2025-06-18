package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello  world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  Go is  awesome!",
			expected: []string{"go", "is", "awesome!"},
		},
		{
			input:    "  REPL  testing  ",
			expected: []string{"repl", "testing"},
		},
		{
			input:    "  Multiple   spaces   ",
			expected: []string{"multiple", "spaces"},
		},
		{
			input:    "  Leading and trailing spaces  ",
			expected: []string{"leading", "and", "trailing", "spaces"},
		},
		{
			input:    "  Special characters! @#  ",
			expected: []string{"special", "characters!", "@#"},
		},
		{input: "Mixed   Case  Input",
			expected: []string{"mixed", "case", "input"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		result := cleanInput(c.input)
		if len(result) != len(c.expected) {
			t.Errorf("cleanInput(%q) = %v; expected %v", c.input, result, c.expected)
			continue
		}
		for i := range result {
			word := result[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%q)[%d] = %q; expected %q", c.input, i, word, expectedWord)
			}
		}
	}
}
