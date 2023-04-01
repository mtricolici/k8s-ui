package main

import (
	"k8s_ui/ncurses"
	"k8s_ui/ui"
)

func main() {
	screen := ncurses.Init()
	defer ncurses.Done()

	menu := ui.NewMenuNamespaces(screen)
	err := menu.Load()
	if err != nil {
		ncurses.MessageBox("Error", err.Error(), 2000)
	} else {
		menu.Show()
	}
}
