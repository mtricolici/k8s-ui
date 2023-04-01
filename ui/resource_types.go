package ui

import (
	"k8s_ui/k8s"

	gc "github.com/rthornton128/goncurses"
)

type MenuResourceTypes struct {
	screen *gc.Window
	k8sc   *k8s.K8SClient
	menu   *Menu

	ns string
}

func NewResourceTypesMenu(screen *gc.Window, namespace string) *MenuResourceTypes {

	mnu := MenuResourceTypes{
		screen: screen,
		k8sc:   k8s.NewK8SClient(),
		menu:   nil,
		ns:     namespace,
	}
	mnu.load()
	return &mnu
}

func (m *MenuResourceTypes) load() {
	data := [][]string{
		{""},
		{"all"},
		{"pod"},
		{"service"},
		{"deployment"},
		{"ingress"},
		{"pvc"},
		{"daemonset"},
		{"replicaset"},
		{"statefulset"},
		{"custom ..."},
	}
	m.menu = NewMenu(m.screen, data)
	m.menu.FuncHandleKey = m.HandleKey

	size_x := 16
	size_y := 16

	x := 20
	y := 1

	m.menu.SetCustomPosition(x, y, size_x, size_y, false)
}

func (m *MenuResourceTypes) Show() {
	m.menu.Show()
}

func (m *MenuResourceTypes) HandleKey(key gc.Key, selectedItem *[]string) bool {
	return false
}
