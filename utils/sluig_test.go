package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello-world-123", true},
		{"hello_world", false},
		{"123", true},
		{"abc123", true},
		{"abc-def_123", false},
	}

	for _, test := range tests {
		result := IsSlug(test.input)
		require.Equal(t, test.expected, result)
	}
}
