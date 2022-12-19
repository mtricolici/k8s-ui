package menu

import (
	"fmt"
	"k8s_ui/ncurses"
	"time"

	l "k8s_ui/logger"

	gc "github.com/rthornton128/goncurses"
)

const (
	MENU_TOP_LEFT_X = 10 // Vertical Column
	MENU_TOP_LEFT_Y = 2  // Horizontal Line
)

type (
	MenuHeaderFunc    func(screen *gc.Window)
	MenuHandleKeyFunc func(screen *gc.Window, key gc.Key, selectedItem []string) bool
)

func BuildMenuItems(items [][]string) []string {
	defer l.LogExecutedTime("BuildMenuItems")()

	menuItems := make([]string, len(items))
	max := make([]int, len(items[0]))

	for _, line := range items {
		for col, colValue := range line {
			if max[col] < len(colValue) {
				max[col] = len(colValue)
			}
		}
	}

	for i, line := range items {
		menuItems[i] = ""

		for col, colValue := range line {
			format := fmt.Sprintf("%s-%d%s", "%", max[col], "s")
			menuItems[i] += fmt.Sprintf(format, colValue)
			if col < len(line)-1 {
				menuItems[i] += " "
			}
		}
	}

	if len(menuItems) == 1 {
		menuItems[0] = " No resources found"
	}

	return menuItems
}

func drawVerticalLineTop(y int, x int, screen *gc.Window, count int) {
	defer l.LogExecutedTime("drawVerticalLineTop")()
	screen.MoveAddChar(y, x, gc.ACS_ULCORNER)
	for i := 0; i < count; i++ {
		screen.MoveAddChar(y, x+1+i, gc.ACS_HLINE)
	}
	screen.MoveAddChar(y, x+count+1, gc.ACS_URCORNER)
}

func drawVerticalLineBottom(y int, x int, screen *gc.Window, count int) {
	defer l.LogExecutedTime("drawVerticalLineBottom")()
	screen.MoveAddChar(y, x, gc.ACS_LLCORNER)
	for i := 0; i < count; i++ {
		screen.MoveAddChar(y, x+1+i, gc.ACS_HLINE)
	}
	screen.MoveAddChar(y, x+count+1, gc.ACS_LRCORNER)
}

func drawMenu(screen *gc.Window, items []string, selectedIndex int, drawIndexFrom int, drawIndexTo int) {
	defer l.LogExecutedTime("drawMenu")()
	if len(items) == 1 {
		screen.MovePrint(MENU_TOP_LEFT_Y, MENU_TOP_LEFT_X, items[0])
		return
	}

	windowHorizontalSize := len(items[0])

	x := MENU_TOP_LEFT_X
	y := MENU_TOP_LEFT_Y

	for i, item := range items {

		if i == 0 {
			screen.MovePrint(y, x+1, item)
			y++ // Move to next line
			drawVerticalLineTop(y, x, screen, windowHorizontalSize)
			y++ // Move to next line
			// screen.ColorOn(ncurses.COLOR_TABLE_HEADER)
			// screen.ColorOff(ncurses.COLOR_TABLE_HEADER)
			continue
		}

		if i < drawIndexFrom || i > drawIndexTo {
			continue // Ingore hidden items
		}

		if i == selectedIndex {
			screen.MoveAddChar(y, x, gc.ACS_VLINE)
			screen.ColorOn(ncurses.COLOR_SELECTED)
			screen.MovePrint(y, x+1, item)
			screen.ColorOff(ncurses.COLOR_SELECTED)
			screen.AddChar(gc.ACS_VLINE)
			y++ // Move to next line
		} else {
			screen.MoveAddChar(y, x, gc.ACS_VLINE)
			screen.MovePrint(y, x+1, item)
			screen.AddChar(gc.ACS_VLINE)
			y++ // Move to next line
		}
	}

	drawVerticalLineBottom(y, x, screen, windowHorizontalSize)
}

func ShowMenu(
	screen *gc.Window,
	items [][]string,
	menuItems []string,
	headerFunc MenuHeaderFunc,
	handleKeyFunc MenuHandleKeyFunc) {

	menuIdx := 1

	screen_max_lines, _ := screen.MaxYX()

	initial_index_from := 1
	initial_index_to := screen_max_lines - MENU_TOP_LEFT_Y - 10

	drawIndexFrom := initial_index_from
	drawIndexTo := initial_index_to

	for {
		screen.Clear()
		headerFunc(screen) // Draw custom header

		drawMenu(screen, menuItems, menuIdx, drawIndexFrom, drawIndexTo)

		screen.Refresh()
		key := screen.GetChar()
		if !handleKeyFunc(screen, key, items[menuIdx]) {
			switch key {
			case gc.KEY_DOWN:
				menuIdx++
				if menuIdx >= len(menuItems) {
					menuIdx = len(menuItems) - 1
				}
				if menuIdx > drawIndexTo {
					drawIndexFrom++
					drawIndexTo++
				}
			case gc.KEY_UP:
				menuIdx--
				if menuIdx < 1 {
					menuIdx = 1
				}
				if menuIdx < drawIndexFrom {
					drawIndexFrom--
					drawIndexTo--
				}
			case gc.KEY_ESC:
				return
			case gc.KEY_BACKSPACE:
				return
			default:
				ShowWarning(screen, "Warning: key not bound")
			}
		}
	}
}

func ShowWarning(screen *gc.Window, message string) {
	screen.ColorOn(ncurses.COLOR_WARNING)
	screen.MovePrint(0, 30, message)
	screen.ColorOff(ncurses.COLOR_WARNING)
	screen.Refresh()
	time.Sleep(300 * time.Millisecond)
}
