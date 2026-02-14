package validators

import (
	"errors"
	"strings"
)

// Allowed turf status values
var allowedTurfStatuses = map[string]bool{
	"active": true, "inactive": true,
}

const (
	minHour       = 0
	maxHour       = 23
	maxNoOfFields = 50
)

// ValidateTurfInput validates turf name, start/end time, status, and noOfFields.
// startTime and endTime are hours of day (0-23).
func ValidateTurfInput(name string, startTime, endTime int, status string, noOfFields int, address string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}

	if startTime < minHour || startTime > maxHour {
		return errors.New("startTime must be between 0 and 23 (hour of day)")
	}
	if endTime < minHour || endTime > maxHour {
		return errors.New("endTime must be between 0 and 23 (hour of day)")
	}
	if endTime <= startTime {
		return errors.New("endTime must be after startTime")
	}

	if status != "" && !allowedTurfStatuses[strings.ToLower(status)] {
		return errors.New("status must be one of: active, inactive")
	}

	if noOfFields < 1 {
		return errors.New("noOfFields must be at least 1")
	}
	if noOfFields > maxNoOfFields {
		return errors.New("noOfFields exceeds maximum allowed")
	}

	if strings.TrimSpace(address) == "" {
		return errors.New("address is required")
	}

	return nil
}
