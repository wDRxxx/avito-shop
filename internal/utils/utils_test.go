package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCustomDuration(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect time.Duration
		hasErr bool
	}{
		{
			name:   "valid hours",
			input:  "2h",
			expect: 2 * time.Hour,
			hasErr: false,
		},
		{
			name:   "valid minutes",
			input:  "30m",
			expect: 30 * time.Minute,
			hasErr: false,
		},
		{
			name:   "valid days",
			input:  "3d",
			expect: 3 * 24 * time.Hour,
			hasErr: false,
		},
		{
			name:   "invalid format",
			input:  "10x",
			expect: 0,
			hasErr: true,
		},
		{
			name:   "empty string",
			input:  "",
			expect: 0,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration, err := ParseCustomDuration(tt.input)
			if tt.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expect, duration)
			}
		})
	}
}
