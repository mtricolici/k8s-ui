package ncurses

import (
	"log"
	"time"

	gc "github.com/rthornton128/goncurses"
)

const (
	COLOR_DEFAULT            int16 = 1
	COLOR_HEADER             int16 = 2
	COLOR_MENU_HEADER        int16 = 3
	COLOR_MENU_ITEM          int16 = 4
	COLOR_MENU_ITEM_SELECTED int16 = 5
	COLOR_HINTS_TEXT         int16 = 6
	COLOR_HINTS_SHORTCUT     int16 = 7
	COLOR_MESSAGEBOX         int16 = 8
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
	gc.InitPair(COLOR_HEADER, gc.C_BLACK, gc.C_WHITE)
	gc.InitPair(COLOR_MENU_HEADER, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(COLOR_MENU_ITEM, gc.C_WHITE, gc.C_BLUE)
	gc.InitPair(COLOR_MENU_ITEM_SELECTED, gc.C_BLACK, gc.C_CYAN)

	gc.InitPair(COLOR_HINTS_TEXT, gc.C_WHITE, gc.C_BLACK)
	gc.InitPair(COLOR_HINTS_SHORTCUT, gc.C_YELLOW, gc.C_BLACK)

	gc.InitPair(COLOR_MESSAGEBOX, gc.C_WHITE, gc.C_CYAN)

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

func HLine(color int16, y, x int, ach gc.Char, width int) {
	screen.ColorOn(color)
	screen.HLine(y, x, ach, width)
	screen.ColorOff(color)
}

// Shows a dialog and don't wait. you should close the window with .Delete() !
func MessageBoxAsync(title, message string) *gc.Window {
	max_y, max_x := screen.MaxYX()

	win_height := 3
	win_width := 2 + len(message)
	win_pos_x := max_x/2 - win_width/2
	win_pos_y := max_y/2 - win_height/2

	win, _ := gc.NewWindow(win_height, win_width, win_pos_y, win_pos_x)
	win.ColorOn(COLOR_MESSAGEBOX)
	win.Border(gc.ACS_VLINE, gc.ACS_VLINE, gc.ACS_HLINE, gc.ACS_HLINE,
		gc.ACS_ULCORNER, gc.ACS_URCORNER, gc.ACS_LLCORNER, gc.ACS_LRCORNER)

	if len(title) > 0 {
		win.MovePrintf(0, 2, " %s ", title)
	}

	win.MovePrint(1, 1, message)
	win.ColorOff(COLOR_MESSAGEBOX)
	win.NoutRefresh()
	gc.Update()
	return win
}

func MessageBox(title, message string, duration int) {
	win := MessageBoxAsync(title, message)
	time.Sleep(time.Duration(duration) * time.Millisecond)
	win.Delete()
}
