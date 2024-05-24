package pcommon

import (
	"regexp"
	"time"
)

func ContainsDigit(s string) bool {
	matched, _ := regexp.MatchString(`[0-9]`, s)
	return matched
}

func ChunkString(s string, n int) []string {
	if n <= 0 {
		return []string{}
	}

	var result []string
	runes := []rune(s) // Convert the string to runes to handle Unicode characters properly

	for i := 0; i < len(runes); i += n {
		end := i + n
		if end > len(runes) {
			end = len(runes)
		}
		result = append(result, string(runes[i:end]))
	}

	return result
}

func Contains[T time.Duration | int64](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func Unique[T time.Duration | string](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
