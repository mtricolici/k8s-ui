package ui

import (
	"fmt"

	"k8s_ui/ncurses"
	"k8s_ui/utils"

	gc "github.com/rthornton128/goncurses"
)

const (
	max_name_size = 35
)

type (
	MenuHeaderFunc       func()
	MenuHandleKeyFunc    func(key gc.Key, selectedItem *string) bool
	MenuCustomStartIndex func() int
)

type Menu struct {
	screen *gc.Window

	names []string
	items []string

	// max lenght per column - or number of characters in a column
	max []int

	Hints1 [][]string
	Hints2 [][]string

	Index                int
	FuncHeader           MenuHeaderFunc
	FuncHandleKey        MenuHandleKeyFunc
	FuncCustomStartIndex MenuCustomStartIndex

	CloseMenu bool

	// menu position variables
	top_left_x   int
	top_left_y   int
	menu_size_x  int
	menu_size_y  int
	show_header  bool
	erase_screen bool

	// navigation variables
	draw_index_from int
	draw_index_to   int
}

func NewMenu(screen *gc.Window, data [][]string) *Menu {
	max_y, max_x := screen.MaxYX()
	menu := Menu{
		CloseMenu:     false,
		screen:        screen,
		items:         nil,
		Hints1:        nil,
		Hints2:        nil,
		FuncHeader:    nil,
		FuncHandleKey: nil,
		// set default to full screen without 2 rows for header
		top_left_x:   0,
		top_left_y:   3,
		menu_size_x:  max_x,
		menu_size_y:  max_y - 3,
		show_header:  true,
		erase_screen: true,
	}
	menu.buildItems(data)
	return &menu
}

func (m *Menu) SetCustomPosition(x, y, size_x, size_y int, show_header bool) {
	m.top_left_x = x
	m.top_left_y = y
	m.menu_size_x = size_x
	m.menu_size_y = size_y
	m.show_header = show_header
	m.erase_screen = false
}

func (m *Menu) buildItems(data [][]string) {
	m.items = make([]string, len(data))
	m.names = make([]string, len(data))
	m.max = make([]int, len(data[0]))

	// make names shorter!
	for i := range data {
		m.names[i] = data[i][0]
		data[i][0] = utils.ShortString(data[i][0], max_name_size)
	}

	for _, line := range data {
		for col, colValue := range line {
			if m.max[col] < len(colValue) {
				m.max[col] = len(colValue)
			}
		}
	}

	for i, line := range data {
		m.items[i] = ""

		for col, colValue := range line {
			if i > 0 && columnRightAlign(data[0][col]) {
				format := fmt.Sprintf("%s%d%s ", "%", m.max[col], "s")
				m.items[i] += fmt.Sprintf(format, colValue)
			} else {
				format := fmt.Sprintf("%s-%d%s ", "%", m.max[col], "s")
				m.items[i] += fmt.Sprintf(format, colValue)
			}
		}
	}

	if len(m.items) == 1 {
		m.items[0] = " No resources found"
	}
}

func (m *Menu) Reload(newData [][]string) {
	m.buildItems(newData)
	m.navigateTo(0)
}

func (m *Menu) handleKey(key gc.Key) bool {
	if m.Index > 0 && m.Index < len(m.names) {
		return m.FuncHandleKey(key, &m.names[m.Index])
	}

	return m.FuncHandleKey(key, nil)
}

func (m *Menu) calcNavigationVars() {

	m.draw_index_from = 1
	// 3 lines are: header, shortcuts and table title

	if m.show_header {
		m.draw_index_to = m.menu_size_y - 3
	} else {
		m.draw_index_to = m.menu_size_y - 2
	}

	if m.Index < m.draw_index_from {
		m.Index = m.draw_index_from
	}

	if m.Index > m.draw_index_to {
		m.Index = m.draw_index_to
	}
}

func (m *Menu) adjustDrawIndexes() {
	if m.Index < m.draw_index_from {
		delta := m.draw_index_from - m.Index
		m.draw_index_from -= delta
		m.draw_index_to -= delta
	} else if m.Index > m.draw_index_to {
		delta := m.Index - m.draw_index_to
		m.draw_index_from += delta
		m.draw_index_to += delta
	}
}

// Jump can be positive (move down) or negative (move up)
func (m *Menu) navigateTo(jump int) {
	items_count := len(m.items)

	if m.Index+jump < 1 {
		m.Index = 1
		m.adjustDrawIndexes()
	} else if m.Index+jump >= items_count {
		m.Index = items_count - 1
		m.adjustDrawIndexes()
	} else {
		m.Index += jump
		m.adjustDrawIndexes()
	}
}

func (m *Menu) Show() {
	if m.FuncCustomStartIndex != nil {
		m.Index = m.FuncCustomStartIndex()
	} else {
		m.Index = 1
	}

	m.calcNavigationVars()

	for {
		if m.erase_screen {
			m.screen.Erase()
		}

		if m.FuncHeader != nil {
			m.FuncHeader() // Draw custom header
		}

		m.drawHints(1, m.Hints1)
		m.drawHints(2, m.Hints2) // Draw shortcut hints
		m.drawMenu()

		m.screen.Refresh()
		key := m.screen.GetChar()

		if !m.handleKey(key) {
			switch key {
			case gc.KEY_DOWN:
				m.navigateTo(1)
			case gc.KEY_UP:
				m.navigateTo(-1)
			case gc.KEY_PAGEUP:
				m.navigateTo(m.draw_index_from - m.draw_index_to - 1)
			case gc.KEY_PAGEDOWN:
				m.navigateTo(m.draw_index_to - m.draw_index_from + 1)
			case gc.KEY_ESC:
				return // close menu
			default:
				ncurses.MessageBox("warning", "key not bound!", 300)
			}
		}

		if m.CloseMenu {
			break
		}
	}
}

func (m *Menu) drawMenu() {
	if len(m.items) == 1 {
		m.screen.MovePrint(m.top_left_y, m.top_left_x, m.items[0])
		return
	}

	// 2 - 1 left border and 1 right border
	windowHorizontalSize := m.menu_size_x - 2

	x := m.top_left_x
	y := m.top_left_y

	for i, item := range m.items {

		// Always print 1st row: this is TITLE
		if i == 0 {
			if m.show_header {
				ncurses.HLine(ncurses.COLOR_MENU_HEADER, y, x, ' ', windowHorizontalSize+2)
				ncurses.AddTextMaxWidth(ncurses.COLOR_MENU_HEADER, y, x+1, windowHorizontalSize, item)
				m.drawColumnsDelimiters(y, ncurses.COLOR_MENU_HEADER_DELIMITER, ncurses.COLOR_MENU_HEADER_DELIMITER)
				y++ // Move to next line
			}
			m.drawVerticalLineTop(y, x, windowHorizontalSize)
			y++ // Move to next line
			continue
		}

		if i < m.draw_index_from || i > m.draw_index_to {
			continue // Ingore hidden items
		}

		//ncurses.AddChar(ncurses.COLOR_MENU_ITEM, y, x, gc.ACS_VLINE)

		if i == m.Index {
			ncurses.HLine(ncurses.COLOR_MENU_ITEM_SELECTED, y, x+1, ' ', windowHorizontalSize)
			ncurses.AddTextMaxWidth(ncurses.COLOR_MENU_ITEM_SELECTED, y, x+1, windowHorizontalSize, item)
			m.drawColumnsDelimiters(y, ncurses.COLOR_MENU_ITEM_DELIMITER, ncurses.COLOR_MENU_ITEM_SELECTED)
		} else {
			ncurses.HLine(ncurses.COLOR_MENU_ITEM, y, x+1, ' ', windowHorizontalSize)
			ncurses.AddTextMaxWidth(ncurses.COLOR_MENU_ITEM, y, x+1, windowHorizontalSize, item)
			m.drawColumnsDelimiters(y, ncurses.COLOR_MENU_ITEM_DELIMITER, ncurses.COLOR_MENU_ITEM_DELIMITER)
		}

		//ncurses.AddChar(ncurses.COLOR_MENU_ITEM, y, x+windowHorizontalSize+1, gc.ACS_VLINE)
		y++ // Move to next line
	}

	m.drawVerticalLineBottom(y, x, windowHorizontalSize)
}

func (m *Menu) drawColumnsDelimiters(y int, color1, color2 int16) {
	x := m.top_left_x

	max_x := x + m.menu_size_x - 1

	for i, max := range m.max {
		if x < max_x {
			if i == 0 {
				ncurses.AddChar(color1, y, x, gc.ACS_VLINE)
			} else {
				ncurses.AddChar(color2, y, x, gc.ACS_VLINE)
			}
		}

		x += max + 1
	}
	ncurses.AddChar(color1, y, max_x, gc.ACS_VLINE)
}

func (m *Menu) drawVerticalLineTop(y int, x int, count int) {
	m.screen.ColorOn(ncurses.COLOR_MENU_ITEM_DELIMITER)
	m.screen.MoveAddChar(y, x, gc.ACS_LTEE)
	m.screen.HLine(y, x+1, gc.ACS_HLINE, count)
	m.screen.MoveAddChar(y, x+count+1, gc.ACS_RTEE)

	max_x := x + count + 1
	columns_count := len(m.max)
	for i, max := range m.max {
		x += max + 1
		if x < max_x && i < columns_count-1 {
			m.screen.MoveAddChar(y, x, gc.ACS_PLUS)
		}
	}
	m.screen.ColorOff(ncurses.COLOR_MENU_ITEM_DELIMITER)
}

func (m *Menu) drawVerticalLineBottom(y int, x int, count int) {

	m.screen.ColorOn(ncurses.COLOR_MENU_ITEM_DELIMITER)
	m.screen.MoveAddChar(y, x, gc.ACS_LLCORNER)
	m.screen.HLine(y, x+1, gc.ACS_HLINE, count)
	m.screen.MoveAddChar(y, x+count+1, gc.ACS_LRCORNER)
	max_x := x + count + 1
	columns_count := len(m.max)
	for i, max := range m.max {
		x += max + 1
		if x < max_x && i < columns_count-1 {
			m.screen.MoveAddChar(y, x, gc.ACS_BTEE)
		}
	}
	m.screen.ColorOff(ncurses.COLOR_MENU_ITEM_DELIMITER)
}

func (m *Menu) drawHints(y int, hints [][]string) {
	if hints == nil || len(hints) < 1 {
		return // no hints for this menu
	}

	x := 1

	ncurses.HLine(ncurses.COLOR_HINTS_TEXT, y, m.top_left_x, ' ', m.menu_size_x)

	for _, hint := range hints {
		ncurses.AddText(ncurses.COLOR_HINTS_SHORTCUT, y, x, hint[1])
		x += len(hint[1]) + 1
		ncurses.AddChar(ncurses.COLOR_HINTS_SHORTCUT, y, x-1, gc.ACS_BULLET)
		ncurses.AddText(ncurses.COLOR_HINTS_TEXT, y, x, hint[0])
		x += len(hint[0]) + 2
	}
}
