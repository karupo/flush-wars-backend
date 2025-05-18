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

// GetMilestones based on xp and streak
func GetMilestones(xp int, streak int) []string {
	var milestones []string
	if streak >= 5 {
		milestones = append(milestones, "5 Days Streak")
	}
	if xp >= 100 {
		milestones = append(milestones, "Level 2")
	}
	return milestones
}
