package main

import (
	l "k8s_ui/logger"
	"k8s_ui/ncurses"
	"k8s_ui/ui"
)

func main() {
	l.Init()
	defer l.Close()

	screen := ncurses.Init()
	defer ncurses.Done()

	menu := ui.NewMenuNamespaces(screen)
	menu.Load()
	menu.Show()
}
