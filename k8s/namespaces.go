package k8s

import (
	"k8s_ui/menu"
	"k8s_ui/ncurses"

	gc "github.com/rthornton128/goncurses"
)

var (
	currentNamespace string
	namespacesCount  int
)

func ShowNamespaces(screen *gc.Window) {

	namespaces := exec_get_namespaces()
	menuItems := menu.BuildMenuItems(namespaces)
	namespacesCount = len(menuItems) - 1

	menu.ShowMenu(screen, namespaces, menuItems, drawNamespacesHeader, handleNamespacesKey)
}

func drawNamespacesHeader(screen *gc.Window) {

	screen.ColorOn(ncurses.COLOR_HEADER)
	screen.Printf(" Namespaces: %d\n", namespacesCount)
	screen.ColorOff(ncurses.COLOR_HEADER)
}

func handleNamespacesKey(screen *gc.Window, key gc.Key, selectedItem []string) bool {
	if key == gc.KEY_BACKSPACE {
		// don't do anything, just ignore this key
		return true
	}

	if key == gc.KEY_RETURN {
		currentNamespace = selectedItem[0]
		ShowPods(screen)
		return true
	}
	return false
}
