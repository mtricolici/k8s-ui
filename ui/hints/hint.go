package hints

type Hint struct {
	text      string
	shortcut  string
	resources []string
}

func newHint(text, shortcut string, resources ...string) Hint {
	return Hint{
		text:      text,
		shortcut:  shortcut,
		resources: resources,
	}
}

func (h *Hint) VisibleToRes(resource string) bool {
	for _, res := range h.resources {
		if res == "all" || res == resource {
			return true
		}
	}
	return false
}
