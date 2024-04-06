package utils

import (
	"fmt"
	"strings"
	"time"
)

type Predicate[T any] func(T) bool

// Filter an array of type T using a predicate function.
func Filter[T any](slice []T, predicate Predicate[T]) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Find an element of type T in an array using a predicate function.
func Find[T any](slice []T, predicate Predicate[T]) (T, bool) {
	for _, v := range slice {
		if predicate(v) {
			return v, true
		}
	}
	var notFound T
	return notFound, false
}

// Humanize converts a time.Duration to a human-readable string.
func Humanize(d time.Duration, limit int) string {
	if d == 0 {
		return "0s"
	}
	var result []string
	var out string
	if d.Hours() > 8760 {
		result = append(result, fmt.Sprintf("%dy", int(d.Hours()/8760)))
		d -= time.Duration(int(d.Hours()/8760)) * time.Hour * 8760
	}
	if d.Hours() > 720 {
		result = append(result, fmt.Sprintf("%dmo", int(d.Hours()/720)))
		d -= time.Duration(int(d.Hours()/720)) * time.Hour * 720
	}
	if d.Hours() > 24 {
		result = append(result, fmt.Sprintf("%dd", int(d.Hours()/24)))
		d -= time.Duration(int(d.Hours()/24)) * time.Hour * 24
	}
	if d.Hours() > 0 {
		result = append(result, fmt.Sprintf("%dh", int(d.Hours())))
	}
	if d.Minutes() > 0 {
		result = append(result, fmt.Sprintf("%dm", int(d.Minutes())%60))
	}
	if d.Seconds() > 0 {
		result = append(result, fmt.Sprintf("%ds", int(d.Seconds())%60))
	}

	for _, v := range result[:limit] {
		if strings.HasPrefix(v, "0") {
			continue
		}
		out += v + " "
	}
	
	return strings.Trim(out, " ")
}