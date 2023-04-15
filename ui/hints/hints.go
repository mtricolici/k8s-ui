package hints

var (
	instance *Hints
)

type Hints struct {
	hints []Hint
}

func createHints() *Hints {
	hints := []Hint{
		newHint("Describe", "d", "all"),
		newHint("Wide", "o", "all"),
		newHint("Logs", "l", "Pod", "Deployment", "ReplicaSet", "StatefulSet", "DaemonSet"),
		newHint("Prev.logs", "p", "Pod", "Deployment", "ReplicaSet", "StatefulSet", "DaemonSet"),
		newHint("Exec", "e", "Pod"),
		newHint("Port-Forward", "p", "Service"),
		newHint("Restart", "r", "Deployment"),
	}

	return &Hints{
		hints: hints,
	}
}

func (h *Hints) GetForResource(resource string) [][]string {
	//TODO: implement caching. Store in a map?!
	hints := [][]string{}

	for _, hint := range h.hints {
		if hint.VisibleToRes(resource) {
			hints = append(hints, []string{hint.text, hint.shortcut})
		}
	}

	return hints
}

func GetHints() *Hints {
	if instance == nil {
		instance = createHints()
	}
	return instance
}
