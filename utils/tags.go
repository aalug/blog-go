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
