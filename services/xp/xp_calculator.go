package xp

import (
	"time"
)

// CalculateXP points based on the stool type. If streak gain more points
func CalculateXP(stoolType string, timestamp time.Time, previousLogTime *time.Time) int {
	baseXP := map[string]int{
		"Normal": 20,
		"Soft":   10,
		"Hard":   5,
		"Runny":  5,
	}

	xp := baseXP[stoolType]

	// Example streak bonus: +10 XP if last log was yesterday
	if previousLogTime != nil {
		if timestamp.Sub(*previousLogTime).Hours() <= 48 {
			xp += 10
		}
	}

	return xp
}

// CalculateLevel based on xp
func CalculateLevel(xp int) int {
	return xp / 100
}

// XPForLevel Simple formula example: XP = 100 * level
func ForLevel(level int) int {
	return 100 * level
}
