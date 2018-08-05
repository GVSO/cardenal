package timeutils

// GetSeconds return the amount of seconds in n days.
func GetSeconds(days int) int {
	return days * 24 * 60 * 60
}
