package ui

import (
	"fmt"
	"k8s_ui/k8s"
	"k8s_ui/ncurses"

	gc "github.com/rthornton128/goncurses"
)

type MenuNamespaces struct {
	screen *gc.Window
	menu   *Menu

	namespacesCount int
}

func NewMenuNamespaces(screen *gc.Window) *MenuNamespaces {
	mnu := MenuNamespaces{
		screen: screen,
		menu:   nil,
	}
	return &mnu
}

func (m *MenuNamespaces) Load() error {
	namespaces, err := k8s.GetNamespaces()
	if err != nil {
		return err
	}

	m.namespacesCount = len(namespaces) - 1 // 1st is header
	if m.menu == nil {
		m.menu = NewMenu(m.screen, namespaces)
		m.menu.FuncHeader = m.DrawHeader
		m.menu.FuncHandleKey = m.HandleKey
		m.menu.Hints1 = [][]string{
			{"Exit", "ESC"},
			{"View yaml", "F3"},
			{"Edit", "F4"},
			{"Refresh", "F5"},
			{"Filter", "F7"},
		}
		m.menu.Hints2 = [][]string{
			{"Describe", "d"},
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

func (m *MenuNamespaces) HandleKey(key gc.Key, selectedItem *string) bool {
	switch key {
	case gc.KEY_RETURN, gc.KEY_ENTER:
		if selectedItem != nil {
			ns := (*selectedItem) // Column 0 is ns name (we don't need other columns here)
			podsMenu := NewResourcesMenu(m.screen, ns)
			err := podsMenu.Load()
			if err != nil {
				ncurses.MessageBox("error", err.Error(), 1000)
			} else {
				podsMenu.Show()
			}
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
			ns := (*selectedItem)
			cmd := fmt.Sprintf("kubectl describe ns %s | less -S", ns)
			ncurses.ExecuteCommand(cmd)
		}
		return true
	case gc.KEY_F3:
		if selectedItem != nil {
			ns := (*selectedItem)
			cmd := fmt.Sprintf("kubectl get ns %s -o yaml|less -S", ns)
			ncurses.ExecuteCommand(cmd)
		}
		return true
	case gc.KEY_F4:
		if selectedItem != nil {
			ns := (*selectedItem)
			cmd := fmt.Sprintf("kubectl edit ns %s", ns)
			ncurses.ExecuteCommand(cmd)
		}
		return true
	case gc.KEY_F7:
		ok, filter := ncurses.InputDialog("Filer (regex allowed)", 30)
		if ok {
			ncurses.MessageBox("warning", fmt.Sprintf("Not implemented yet. filter: '%s'", filter), 1000)
		}

		return true

	}

	return false
}
