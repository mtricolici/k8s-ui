package main

import (
	"k8s_ui/k8s"
	"k8s_ui/menu"
	"k8s_ui/ncurses"
	"time"

	gc "github.com/rthornton128/goncurses"
)

func main() {
	scr := ncurses.Init()
	defer ncurses.Done()

	ns := k8s.K8s_namespaces()
	menuItems := menu.BuildMenuItems(ns)

	shouldStop := false
	menuIdx := 0

	for {
		scr.Clear()
		scr.ColorOn(ncurses.COLOR_HEADER)
		scr.Println("*** Namespaces:")
		scr.ColorOff(ncurses.COLOR_HEADER)

		menu.ShowMenu(scr, menuItems, menuIdx)

		scr.Refresh()
		key := scr.GetChar()
		switch key {
		case gc.KEY_DOWN:
			menuIdx++
			if menuIdx >= len(menuItems) {
				menuIdx = 0
			}
		case gc.KEY_UP:
			menuIdx--
			if menuIdx < 0 {
				menuIdx = len(menuItems) - 1
			}
		case gc.KEY_ESC:
			shouldStop = true
		default:
			scr.ColorOn(ncurses.COLOR_WARNING)
			scr.MovePrint(0, 30, "Warning: key not bound")
			scr.ColorOff(ncurses.COLOR_WARNING)
			scr.Refresh()
			time.Sleep(300 * time.Millisecond)
		}

		if shouldStop {
			break
		}
	}
}
