package ui

import (
	gc "github.com/rthornton128/goncurses"
)

type MenuResourceTypes struct {
	screen *gc.Window
	menu   *Menu

	ns             string
	menuStartIndex int
	SelectedType   string
}

func NewResourceTypesMenu(screen *gc.Window, namespace string) *MenuResourceTypes {

	mnu := MenuResourceTypes{
		screen: screen,
		menu:   nil,
		ns:     namespace,
	}
	mnu.load()
	return &mnu
}

func (m *MenuResourceTypes) load() {
	data := [][]string{
		{""},
		{"Pod"},
		{"Service"},
		{"Deployment"},
		{"Ingress"},
		{"pvc"},
		{"DaemonSet"},
		{"ReplicaSet"},
		{"StatefulSet"},
		{"HorizontalPodAutoscaler"},
	}

	m.menu = NewMenu(m.screen, data)
	m.menu.FuncHandleKey = m.HandleKey
	m.menu.FuncCustomStartIndex = m.CustomStartIndex
	m.menuStartIndex = findSelectedIndex(data)

	size_x := 30
	size_y := 16

	x := 20
	y := 1

	m.menu.SetCustomPosition(x, y, size_x, size_y, false)
}

func findSelectedIndex(data [][]string) int {
	for i, v := range data {
		if len(v) > 0 && v[0] == ui_resource_type {
			return i
		}
	}
	return 1 // by default show select 1st element
}

func (m *MenuResourceTypes) Show() {
	m.menu.Show()
}

func (m *MenuResourceTypes) HandleKey(key gc.Key, selectedItem *string) bool {
	if key == gc.KEY_RETURN {
		m.menu.CloseMenu = true
		if selectedItem != nil {
			m.SelectedType = (*selectedItem)
		}
		return true
	}
	return false
}

func (m *MenuResourceTypes) CustomStartIndex() int {
	return m.menuStartIndex
}
