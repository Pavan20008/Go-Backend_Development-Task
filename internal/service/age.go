package service

import "time"

// CalculateAge returns the age in completed years for a given date of birth
// relative to the supplied reference time. The reference time is passed in
// (rather than read from time.Now internally) so the calculation is pure and
// easy to unit test.
//
// A person's age only increments once their birthday (month and day) has been
// reached in the reference year.
func CalculateAge(dob, reference time.Time) int {
	// Normalise both to date-only in UTC to avoid timezone/time-of-day drift.
	dob = time.Date(dob.Year(), dob.Month(), dob.Day(), 0, 0, 0, 0, time.UTC)
	reference = time.Date(reference.Year(), reference.Month(), reference.Day(), 0, 0, 0, 0, time.UTC)

	if reference.Before(dob) {
		return 0
	}

	age := reference.Year() - dob.Year()

	// If the birthday hasn't occurred yet this year, subtract one.
	beforeBirthday := reference.Month() < dob.Month() ||
		(reference.Month() == dob.Month() && reference.Day() < dob.Day())
	if beforeBirthday {
		age--
	}

	return age
}
