package ui

import (
	"k8s_ui/k8s"
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

func resourceHasLogs() bool {
	switch ui_resource_type {
	case "Pod":
		return true
	case "Deployment":
		return true
	case "ReplicaSet":
		return true
	case "StatefulSet":
		return true
	case "DaemonSet":
		return true
	default:
		return false
	}
}

func (m *MenuResources) chooseContainer(title string, podName string) string {
	containers, err := k8s.GetPodContainerNames(m.ns, podName)
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
