package ncurses

func shortString(str string, max int) string {
	slen := len(str)
	if slen <= max {
		return str
	}

	return str[:max-3] + "..."
}
