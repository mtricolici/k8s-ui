package ui

import (
	"k8s_ui/k8s"
	"k8s_ui/ncurses"

	gc "github.com/rthornton128/goncurses"
)

type MenuResources struct {
	screen *gc.Window
	k8sc   *k8s.K8SClient
	menu   *Menu

	ns         string
	itemsCount int
	wide       bool
}

func NewResourcesMenu(screen *gc.Window, namespace string) *MenuResources {
	mnu := MenuResources{
		screen: screen,
		k8sc:   k8s.NewK8SClient(),
		menu:   nil,
		ns:     namespace,
		wide:   false,
	}

	return &mnu
}

func (m *MenuResources) Load() error {
	pods, err := m.k8sc.GetPods(m.ns, m.wide)
	if err != nil {
		return err
	}

	m.itemsCount = len(pods) - 1 // 1st is header
	if m.menu == nil {
		m.menu = NewMenu(m.screen, pods)
		m.menu.FuncHeader = m.DrawHeader
		m.menu.FuncHandleKey = m.HandleKey

		m.menu.Hints = [][]string{
			{"Back", "ESC"},
			{"Wide", "o"},
			{"Logs", "l"},
			{"Prev.logs", "p"},
			{"Describe", "d"},
			{"Exec", "e"},
			{"Filter", "F3"},
			{"Refresh", "F5"},
		}
	} else {
		m.menu.Reload(pods)
	}

	return nil
}

func (m *MenuResources) Show() {
	m.menu.Show()
}

func (m *MenuResources) DrawHeader() {
	_, maxy := m.screen.MaxYX()
	m.screen.ColorOn(ncurses.COLOR_HEADER)
	m.screen.HLine(0, 0, gc.ACS_HLINE, maxy)
	m.screen.MovePrintf(0, 3, " Namespace '%s' resource %d of %d ", m.ns, m.menu.Index, m.itemsCount)
	m.screen.ColorOff(ncurses.COLOR_HEADER)
}

func (m *MenuResources) HandleKey(key gc.Key, selectedItem *[]string) bool {
	switch key {
	case 111: // key 'o' has been pressed
		win := ncurses.MessageBoxAsync("", "Loading ...")
		m.wide = !m.wide
		err := m.Load()
		win.Delete() // close 'Loading' dialog ...
		if err != nil {
			ncurses.MessageBox("Error", err.Error(), 1000)
		}
		return true
	case gc.KEY_F5:
		win := ncurses.MessageBoxAsync("", "reloading ...")
		err := m.Load() // reload
		win.Delete()    // close 'Reloading' dialog ...
		if err != nil {
			ncurses.MessageBox("Error", err.Error(), 1000)
		}
		return true
	}

	return false
}
