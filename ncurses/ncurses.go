package ncurses

import (
	"log"

	gc "github.com/rthornton128/goncurses"
)

const (
	COLOR_DEFAULT  int16 = 1
	COLOR_SELECTED int16 = 2
	COLOR_HEADER   int16 = 3
	COLOR_WARNING  int16 = 4
)

func Done() {
	gc.End()
}

func Init() *gc.Window {
	stdscr, err := gc.Init()

	if err != nil {
		log.Fatal(err)
	}

	if !gc.HasColors() {
		log.Fatal("Application requires a colour capable terminal")
	}

	if err := gc.StartColor(); err != nil {
		log.Fatal(err)
	}

	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Keypad(true)

	gc.InitPair(COLOR_DEFAULT, gc.C_WHITE, gc.C_BLACK)
	gc.InitPair(COLOR_SELECTED, gc.C_BLACK, gc.C_WHITE)
	gc.InitPair(COLOR_HEADER, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(COLOR_WARNING, gc.C_WHITE, gc.C_MAGENTA)

	stdscr.SetBackground(gc.ColorPair(1))

	return stdscr
}
