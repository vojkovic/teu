package common

import (
	"fmt"
	"strconv"
	"time"
)

// ConvertUnixToHumanReadable converts Unix timestamp to human-readable format
func ConvertUnixToHumanReadable(unixTimeStr string) string {
	unixTime, err := strconv.ParseInt(unixTimeStr, 10, 64)
	if err != nil {
		return "Invalid Unix timestamp"
	}

	// Convert Unix timestamp to time object
	timestamp := time.Unix(unixTime, 0)

	// Calculate time difference
	diff := time.Since(timestamp)

	// Convert time difference to human-readable format
	switch {
	case diff.Seconds() < 60:
		if int(diff.Seconds()) == 1 {
			return "1 second ago"
		}
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	case diff.Minutes() < 60:
		if int(diff.Minutes()) == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff.Hours() < 24:
		if int(diff.Hours()) == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	default:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}