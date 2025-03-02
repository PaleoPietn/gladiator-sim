package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// ShowStartScreen displays the welcome screen and gets the player's name
func ShowStartScreen(screen tcell.Screen) string {
	screen.Clear()

	defaultStyle := tcell.StyleDefault
	titleStyle := defaultStyle.Bold(true).Foreground(tcell.ColorYellow)
	promptStyle := defaultStyle.Foreground(tcell.ColorWhite)
	inputStyle := defaultStyle.Foreground(tcell.ColorGreen)

	printText := func(x, y int, text string, style tcell.Style) {
		posX := x
		for _, c := range text {
			width := runewidth.RuneWidth(c)
			screen.SetContent(posX, y, c, nil, style)
			posX += width
		}
	}

	// Draw title and prompt
	printText(10, 5, "WELCOME TO ROGUELIKE GLADIATOR ARENA", titleStyle)
	printText(10, 8, "Enter your name, brave warrior:", promptStyle)

	playerName := ""

	// Draw input field
	drawInputField := func() {
		screen.Clear()
		printText(10, 5, "WELCOME TO ROGUELIKE GLADIATOR ARENA", titleStyle)
		printText(10, 8, "Enter your name, brave warrior:", promptStyle)
		printText(10, 10, playerName, inputStyle)

		// Draw cursor
		screen.SetContent(10+len(playerName), 10, '_', nil, inputStyle)

		printText(10, 14, "Press ENTER when done", promptStyle)
		screen.Show()
	}

	drawInputField()

	// Input handling loop
	for {
		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEnter:
				if strings.TrimSpace(playerName) == "" {
					return "Hero"
				}
				return playerName

			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(playerName) > 0 {
					playerName = playerName[:len(playerName)-1]
				}

			case tcell.KeyEscape:
				return "Hero"

			default:
				if ev.Key() == tcell.KeyRune {
					if len(playerName) < 30 {
						playerName += string(ev.Rune())
					}
				}
			}

			drawInputField()

		case *tcell.EventResize:
			screen.Sync()
			drawInputField()
		}
	}
}
