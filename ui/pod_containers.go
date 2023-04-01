package ui

import (
	gc "github.com/rthornton128/goncurses"
)

type MenuChooseContainer struct {
	screen *gc.Window
	menu   *Menu

	SelectedContainer string
}

func NewMenuChooseContainer(screen *gc.Window, title string, containers []string) *MenuChooseContainer {

	mnu := MenuChooseContainer{
		screen: screen,
		menu:   nil,
	}

	mnu.load(title, containers)
	return &mnu
}

func (m *MenuChooseContainer) load(title string, containers []string) {
	data := [][]string{
		{title},
	}

	for _, container := range containers {
		data = append(data, []string{container})
	}

	m.menu = NewMenu(m.screen, data)
	m.menu.FuncHandleKey = m.HandleKey

	max_y, max_x := m.screen.MaxYX()

	size_x := 30
	size_y := 20

	x := max_x/2 - size_x/2
	y := max_y/2 - size_y/2

	m.menu.SetCustomPosition(x, y, size_x, size_y, true)
}

func (m *MenuChooseContainer) Show() {
	m.menu.Show()
}

func (m *MenuChooseContainer) HandleKey(key gc.Key, selectedItem *string) bool {
	if key == gc.KEY_RETURN {
		m.menu.CloseMenu = true
		if selectedItem != nil {
			m.SelectedContainer = (*selectedItem)
		}
		return true
	}
	return false
}
