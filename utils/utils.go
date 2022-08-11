package utils

import (
	"log"
	"strconv"
)

// ParseFloat parses string to float64.
func ParseFloat(value string) float64 {
	fValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("can't parse float: %s", value)
		return 0.0
	}
	return fValue
}

// FormatFloat parses float64 to string.
func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 10, 64)
}
