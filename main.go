package main

import (
	"k8s_ui/k8s"
	"k8s_ui/menu"
	"k8s_ui/ncurses"
)

func main() {
	scr := ncurses.Init()
	defer ncurses.Done()

	ns := k8s.K8s_namespaces()
	menuItems := menu.BuildMenuItems(ns)

	scr.Clear()
	scr.ColorOn(ncurses.COLOR_HEADER)
	scr.Println("*** Namespaces:")
	scr.ColorOff(ncurses.COLOR_HEADER)

	menu.ShowMenu(scr, menuItems, 2)

	scr.Refresh()
	scr.GetChar()
}
