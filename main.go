package main

import (
	"fmt"
	"k8s_ui/k8s"
	"strings"
)

const (
	clearScreen string = "\033c"
	colorReset  string = "\033[0m"
	fgColor     string = "\u001b[30m"
	bgColor     string = "\u001b[47;1m"
)

func buildMenuItems(items [][]string) []string {
	menuItems := make([]string, len(items))

	for i, line := range items {
		menuItems[i] = strings.Join(line[:], " ")
	}

	return menuItems
}

func showMenu(items []string, selectedIndex int) {
	fmt.Print(clearScreen)
	for i, item := range items {
		if i == selectedIndex {
			fmt.Print(string(fgColor), string(bgColor), item, string(colorReset), "\n")
		} else {
			fmt.Println(item)
		}
	}
}

func main() {
	ns := k8s.K8s_namespaces()
	menuItems := buildMenuItems(ns)
	showMenu(menuItems, 2)
}
