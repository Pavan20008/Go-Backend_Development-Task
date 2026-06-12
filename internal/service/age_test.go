package service

import (
	"testing"
	"time"
)

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name      string
		dob       time.Time
		reference time.Time
		want      int
	}{
		{
			name:      "birthday already passed this year",
			dob:       date(1990, time.May, 10),
			reference: date(2025, time.June, 12),
			want:      35,
		},
		{
			name:      "birthday not yet reached this year",
			dob:       date(1990, time.May, 10),
			reference: date(2025, time.April, 1),
			want:      34,
		},
		{
			name:      "exactly on birthday",
			dob:       date(2000, time.January, 1),
			reference: date(2025, time.January, 1),
			want:      25,
		},
		{
			name:      "day before birthday",
			dob:       date(2000, time.March, 15),
			reference: date(2025, time.March, 14),
			want:      24,
		},
		{
			name:      "newborn",
			dob:       date(2025, time.June, 12),
			reference: date(2025, time.June, 12),
			want:      0,
		},
		{
			name:      "leap day birthday on non-leap year",
			dob:       date(2000, time.February, 29),
			reference: date(2025, time.February, 28),
			want:      24,
		},
		{
			name:      "leap day birthday on leap year",
			dob:       date(2000, time.February, 29),
			reference: date(2024, time.February, 29),
			want:      24,
		},
		{
			name:      "future date of birth returns zero",
			dob:       date(2030, time.January, 1),
			reference: date(2025, time.June, 12),
			want:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateAge(tt.dob, tt.reference); got != tt.want {
				t.Errorf("CalculateAge(%v, %v) = %d, want %d", tt.dob, tt.reference, got, tt.want)
			}
		})
	}
}
