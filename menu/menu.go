package menu

import (
	"fmt"
	"k8s_ui/ncurses"
	"time"

	gc "github.com/rthornton128/goncurses"
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

	return menuItems
}

func drawMenu(screen *gc.Window, items []string, selectedIndex int) {
	for i, item := range items {
		if i == 0 {
			screen.ColorOn(ncurses.COLOR_TABLE_HEADER)
			screen.Println(item)
			screen.ColorOff(ncurses.COLOR_TABLE_HEADER)
		} else if i == selectedIndex {
			screen.ColorOn(ncurses.COLOR_SELECTED)
			screen.Println(item)
			screen.ColorOff(ncurses.COLOR_SELECTED)
		} else {
			screen.Println(item)
		}
	}
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
