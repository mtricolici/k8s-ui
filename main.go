package main

import (
	"k8s_ui/k8s"
	"k8s_ui/ncurses"
)

func main() {
	screen := ncurses.Init()
	defer ncurses.Done()

	k8s.ShowNamespaces(screen)
}
