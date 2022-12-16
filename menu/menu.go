package menu

import (
	"fmt"
	"k8s_ui/ncurses"
	"time"

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
	screen.MoveAddChar(y, x, gc.ACS_ULCORNER)
	for i := 0; i < count; i++ {
		screen.MoveAddChar(y, x+1+i, gc.ACS_HLINE)
	}
	screen.MoveAddChar(y, x+count+1, gc.ACS_URCORNER)
}

func drawVerticalLineBottom(y int, x int, screen *gc.Window, count int) {
	screen.MoveAddChar(y, x, gc.ACS_LLCORNER)
	for i := 0; i < count; i++ {
		screen.MoveAddChar(y, x+1+i, gc.ACS_HLINE)
	}
	screen.MoveAddChar(y, x+count+1, gc.ACS_LRCORNER)
}

func drawMenu(screen *gc.Window, items []string, selectedIndex int) {
	if len(items) == 1 {
		screen.MovePrint(MENU_TOP_LEFT_Y, MENU_TOP_LEFT_X, items[0])
		return
	}

	windowHorizontalSize := len(items[0])

	for i, item := range items {
		x := MENU_TOP_LEFT_X + 1 // column +1 (for border)
		y := MENU_TOP_LEFT_Y + i // line +1 (for border)

		if i == 0 {
			// screen.ColorOn(ncurses.COLOR_TABLE_HEADER)
			screen.MovePrint(y, x, item)
			// screen.ColorOff(ncurses.COLOR_TABLE_HEADER)
			screen.Println()
			drawVerticalLineTop(y+1, x-1, screen, windowHorizontalSize)
			screen.Println()
			continue
		} else if i == selectedIndex {
			screen.MoveAddChar(y+1, x-1, gc.ACS_VLINE)
			screen.ColorOn(ncurses.COLOR_SELECTED)
			screen.MovePrint(y+1, x, item)
			screen.ColorOff(ncurses.COLOR_SELECTED)
			screen.AddChar(gc.ACS_VLINE)
		} else {
			screen.MoveAddChar(y+1, x-1, gc.ACS_VLINE)
			screen.MovePrint(y+1, x, item)
			screen.AddChar(gc.ACS_VLINE)
		}
		screen.Println()
	}

	x := MENU_TOP_LEFT_X
	y := MENU_TOP_LEFT_Y + len(items) + 1
	drawVerticalLineBottom(y, x, screen, windowHorizontalSize)
}

func ShowMenu(
	screen *gc.Window,
	items [][]string,
	menuItems []string,
	headerFunc MenuHeaderFunc,
	handleKeyFunc MenuHandleKeyFunc) {

	menuIdx := 1

	for {
		screen.Clear()
		headerFunc(screen) // Draw custom header

		drawMenu(screen, menuItems, menuIdx)

		screen.Refresh()
		key := screen.GetChar()
		switch key {
		case gc.KEY_DOWN:
			menuIdx++
			if menuIdx >= len(menuItems) {
				menuIdx = 1
			}
		case gc.KEY_UP:
			menuIdx--
			if menuIdx < 1 {
				menuIdx = len(menuItems) - 1
			}
		case gc.KEY_ESC:
			return
		default:
			if !handleKeyFunc(screen, key, items[menuIdx]) {
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
