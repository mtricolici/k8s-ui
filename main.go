package main

import (
	"k8s_ui/k8s"
	"k8s_ui/menu"
)

func main() {
	ns := k8s.K8s_namespaces()
	menuItems := menu.BuildMenuItems(ns)

	menu.ClearScreen()
	menu.ShowMenu(menuItems, 2)
}
