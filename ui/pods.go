package ui

import (
	"k8s_ui/k8s"
	"k8s_ui/ncurses"

	gc "github.com/rthornton128/goncurses"
)

type MenuPods struct {
	screen *gc.Window
	k8sc   *k8s.K8SClient
	menu   *Menu

	ns        string
	podsCount int
}

func NewMenuPods(screen *gc.Window, namespace string) *MenuPods {
	mnu := MenuPods{
		screen: screen,
		k8sc:   k8s.NewK8SClient(),
		menu:   nil,
		ns:     namespace,
	}
	return &mnu
}

func (m *MenuPods) Load() error {
	pods, err := m.k8sc.GetPods(m.ns)
	if err != nil {
		return err
	}

	m.podsCount = len(pods)
	m.menu = NewMenu(m.screen, pods)
	m.menu.FuncHeader = m.DrawHeader
	m.menu.FuncHandleKey = m.HandleKey
	return nil
}

func (m *MenuPods) Show() {
	m.menu.Show()
}

func (m *MenuPods) DrawHeader() {
	_, maxy := m.screen.MaxYX()
	m.screen.ColorOn(ncurses.COLOR_HEADER)
	m.screen.HLine(0, 0, gc.ACS_HLINE, maxy)
	m.screen.MovePrintf(0, 3, " Namespace '%s' pods: %d of %d ", m.ns, m.menu.Index, m.podsCount)
	m.screen.ColorOff(ncurses.COLOR_HEADER)
}

func (m *MenuPods) HandleKey(key gc.Key, selectedItem []string) bool {
	return false
}
