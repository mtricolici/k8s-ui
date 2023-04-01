package ui

import (
	"fmt"
	"k8s_ui/k8s"
	"k8s_ui/ncurses"
	"strings"

	gc "github.com/rthornton128/goncurses"
)

type MenuResources struct {
	screen *gc.Window
	k8sc   *k8s.K8SClient
	menu   *Menu

	ns         string
	itemsCount int
	wide       bool

	resourceType string
}

func NewResourcesMenu(screen *gc.Window, namespace string) *MenuResources {
	mnu := MenuResources{
		screen:       screen,
		k8sc:         k8s.NewK8SClient(),
		menu:         nil,
		ns:           namespace,
		wide:         false,
		resourceType: "pod", //default show pods in a namespace
	}

	return &mnu
}

func (m *MenuResources) Load() error {
	//TODO: load resources according to m.resourceType
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

	ncurses.HLine(ncurses.COLOR_HEADER, 0, 0, ' ', maxy)
	ncurses.AddText(ncurses.COLOR_HEADER, 0, 1, "Namespace:")
	ncurses.AddText(ncurses.COLOR_HEADER_HIGH, 0, 11, m.ns)
	x := 12 + len(m.ns)

	ncurses.AddText(ncurses.COLOR_HEADER, 0, x, "view:")
	x += 5
	ncurses.AddText(ncurses.COLOR_HEADER_HIGH, 0, x, m.resourceType)
	x += len(m.resourceType)
	ncurses.AddText(ncurses.COLOR_HEADER_HINT, 0, x, "<F2>")
	x += 5

	ncurses.AddText(ncurses.COLOR_HEADER, 0, x, fmt.Sprintf("%d of %d", m.menu.Index, m.itemsCount))
}

func (m *MenuResources) HandleKey(key gc.Key, selectedItem *[]string) bool {
	switch key {
	case 111: // key 'o' has been pressed
		m.wide = !m.wide
		m.reload()
		return true
	case gc.KEY_F5:
		m.reload()
		return true
	case gc.KEY_F2:
		mnu := NewResourceTypesMenu(m.screen, m.ns)
		mnu.Show()

		if len(mnu.SelectedType) > 0 {
			if mnu.SelectedType == "all" {
				ncurses.MessageBox("Error", "ALL not implemented yet", 1000)
			} else if strings.HasPrefix(mnu.SelectedType, "custom") {
				ncurses.MessageBox("Error", "CUSTOM not implemented yet", 1000)
			} else {
				m.resourceType = mnu.SelectedType
				m.reload()
			}
		}

		return true
	}

	return false
}

func (m *MenuResources) reload() {
	win := ncurses.MessageBoxAsync("", "Loading ...")
	err := m.Load()
	win.Delete() // close 'Loading' dialog ...
	if err != nil {
		ncurses.MessageBox("Error", err.Error(), 1000)
	}
	win.Delete() // close 'Reloading' dialog ...
}
