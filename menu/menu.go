package menu

import (
	"fmt"
	"k8s_ui/ncurses"

	gc "github.com/rthornton128/goncurses"
)

func BuildMenuItems(items [][]string) []string {
	menuItems := make([]string, len(items))
	max := make([]int, len(items))

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

func ShowMenu(screen *gc.Window, items []string, selectedIndex int) {
	for i, item := range items {
		if i == selectedIndex {
			screen.ColorOn(ncurses.COLOR_SELECTED)
			screen.Println(item)
			screen.ColorOff(ncurses.COLOR_SELECTED)
		} else {
			screen.Println(item)
		}
	}
}
