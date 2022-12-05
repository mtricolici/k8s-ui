package k8s

import (
	"k8s_ui/menu"
	"k8s_ui/ncurses"

	gc "github.com/rthornton128/goncurses"
)

func ShowNamespaces(screen *gc.Window) {

	namespaces := exec_get_namespaces()
	menuItems := menu.BuildMenuItems(namespaces)

	menu.ShowMenu(screen, namespaces, menuItems, drawNamespacesHeader, handleNamespacesKey)
}

func drawNamespacesHeader(screen *gc.Window) {

	screen.ColorOn(ncurses.COLOR_HEADER)
	screen.Println("*** Namespaces:")
	screen.ColorOff(ncurses.COLOR_HEADER)
}

func handleNamespacesKey(screen *gc.Window, key gc.Key, selectedItem []string) bool {
	if key == gc.KEY_RETURN {
		ShowPods(screen, selectedItem[0])
		return true
	}
	return false
}
