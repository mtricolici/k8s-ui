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

func (m *MenuResources) chooseContainer(podName string) string {
	containers, err := m.k8sc.GetPodContainerNames(m.ns, podName)
	if err != nil {
		ncurses.MessageBox("Error", err.Error(), 1000)
		return ""
	} else {
		//TODO: show a menu that allows you to choose container
		return containers[0]
	}
}
