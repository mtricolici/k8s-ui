package ui

import (
	"fmt"
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

func (m *MenuNamespaces) Load() error {
	namespaces, err := m.k8sc.GetNamespaces()
	if err != nil {
		return err
	}

	m.namespacesCount = len(namespaces) - 1 // 1st is header
	if m.menu == nil {
		m.menu = NewMenu(m.screen, namespaces)
		m.menu.FuncHeader = m.DrawHeader
		m.menu.FuncHandleKey = m.HandleKey
		m.menu.Hints = [][]string{
			{"Exit", "ESC"},
			{"Describe", "d"},
			{"Edit", "F4"},
			{"Filter", "F3"},
			{"Refresh", "F5"},
		}
	} else {
		m.menu.Reload(namespaces)
	}

	return nil
}

func (m *MenuNamespaces) Show() {
	m.menu.Show()
}

func (m *MenuNamespaces) DrawHeader() {
	_, maxy := m.screen.MaxYX()
	m.screen.ColorOn(ncurses.COLOR_HEADER)
	m.screen.HLine(0, 0, ' ', maxy)
	m.screen.MovePrintf(0, 1, "Namespaces: %d of %d ", m.menu.Index, m.namespacesCount)
	m.screen.ColorOff(ncurses.COLOR_HEADER)
}

func (m *MenuNamespaces) HandleKey(key gc.Key, selectedItem *[]string) bool {

	switch key {
	case gc.KEY_RETURN:
		ns := (*selectedItem)[0] // Column 0 is ns name (we don't need other columns here)
		podsMenu := NewResourcesMenu(m.screen, ns)
		err := podsMenu.Load()
		if err != nil {
			ncurses.MessageBox("error", err.Error(), 1000)
		} else {
			podsMenu.Show()
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
	case 100: // character 'd'
		if selectedItem != nil {
			ns := (*selectedItem)[0]
			cmd := fmt.Sprintf("kubectl describe ns %s | less -S", ns)
			ncurses.ExecuteCommand(cmd)
		}
		return true
	case gc.KEY_F4:
		if selectedItem != nil {
			ns := (*selectedItem)[0]
			cmd := fmt.Sprintf("kubectl edit ns %s", ns)
			ncurses.ExecuteCommand(cmd)
		}
		return true

	}

	return false
}
