package k8s

import (
	"fmt"
	"k8s_ui/menu"
	"k8s_ui/ncurses"

	gc "github.com/rthornton128/goncurses"
)

func ShowPods(screen *gc.Window) {

	pods := exec_get_pods(currentNamespace)

	menuItems := menu.BuildMenuItems(pods)

	menu.ShowMenu(screen, pods, menuItems, drawPodsHeader, handlePodsKey)
}

func drawPodsHeader(screen *gc.Window) {

	screen.ColorOn(ncurses.COLOR_HEADER)
	screen.Println(fmt.Sprintf("*** ns: '%s' pods:", currentNamespace))
	screen.ColorOff(ncurses.COLOR_HEADER)
}

func handlePodsKey(screen *gc.Window, key gc.Key, selectedItem []string) bool {
	msg := fmt.Sprintf("%s", selectedItem)
	menu.ShowWarning(screen, msg)
	return false
}
