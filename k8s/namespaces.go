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
	client := NewK8SClient()

	namespaces, _ := client.GetNamespaces()
	menuItems := menu.BuildMenuItems(namespaces)
	namespacesCount = len(menuItems) - 1

	menu.ShowMenu(screen, namespaces, menuItems, drawNamespacesHeader, handleNamespacesKey)
}

func drawNamespacesHeader(screen *gc.Window) {

	screen.ColorOn(ncurses.COLOR_HEADER)
	screen.Printf(" Namespaces: %d of %d", menu.MenuCurrentItemIndex, namespacesCount)
	screen.ColorOff(ncurses.COLOR_HEADER)
}

func handleNamespacesKey(screen *gc.Window, key gc.Key, selectedItem []string) bool {

	if key == gc.KEY_RETURN {
		currentNamespace = selectedItem[0]
		saveIndex := menu.MenuCurrentItemIndex
		ShowPods(screen)
		// ugly yaaa :D sorry. that's open source lol
		menu.MenuCurrentItemIndex = saveIndex
		return true
	}

	return false
}
