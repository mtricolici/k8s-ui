package ncurses

import (
	"fmt"
	"k8s_ui/utils"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	gc "github.com/rthornton128/goncurses"
)

const (
	COLOR_DEFAULT                      int16 = 1
	COLOR_HEADER                       int16 = 2
	COLOR_HEADER_HIGH                  int16 = 3
	COLOR_HEADER_HINT                  int16 = 4
	COLOR_MENU_HEADER                  int16 = 5
	COLOR_MENU_HEADER_DELIMITER        int16 = 6
	COLOR_MENU_ITEM                    int16 = 7
	COLOR_MENU_ITEM_DELIMITER          int16 = 8
	COLOR_MENU_ITEM_SELECTED           int16 = 9
	COLOR_MENU_ITEM_SELECTED_DELIMITER int16 = 10
	COLOR_HINTS_TEXT                   int16 = 11
	COLOR_HINTS_SHORTCUT               int16 = 12
	COLOR_MESSAGEBOX                   int16 = 13
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
	gc.SetEscDelay(200)
	stdscr.Keypad(true)

	delimiter := int16(gc.C_GREEN)

	gc.InitPair(COLOR_DEFAULT, gc.C_WHITE, gc.C_BLACK)
	gc.InitPair(COLOR_HEADER, gc.C_BLACK, gc.C_WHITE)
	gc.InitPair(COLOR_HEADER_HIGH, gc.C_GREEN, gc.C_WHITE)
	gc.InitPair(COLOR_HEADER_HINT, gc.C_YELLOW, gc.C_WHITE)

	gc.InitPair(COLOR_MENU_HEADER, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(COLOR_MENU_HEADER_DELIMITER, delimiter, gc.C_BLACK)
	gc.InitPair(COLOR_MENU_ITEM, gc.C_WHITE, gc.C_BLUE)
	gc.InitPair(COLOR_MENU_ITEM_DELIMITER, delimiter, gc.C_BLUE)
	gc.InitPair(COLOR_MENU_ITEM_SELECTED, gc.C_BLACK, gc.C_CYAN)
	gc.InitPair(COLOR_MENU_ITEM_SELECTED_DELIMITER, delimiter, gc.C_CYAN)

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

func AddTextMaxWidth(color int16, y, x, max_width int, text string) {
	screen.ColorOn(color)
	screen.MovePrint(y, x, utils.ShortString(text, max_width))
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

func InputDialog(title string, max int) (bool, string) {
	win := MessageBoxAsync(title, strings.Repeat(" ", max+1))
	win.Keypad(true)
	win.Move(1, 1)

	// show cursor
	gc.Cursor(1)

	var input string
	success := true

	for {
		ch := win.GetChar()
		switch ch {
		case gc.KEY_ESC:
			input = ""
			success = false
			goto input_done
		case gc.KEY_RETURN, gc.KEY_ENTER:
			success = true
			goto input_done
		case gc.KEY_BACKSPACE: //TODO: 8 too ? test on other terminal types
			if len(input) > 0 {
				input = input[:len(input)-1]

				win.ColorOn(COLOR_MESSAGEBOX)
				win.HLine(1, 1, ' ', max)
				win.ColorOff(COLOR_MESSAGEBOX)
				win.MovePrint(1, 1, input)
				win.NoutRefresh()
				gc.Update()
			}
		default:
			if ch >= 32 && ch <= 126 {
				if len(input) < max {
					input += string(rune(ch))
					win.ColorOn(COLOR_MESSAGEBOX)
					win.HLine(1, 1, ' ', max)
					win.ColorOff(COLOR_MESSAGEBOX)
					win.MovePrint(1, 1, input)
					win.NoutRefresh()
					gc.Update()
				}
			}
		}
	}

input_done:
	// Hide cursor
	gc.Cursor(0)

	win.Delete()
	return success, input
}

func Clear_screen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ExecuteCommand(command string) {
	gc.StdScr().Clear()
	gc.End()
	Clear_screen()
	fmt.Printf("Running: %s\n", command)

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		time.Sleep(3 * time.Second)
	}
}
