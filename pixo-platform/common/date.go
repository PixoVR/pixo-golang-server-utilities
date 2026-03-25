package common

import "time"

// SetDefaultDateRange sets default start and end dates if they are not provided.
// If endDate is nil, it sets endDate to the current date.
// If startDate is nil, it sets startDate to 30 days before the endDate.
// Parameters:
// - startDate: A pointer to the start date, which can be nil.
// - endDate: A pointer to the end date, which can be nil.
// Returns:
// - time.Time: The resolved start date.
// - time.Time: The resolved end date.
func SetDefaultDateRange(startDate, endDate *time.Time) (time.Time, time.Time) {
	now := time.Now()

	if endDate == nil {
		endDate = &now
	}

	// If startDate is null, set it to 30 days before the endDate
	if startDate == nil {
		defaultStart := endDate.AddDate(0, 0, -30)
		startDate = &defaultStart
	}

	return *startDate, *endDate
}

// GetDateThirteenMonthsFromNextMonth calculates the date that is exactly
// thirteen months from the first day of the next month based on the given
// startDate.
//
// Parameters:
//   - startDate: The initial date (time.Time) from which the calculation starts.
//
// Returns:
//   - time.Time: The date thirteen months from the first day of the month
//     immediately following startDate.
//
// The function first determines the first day of the month after startDate,
// then adds thirteen months to this date before returning it. The time
// components (hour, minute, second, nanosecond) of the returned date are set to
// zero, and the timezone/location matches that of the input date.
func GetDateThirteenMonthsFromNextMonth(startDate time.Time) time.Time {
	firstDayOfNextMonth := time.Date(startDate.Year(), startDate.Month()+1, 1, 0, 0, 0, 0, startDate.Location())
	return firstDayOfNextMonth.AddDate(1, 1, 0)
}

func IsValidDateRange(startDate, endDate *time.Time) bool {
	if startDate == nil || endDate == nil {
		return false
	}

	if startDate.IsZero() || endDate.IsZero() {
		return false
	}

	if !startDate.Before(*endDate) {
		return false
	}

	return true
}
