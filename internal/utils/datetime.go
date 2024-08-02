package utils

import "time"

/*
ValidateDate validates the date string and returns a valid RFC3339 date string.

If the date string is empty, it returns an empty string.
*/
func ValidateDate(date string) string {
	if date == "" {
		return ""
	}

	// convert date to time
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return ""
	}

	// convert time to string
	return t.Format(time.RFC3339)
}
