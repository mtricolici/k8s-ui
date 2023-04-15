package ui

import (
	"fmt"
	"k8s_ui/k8s"
	"k8s_ui/ncurses"
	"k8s_ui/ui/hints"
	"k8s_ui/utils"

	gc "github.com/rthornton128/goncurses"
)

type MenuResources struct {
	screen *gc.Window
	menu   *Menu

	ns         string
	itemsCount int
}

func NewResourcesMenu(screen *gc.Window, namespace string) *MenuResources {
	mnu := MenuResources{
		screen: screen,
		menu:   nil,
		ns:     namespace,
	}

	return &mnu
}

func (m *MenuResources) Load() error {
	resources, err := k8s.GetResources(m.ns, ui_resource_type, ui_wide_view)
	if err != nil {
		return err
	}

	m.itemsCount = len(resources) - 1 // 1st is header
	if m.menu == nil {
		m.menu = NewMenu(m.screen, resources)
		m.menu.FuncHeader = m.DrawHeader
		m.menu.FuncHandleKey = m.HandleKey

		m.menu.Hints1 = [][]string{
			{"Back", "ESC"},
			{"View yaml", "F3"},
			{"Edit", "F4"},
			{"Refresh", "F5"},
			{"Filter", "F7"},
		}
	} else {
		m.menu.Reload(resources)
		m.itemsCount = m.menu.GetItemsCount()
	}

	m.menu.Hints2 = hints.GetHints().GetForResource(ui_resource_type)

	return nil
}

func (m *MenuResources) Show() {
	m.menu.Show()
}

func (m *MenuResources) DrawHeader() {
	_, maxy := m.screen.MaxYX()

	ncurses.HLine(ncurses.COLOR_HEADER, 0, 0, ' ', maxy)
	ncurses.AddText(ncurses.COLOR_HEADER, 0, 1, "Namespace:")
	shortName := utils.ShortString(m.ns, max_name_size)
	ncurses.AddText(ncurses.COLOR_HEADER_HIGH, 0, 11, shortName)
	x := 12 + len(shortName)

	ncurses.AddText(ncurses.COLOR_HEADER, 0, x, "view:")
	x += 5
	ncurses.AddText(ncurses.COLOR_HEADER_HIGH, 0, x, ui_resource_type)
	x += len(ui_resource_type)
	ncurses.AddText(ncurses.COLOR_HEADER_HINT, 0, x, "<F2>")
	x += 5

	ncurses.AddText(ncurses.COLOR_HEADER, 0, x, fmt.Sprintf("%d of %d", m.menu.Index, m.itemsCount))
}

func (m *MenuResources) HandleKey(key gc.Key, selectedItem *string) bool {
	switch key {
	case 111: // key 'o' has been pressed
		ui_wide_view = !ui_wide_view
		m.reload()
		return true
	case gc.KEY_F5:
		m.reload()
		return true
	case gc.KEY_F2:
		mnu := NewResourceTypesMenu(m.screen, m.ns)
		mnu.Show()
		if len(mnu.SelectedType) > 0 {
			ui_resource_type = mnu.SelectedType
			m.reload()
		}
		return true
	case 100: // character 'd' - describe resource
		if selectedItem != nil {
			ncurses.ExecuteCommand(true, "kubectl describe %s %s -n %s", ui_resource_type, *selectedItem, m.ns)
		}
		return true
	case 108: // character 'l' - view logs
		m.showLogs(selectedItem, "")
		return true
	case 112: // character 'p' - view previous logs
		m.showLogs(selectedItem, "-p")
		return true
	case 101: // character 'e' - execute a shell inside container
		m.executeShell(selectedItem)
		return true
	case gc.KEY_F4: // edit resource as yaml
		if selectedItem != nil {
			ncurses.ExecuteCommand(false, "kubectl edit %s %s -n %s", ui_resource_type, *selectedItem, m.ns)
		}
		return true
	case gc.KEY_F3: // View resource yaml
		if selectedItem != nil {
			ncurses.ExecuteCommand(true, "kubectl get %s %s -n %s -o yaml", ui_resource_type, *selectedItem, m.ns)
		}
		return true
	case gc.KEY_F7:
		ok, filter := ncurses.InputDialog("Filer (regex)", 30)
		if ok {
			err := m.menu.SetFilter(filter)
			if err != nil {
				ncurses.MessageBox("error", err.Error(), 1000)
			} else {
				m.reload()
			}
		}

		return true
	case 115: // 's' was pressed
		if selectedItem != nil && ui_resource_type == "helm" {
			ncurses.ExecuteCommand(true, "helm status %s -n %s", *selectedItem, m.ns)
		}
		return true
	case 104: // 'h' was pressed
		if selectedItem != nil && ui_resource_type == "helm" {
			ncurses.ExecuteCommand(true, "helm history %s -n %s", *selectedItem, m.ns)
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
}

func (m *MenuResources) showLogs(resourceName *string, options string) {
	if resourceName != nil && resourceHasLogs() {
		name := *resourceName
		if ui_resource_type == "Pod" {
			container := m.chooseContainer("Logs for ?", name)
			if len(container) > 0 {
				ncurses.ExecuteCommand(true, "kubectl logs %s %s -n %s -c %s", options, name, m.ns, container)
			}
		} else {
			ncurses.ExecuteCommand(true, "kubectl logs %s %s/%s -n %s --all-containers=true", options, ui_resource_type, name, m.ns)
		}
	}
}

func (m *MenuResources) executeShell(resourceName *string) {
	if resourceName != nil && ui_resource_type == "Pod" {
		pod := *resourceName
		container := m.chooseContainer("Execute where ?", pod)
		if len(container) > 0 {
			ncurses.ExecuteCommand(false, "kubectl exec -it %s -n %s -c %s -- sh", pod, m.ns, container)
		}
	}
}
