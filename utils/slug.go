package utils

import "regexp"

func IsSlug(str string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)
	return regex.MatchString(str)
}
