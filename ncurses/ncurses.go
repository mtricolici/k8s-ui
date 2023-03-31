package ncurses

import (
	"log"

	gc "github.com/rthornton128/goncurses"
)

const (
	COLOR_DEFAULT            int16 = 1
	COLOR_TITLE              int16 = 2
	COLOR_WARNING            int16 = 3
	COLOR_HEADER             int16 = 4
	COLOR_MENU_ITEM          int16 = 5
	COLOR_MENU_ITEM_SELECTED int16 = 6
)

var (
	screen *gc.Window
)

func Done() {
	gc.End()
}

func Init() *gc.Window {
	stdscr, err := gc.Init()

	if err != nil {
		log.Panic(err)
	}

	if !gc.HasColors() {
		log.Panic("Application requires a colour capable terminal")
	}

	if err := gc.StartColor(); err != nil {
		log.Panic(err)
	}

	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	gc.SetEscDelay(50)
	stdscr.Keypad(true)

	gc.InitPair(COLOR_DEFAULT, gc.C_WHITE, gc.C_BLACK)
	gc.InitPair(COLOR_TITLE, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(COLOR_WARNING, gc.C_WHITE, gc.C_MAGENTA)
	gc.InitPair(COLOR_HEADER, gc.C_BLACK, gc.C_WHITE)
	gc.InitPair(COLOR_MENU_ITEM, gc.C_WHITE, gc.C_BLUE)
	gc.InitPair(COLOR_MENU_ITEM_SELECTED, gc.C_BLACK, gc.C_CYAN)

	stdscr.SetBackground(gc.ColorPair(1))

	screen = stdscr

	return stdscr
}

func AddChar(color int16, y, x int, ach gc.Char) {
	screen.ColorOn(color)
	screen.MoveAddChar(y, x, ach)
	screen.ColorOff(color)
}

func AddText(color int16, y, x int, text string) {
	screen.ColorOn(color)
	screen.MovePrint(y, x, text)
	screen.ColorOff(color)
}
