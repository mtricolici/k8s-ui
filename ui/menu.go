package ui

import (
	"fmt"

	l "k8s_ui/logger"
	"k8s_ui/ncurses"

	gc "github.com/rthornton128/goncurses"
)

const (
	top_left_x = 0 // Vertical Column
	top_left_y = 2 // Horizontal Line
)

type (
	MenuHeaderFunc    func()
	MenuHandleKeyFunc func(key gc.Key, selectedItem *[]string) bool
)

type Menu struct {
	screen *gc.Window
	data   [][]string
	items  []string

	Hints [][]string

	Index         int
	FuncHeader    MenuHeaderFunc
	FuncHandleKey MenuHandleKeyFunc

	// navigation variables
	max_x           int
	max_y           int
	draw_index_from int
	draw_index_to   int
}

func NewMenu(screen *gc.Window, data [][]string) *Menu {
	menu := Menu{
		screen:        screen,
		data:          data,
		items:         nil,
		Hints:         nil,
		FuncHeader:    nil,
		FuncHandleKey: nil,
	}
	menu.buildItems()
	return &menu
}

func (m *Menu) buildItems() {
	defer l.LogExecutedTime("Menu.buildItems")()

	m.items = make([]string, len(m.data))
	max := make([]int, len(m.data[0]))

	for _, line := range m.data {
		for col, colValue := range line {
			if max[col] < len(colValue) {
				max[col] = len(colValue)
			}
		}
	}

	for i, line := range m.data {
		m.items[i] = ""

		for col, colValue := range line {
			if columnRightAlign(m.data[0][col]) {
				format := fmt.Sprintf("%s%d%s ", "%", max[col], "s")
				m.items[i] += fmt.Sprintf(format, colValue)
			} else {
				format := fmt.Sprintf("%s-%d%s ", "%", max[col], "s")
				m.items[i] += fmt.Sprintf(format, colValue)
			}

		}
	}

	if len(m.items) == 1 {
		m.items[0] = " No resources found"
	}
}

func (m *Menu) Reload(newData [][]string) {
	m.data = newData
	m.buildItems()
}

func (m *Menu) handleKey(key gc.Key) bool {
	if m.Index > 0 && m.Index < len(m.data) {
		return m.FuncHandleKey(key, &m.data[m.Index])
	}

	return m.FuncHandleKey(key, nil)
}

func (m *Menu) calcNavigationVars() {
	m.max_y, m.max_x = m.screen.MaxYX()

	m.draw_index_from = 1
	// 3 lines are: header, shortcuts and table title
	m.draw_index_to = m.max_y - top_left_y - 3

	if m.Index < m.draw_index_from {
		m.Index = m.draw_index_from
	}

	if m.Index > m.draw_index_to {
		m.Index = m.draw_index_to
	}
}

func (m *Menu) navigateUp(jump int) {
	if m.Index-jump < 1 {
		m.Index = 1
		if m.Index < m.draw_index_from {
			delta := m.draw_index_from - m.Index
			m.draw_index_from -= delta
			m.draw_index_to -= delta
		}
	} else {
		m.Index -= jump
		if m.Index < m.draw_index_from {
			delta := m.draw_index_from - m.Index
			m.draw_index_from -= delta
			m.draw_index_to -= delta
		}
	}
}

func (m *Menu) navigateDown(jump int) {
	items_count := len(m.items)

	if m.Index+jump >= items_count {
		m.Index = items_count - 1
		if m.Index > m.draw_index_to {
			delta := m.Index - m.draw_index_to
			m.draw_index_from += delta
			m.draw_index_to += delta
		}

	} else {
		m.Index += jump
		if m.Index > m.draw_index_to {
			delta := m.Index - m.draw_index_to
			m.draw_index_from += delta
			m.draw_index_to += delta
		}
	}
}

func (m *Menu) Show() {
	m.Index = 1
	m.calcNavigationVars()

	for {
		m.screen.Erase()
		m.FuncHeader() // Draw custom header
		m.drawHints()  // Draw shortcut hints
		m.drawMenu()

		m.screen.Refresh()
		key := m.screen.GetChar()

		if !m.handleKey(key) {
			switch key {
			case gc.KEY_DOWN:
				m.navigateDown(1)
			case gc.KEY_UP:
				m.navigateUp(1)
			case gc.KEY_PAGEUP:
				m.navigateUp(m.draw_index_to - m.draw_index_from + 1)
			case gc.KEY_PAGEDOWN:
				m.navigateDown(m.draw_index_to - m.draw_index_from + 1)
			case gc.KEY_ESC:
				return // close menu
			case gc.KEY_BACKSPACE:
				return // close menu
			default:
				ncurses.MessageBox("warning", "key not bound!", 300)
			}
		}
	}
}

func (m *Menu) drawMenu() {
	defer l.LogExecutedTime("drawMenu")()
	if len(m.items) == 1 {
		m.screen.MovePrint(top_left_y, top_left_x, m.items[0])
		return
	}

	// 2 - 1 left border and 1 right border
	windowHorizontalSize := m.max_x - 2 - top_left_x

	x := top_left_x
	y := top_left_y

	for i, item := range m.items {

		// Always print 1st row: this is TITLE
		if i == 0 {
			//m.screen.MovePrint(y, x+1, item)
			ncurses.HLine(ncurses.COLOR_MENU_HEADER, y, x, ' ', windowHorizontalSize+2)
			ncurses.AddText(ncurses.COLOR_MENU_HEADER, y, x+1, item)
			y++ // Move to next line
			m.drawVerticalLineTop(y, x, windowHorizontalSize)
			y++ // Move to next line
			continue
		}

		if i < m.draw_index_from || i > m.draw_index_to {
			continue // Ingore hidden items
		}

		ncurses.AddChar(ncurses.COLOR_MENU_ITEM, y, x, gc.ACS_VLINE)

		if i == m.Index {
			ncurses.HLine(ncurses.COLOR_MENU_ITEM_SELECTED, y, x+1, ' ', windowHorizontalSize)
			ncurses.AddText(ncurses.COLOR_MENU_ITEM_SELECTED, y, x+1, item)
		} else {
			ncurses.HLine(ncurses.COLOR_MENU_ITEM, y, x+1, ' ', windowHorizontalSize)
			ncurses.AddText(ncurses.COLOR_MENU_ITEM, y, x+1, item)
		}

		ncurses.AddChar(ncurses.COLOR_MENU_ITEM, y, x+windowHorizontalSize+1, gc.ACS_VLINE)
		y++ // Move to next line
	}

	m.drawVerticalLineBottom(y, x, windowHorizontalSize)
}

func (m *Menu) drawVerticalLineTop(y int, x int, count int) {
	m.screen.ColorOn(ncurses.COLOR_MENU_ITEM)
	m.screen.MoveAddChar(y, x, gc.ACS_ULCORNER)
	m.screen.HLine(y, x+1, gc.ACS_HLINE, count)
	m.screen.MoveAddChar(y, x+count+1, gc.ACS_URCORNER)
	m.screen.ColorOff(ncurses.COLOR_MENU_ITEM)
}

func (m *Menu) drawVerticalLineBottom(y int, x int, count int) {

	m.screen.ColorOn(ncurses.COLOR_MENU_ITEM)
	m.screen.MoveAddChar(y, x, gc.ACS_LLCORNER)
	m.screen.HLine(y, x+1, gc.ACS_HLINE, count)
	m.screen.MoveAddChar(y, x+count+1, gc.ACS_LRCORNER)
	m.screen.ColorOff(ncurses.COLOR_MENU_ITEM)
}

func (m *Menu) drawHints() {
	if m.Hints == nil || len(m.Hints) < 1 {
		return // no hints for this menu
	}

	x := 1
	y := 1

	ncurses.HLine(ncurses.COLOR_HINTS_TEXT, y, 0, ' ', m.max_x)

	for _, hint := range m.Hints {
		ncurses.AddText(ncurses.COLOR_HINTS_TEXT, y, x, hint[0])
		x += len(hint[0]) + 1
		ncurses.AddText(ncurses.COLOR_HINTS_SHORTCUT, y, x, hint[1])
		x += len(hint[1]) + 2
	}
}
