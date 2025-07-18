package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

const (
	// UI
	startScreenXindex            = 10
	startScreenTitleYIndex       = 5
	startScreenPromptYIndex      = 8
	startScreenPlayerInputYIndex = 10
	startScreenHelperYIndex      = 14

	// Texts
	startScreenTitle  = "WELCOME TO ROGUELIKE GLADIATOR ARENA"
	startScreenPrompt = "Enter your name, brave warrior:"
	startScreenHelper = "Press ENTER when done"
)

// ShowStartScreen displays the welcome screen and gets the player's name
func (ui *UI) StartScreen() string {
	ui.screen.Clear()

	// Draw input field
	playerName := ""
	ui.drawInputField(playerName)

	// Input handling loop
	for {
		ev := ui.screen.PollEvent()

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

			ui.drawInputField(playerName)

		case *tcell.EventResize:
			ui.screen.Sync()
			ui.drawInputField(playerName)
		}
	}
}

func (ui *UI) drawInputField(playerName string) {
	ui.screen.Clear()
	ui.printText(startScreenXindex, startScreenTitleYIndex, startScreenTitle, titleStyle)
	ui.printText(startScreenXindex, startScreenPromptYIndex, startScreenPrompt, infoStyle)
	ui.printText(startScreenXindex, startScreenPlayerInputYIndex, playerName, heroStyle)

	// Draw cursor
	ui.screen.SetContent(startScreenXindex+len(playerName), startScreenPlayerInputYIndex, '_', nil, heroStyle)

	ui.printText(startScreenXindex, startScreenHelperYIndex, startScreenHelper, infoStyle)
	ui.screen.Show()
}
