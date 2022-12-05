package main

import (
	"fmt"
	"k8s_ui/k8s"
	"k8s_ui/menu"
	"k8s_ui/ncurses"

	gc "github.com/rthornton128/goncurses"
)

func main() {
	scr := ncurses.Init()
	defer ncurses.Done()

	namespaces := k8s.K8s_namespaces()
	menuItems := menu.BuildMenuItems(namespaces)

	menu.ShowMenu(scr, namespaces, menuItems, drawHeader, handleKey)
}

func drawHeader(screen *gc.Window) {

	screen.ColorOn(ncurses.COLOR_HEADER)
	screen.Println("*** Namespaces:")
	screen.ColorOff(ncurses.COLOR_HEADER)
}

func handleKey(screen *gc.Window, key gc.Key, selectedItem []string) bool {
	msg := fmt.Sprintf("%s", selectedItem)
	menu.ShowWarning(screen, msg)
	return false
}
