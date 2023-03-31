package ui

func columnRightAlign(column string) bool {
	switch column {
	case "RESTARTS":
		return true
	case "AGE":
		return true
	}
	return false
}
