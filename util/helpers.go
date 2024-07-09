package util

import (
	"errors"
	"time"
)

func ParseDate(dateString string) (time.Time, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, should be YYYY-MM-DD")
	}
	return date, nil
}

func Difference(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	for _, item := range slice2 {
		m[item] = true
	}
	var diff []string
	for _, item := range slice1 {
		if !m[item] {
			diff = append(diff, item)
		}
	}
	return diff
}
