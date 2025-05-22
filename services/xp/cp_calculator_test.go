package xp

import (
	"testing"
	"time"
)

func TestCalculateXP(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	old := now.Add(-72 * time.Hour)

	tests := []struct {
		name            string
		stoolType       string
		timestamp       time.Time
		previousLogTime *time.Time
		expectedXP      int
	}{
		{
			name:            "Normal with streak",
			stoolType:       "Normal",
			timestamp:       now,
			previousLogTime: &yesterday,
			expectedXP:      30, // 20 base + 10 streak
		},
		{
			name:            "Soft no streak",
			stoolType:       "Soft",
			timestamp:       now,
			previousLogTime: &old,
			expectedXP:      10, // 10 base
		},
		{
			name:            "Hard with streak",
			stoolType:       "Hard",
			timestamp:       now,
			previousLogTime: &yesterday,
			expectedXP:      15, // 5 base + 10 streak
		},
		{
			name:            "Runny no previous log",
			stoolType:       "Runny",
			timestamp:       now,
			previousLogTime: nil,
			expectedXP:      5, // 5 base
		},
		{
			name:            "Unknown type",
			stoolType:       "Alien",
			timestamp:       now,
			previousLogTime: nil,
			expectedXP:      0, // not found in map
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateXP(tt.stoolType, tt.timestamp, tt.previousLogTime)
			if got != tt.expectedXP {
				t.Errorf("expected %d, got %d", tt.expectedXP, got)
			}
		})
	}
}
