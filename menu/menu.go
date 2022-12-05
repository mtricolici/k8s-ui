package menu

import "fmt"

const (
	clearScreen string = "\033c"
	colorReset  string = "\033[0m"
	fgColor     string = "\u001b[30m"
	bgColor     string = "\u001b[47;1m"
)

func BuildMenuItems(items [][]string) []string {
	menuItems := make([]string, len(items))
	max := make([]int, len(items))

	for _, line := range items {
		for col, colValue := range line {
			if max[col] < len(colValue) {
				max[col] = len(colValue)
			}
		}
	}

	for i, line := range items {
		menuItems[i] = ""

		for col, colValue := range line {
			format := fmt.Sprintf("%s-%d%s", "%", max[col], "s")
			menuItems[i] += fmt.Sprintf(format, colValue)
			if col < len(line)-1 {
				menuItems[i] += " "
			}
		}
	}

	return menuItems
}

func ClearScreen() {
	fmt.Print(clearScreen)
}

func ShowMenu(items []string, selectedIndex int) {
	for i, item := range items {
		if i == selectedIndex {
			fmt.Print(string(fgColor), string(bgColor), item, string(colorReset), "\n")
		} else {
			fmt.Println(item)
		}
	}
}
