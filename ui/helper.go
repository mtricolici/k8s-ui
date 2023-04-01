package ui

import (
	"k8s_ui/ncurses"
)

func columnRightAlign(column string) bool {
	switch column {
	case "RESTARTS":
		return true
	case "AGE":
		return true
	}
	return false
}

func (m *MenuResources) chooseContainer(title string, podName string) string {
	containers, err := m.k8sc.GetPodContainerNames(m.ns, podName)
	if err != nil {
		ncurses.MessageBox("Error", err.Error(), 1000)
		return ""
	} else {
		if len(containers) == 1 {
			return containers[0] // we have just 1 container
		}
		mnu := NewMenuChooseContainer(m.screen, title, containers)
		mnu.Show()
		return mnu.SelectedContainer
	}
}
