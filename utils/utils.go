package utils

import (
	"fmt"
	"regexp"
	"time"
)

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

func HumanElapsedTime(t time.Time) string {
	now := time.Now()
	elapsed := now.Sub(t)

	if elapsed < 5*time.Minute {
		return fmt.Sprintf("%ds", int(elapsed.Seconds()))
	} else if elapsed < 5*time.Hour {
		return fmt.Sprintf("%dm%02ds", int(elapsed.Minutes()), int(elapsed.Seconds())%60)
	} else if elapsed < 48*time.Hour {
		return fmt.Sprintf("%dh%02dm", int(elapsed.Hours()), int(elapsed.Minutes())%60)
	}

	days := int(elapsed.Hours() / 24)
	if days <= 60 {
		return fmt.Sprintf("%dd", days)
	}

	months := days / 30
	if months <= 12 {
		return fmt.Sprintf("%dmo", months)
	}

	years := months / 12
	return fmt.Sprintf("%dy", years)
}
