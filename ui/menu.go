package ui

import (
	"fmt"
	"time"

	l "k8s_ui/logger"
	"k8s_ui/ncurses"

	gc "github.com/rthornton128/goncurses"
)

const (
	top_left_x = 10 // Vertical Column
	top_left_y = 2  // Horizontal Line
)

type (
	MenuHeaderFunc    func()
	MenuHandleKeyFunc func(key gc.Key, selectedItem []string) bool
)

type Menu struct {
	screen *gc.Window
	data   [][]string
	items  []string

	Index         int
	FuncHeader    MenuHeaderFunc
	FuncHandleKey MenuHandleKeyFunc
}

func NewMenu(screen *gc.Window, data [][]string) *Menu {
	menu := Menu{
		screen:        screen,
		data:          data,
		items:         nil,
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
			format := fmt.Sprintf("%s-%d%s", "%", max[col], "s")
			m.items[i] += fmt.Sprintf(format, colValue)
			if col < len(line)-1 {
				m.items[i] += " "
			}
		}
	}

	if len(m.items) == 1 {
		m.items[0] = " No resources found"
	}
}

func (m *Menu) Show() {
	m.Index = 1

	screen_max_lines, _ := m.screen.MaxYX()

	initial_index_from := 1
	initial_index_to := screen_max_lines - top_left_y - 10

	drawIndexFrom := initial_index_from
	drawIndexTo := initial_index_to

	for {
		m.screen.Erase()
		m.FuncHeader() // Draw custom header

		m.draw(drawIndexFrom, drawIndexTo)

		m.screen.Refresh()
		key := m.screen.GetChar()

		keyHandled := false
		//TODO: investigate why sometimes it can be greater?!?!?!
		if m.Index < len(m.data) {
			keyHandled = m.FuncHandleKey(key, m.data[m.Index])
		}

		if !keyHandled {
			switch key {
			case gc.KEY_DOWN:
				m.Index++
				if m.Index >= len(m.items) {
					m.Index = len(m.items) - 1
				}
				if m.Index > drawIndexTo {
					drawIndexFrom++
					drawIndexTo++
				}
			case gc.KEY_UP:
				m.Index--
				if m.Index < 1 {
					m.Index = 1
				}
				if m.Index < drawIndexFrom {
					drawIndexFrom--
					drawIndexTo--
				}
			case gc.KEY_ESC:
				return
			case gc.KEY_BACKSPACE:
				return
			default:
				m.ShowWarning("Warning: key not bound")
			}
		}
	}
}

func (m *Menu) draw(drawIndexFrom int, drawIndexTo int) {
	defer l.LogExecutedTime("drawMenu")()
	if len(m.items) == 1 {
		m.screen.MovePrint(top_left_y, top_left_x, m.items[0])
		return
	}

	windowHorizontalSize := len(m.items[0])

	x := top_left_x
	y := top_left_y

	for i, item := range m.items {

		if i == 0 {
			m.screen.MovePrint(y, x+1, item)
			y++ // Move to next line
			m.drawVerticalLineTop(y, x, windowHorizontalSize)
			y++ // Move to next line
			continue
		}

		if i < drawIndexFrom || i > drawIndexTo {
			continue // Ingore hidden items
		}

		if i == m.Index {
			m.screen.MoveAddChar(y, x, gc.ACS_VLINE)
			m.screen.ColorOn(ncurses.COLOR_SELECTED)
			m.screen.MovePrint(y, x+1, item)
			m.screen.ColorOff(ncurses.COLOR_SELECTED)
			m.screen.AddChar(gc.ACS_VLINE)
			y++ // Move to next line
		} else {
			m.screen.MoveAddChar(y, x, gc.ACS_VLINE)
			m.screen.MovePrint(y, x+1, item)
			m.screen.AddChar(gc.ACS_VLINE)
			y++ // Move to next line
		}
	}

	m.drawVerticalLineBottom(y, x, windowHorizontalSize)
}

func (m *Menu) ShowWarning(message string) {
	m.screen.ColorOn(ncurses.COLOR_WARNING)
	m.screen.MovePrint(0, 30, message)
	m.screen.ColorOff(ncurses.COLOR_WARNING)
	m.screen.Refresh()
	time.Sleep(300 * time.Millisecond)
}

func (m *Menu) ShowError(err error) {
	m.screen.ColorOn(ncurses.COLOR_WARNING)
	m.screen.MovePrint(0, 30, err.Error())
	m.screen.ColorOff(ncurses.COLOR_WARNING)
	m.screen.Refresh()
	time.Sleep(300 * time.Millisecond)
}

func (m *Menu) drawVerticalLineTop(y int, x int, count int) {
	defer l.LogExecutedTime("drawVerticalLineTop")()
	m.screen.MoveAddChar(y, x, gc.ACS_ULCORNER)
	for i := 0; i < count; i++ {
		m.screen.MoveAddChar(y, x+1+i, gc.ACS_HLINE)
	}
	m.screen.MoveAddChar(y, x+count+1, gc.ACS_URCORNER)
}

func (m *Menu) drawVerticalLineBottom(y int, x int, count int) {
	defer l.LogExecutedTime("drawVerticalLineBottom")()
	m.screen.MoveAddChar(y, x, gc.ACS_LLCORNER)
	for i := 0; i < count; i++ {
		m.screen.MoveAddChar(y, x+1+i, gc.ACS_HLINE)
	}
	m.screen.MoveAddChar(y, x+count+1, gc.ACS_LRCORNER)
}