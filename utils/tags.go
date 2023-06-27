package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// IsTagList check if the string is a tag list
func IsTagList(str string) bool {
	pattern := "^\\d+(,\\d+)*,?$"
	match, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return match
}

// TagsToIntSlice convert tag list string to int32 slice
func TagsToIntSlice(str string) ([]int32, error) {
	numbers := strings.Split(str, ",")
	intSlice := make([]int32, 0, len(numbers))

	for _, num := range numbers {
		if num != "" {
			n, err := strconv.ParseInt(num, 10, 32)
			if err != nil {
				return nil, err
			}
			intSlice = append(intSlice, int32(n))
		}
	}

	return intSlice, nil
}

// CompareTagLists compare two tag lists and return unique tags in each list
func CompareTagLists(slice1 []string, slice2 []string) ([]string, []string) {
	uniqueSlice1 := make([]string, 0)
	uniqueSlice2 := make([]string, 0)

	// Create a map to store the presence of each string in slice2
	lookup := make(map[string]bool)
	for _, s := range slice2 {
		lookup[s] = true
	}

	// Check each string in slice1
	for _, s := range slice1 {
		// If the string is not found in slice2, add it to uniqueSlice1
		if !lookup[s] {
			uniqueSlice1 = append(uniqueSlice1, s)
		}
	}

	// Reset the lookup map to store the presence of each string in slice1
	lookup = make(map[string]bool)
	for _, s := range slice1 {
		lookup[s] = true
	}

	// Check each string in slice2
	for _, s := range slice2 {
		// If the string is not found in slice1, add it to uniqueSlice2
		if !lookup[s] {
			uniqueSlice2 = append(uniqueSlice2, s)
		}
	}

	return uniqueSlice1, uniqueSlice2
}

// contains helper function to check if a string is present in a slice
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
