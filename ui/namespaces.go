package ui

import (
	"k8s_ui/k8s"
	"k8s_ui/ncurses"

	gc "github.com/rthornton128/goncurses"
)

type MenuNamespaces struct {
	screen *gc.Window
	k8sc   *k8s.K8SClient
	menu   *Menu

	namespacesCount int
}

func NewMenuNamespaces(screen *gc.Window) *MenuNamespaces {
	mnu := MenuNamespaces{
		screen: screen,
		k8sc:   k8s.NewK8SClient(),
		menu:   nil,
	}
	return &mnu
}

func (m *MenuNamespaces) Load() {
	//TODO: handle errors here
	namespaces, _ := m.k8sc.GetNamespaces()

	m.namespacesCount = len(namespaces) - 1 // 1st is header
	m.menu = NewMenu(m.screen, namespaces)
	m.menu.FuncHeader = m.DrawHeader
	m.menu.FuncHandleKey = m.HandleKey
}

func (m *MenuNamespaces) Show() {
	m.menu.Show()
}

func (m *MenuNamespaces) DrawHeader() {
	_, maxy := m.screen.MaxYX()
	m.screen.ColorOn(ncurses.COLOR_HEADER)
	m.screen.HLine(0, 0, gc.ACS_HLINE, maxy)
	m.screen.MovePrintf(0, 3, " Namespaces: %d of %d ", m.menu.Index, m.namespacesCount)
	m.screen.ColorOff(ncurses.COLOR_HEADER)
}

func (m *MenuNamespaces) HandleKey(key gc.Key, selectedItem []string) bool {

	if key == gc.KEY_RETURN {
		ns := selectedItem[0] // Column 0 is ns name (we don't need other columns here)
		podsMenu := NewMenuPods(m.screen, ns)
		err := podsMenu.Load()
		if err != nil {
			m.menu.ShowError(err)
		} else {
			podsMenu.Show()
		}

		return true
	}
	return false
}
