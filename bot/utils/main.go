package utils

import (
	"fmt"
	"strconv"
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

func CreateParams(args []string) map[string]interface{} {
	paramsObject := make(map[string]interface{})

	for _, param := range args {
		if strings.Contains(param, ":") {
			splitParam := strings.Split(param, ":")
			key := strings.ToLower(splitParam[0])
			value := splitParam[1]
			if value == "true" || value == "false" {
				boolValue, _ := strconv.ParseBool(value)
				paramsObject[key] = boolValue
			} else {
				paramsObject[key] = value
			}
		} else if strings.HasPrefix(param, "-") {
			key := strings.TrimPrefix(param, "-")
			paramsObject[key] = true
		}
	}

	return paramsObject
}

func CreateHashtags(message string) []string {
	var hashtags []string
	for _, word := range strings.Split(message, " ") {
		if strings.HasPrefix(word, "#") {
		  slicedText := strings.TrimPrefix(word, "#")
			hashtags = append(hashtags, slicedText)
		}
	}
	return hashtags
}