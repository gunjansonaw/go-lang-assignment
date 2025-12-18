package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "Birthday this year already passed",
			dob:      time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			expected: calculateExpectedAge(time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:     "Birthday this year not yet passed",
			dob:      time.Date(2000, 12, 25, 0, 0, 0, 0, time.UTC),
			expected: calculateExpectedAge(time.Date(2000, 12, 25, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:     "Birthday today",
			dob:      time.Now().AddDate(-25, 0, 0),
			expected: 25,
		},
		{
			name:     "Born this year",
			dob:      time.Now().AddDate(0, -6, 0),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			age := CalculateAge(tt.dob)
			if age != tt.expected {
				t.Errorf("CalculateAge() = %v, want %v", age, tt.expected)
			}
		})
	}
}

func calculateExpectedAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

