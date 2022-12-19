package main

import (
	"k8s_ui/k8s"
	l "k8s_ui/logger"
	"k8s_ui/ncurses"
)

func main() {
	l.Init()
	defer l.Close()

	screen := ncurses.Init()
	defer ncurses.Done()

	k8s.ShowNamespaces(screen)
}
