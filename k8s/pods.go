package k8s

// import (
// 	"k8s_ui/menu"
// 	"k8s_ui/ncurses"

// 	gc "github.com/rthornton128/goncurses"
// )

// var (
// 	podsCount int
// )

// func ShowPods(screen *gc.Window) {

// 	client := NewK8SClient()
// 	pods, _ := client.GetPods(currentNamespace)

// 	menuItems := menu.BuildMenuItems(pods)
// 	podsCount = len(menuItems) - 1

// 	menu.ShowMenu(screen, pods, menuItems, drawPodsHeader, handlePodsKey)
// }

// func drawPodsHeader(screen *gc.Window) {

// 	screen.ColorOn(ncurses.COLOR_HEADER)
// 	screen.Printf("Ns '%s' pods: %d of %d", currentNamespace, menu.MenuCurrentItemIndex, podsCount)
// 	screen.ColorOff(ncurses.COLOR_HEADER)
// }

// func handlePodsKey(screen *gc.Window, key gc.Key, selectedItem []string) bool {
// 	// msg := fmt.Sprintf("%s", selectedItem)
// 	// menu.ShowWarning(screen, msg)
// 	return false
// }
