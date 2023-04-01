package utils

import "regexp"

func ShortString(str string, max int) string {
	slen := len(str)
	if slen <= max {
		return str
	}

	return str[:max-3] + "..."
}

func ReplaceSpecialChars(str string) string {
	regex := regexp.MustCompile("[^a-zA-Z0-9_-]+")
	return regex.ReplaceAllString(str, "")
}
