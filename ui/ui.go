package ui

import (
	"gladiator-sim/models"
	"time"

	"github.com/gdamore/tcell/v2"
)

type UI struct {
	screen         tcell.Screen
	quit           chan struct{}
	lastUsedYIndex int // very ugly hack but we can do something better when we revamp the UI
}

func NewUI(screen tcell.Screen, quit chan struct{}) *UI {
	return &UI{
		screen: screen,
		quit:   quit,
	}
}

// Draws the battle based on the battle states provided, returns the index of the upgrade selected
// Could be decoupled if we throw a chan to signal when the logs are printed or create a listener in the engine for "BattleFinished"
func (ui *UI) DrawBattle(states []*models.BattleState) {
	logs := []*models.BattleLog{}
	ui.screen.Clear()

	for _, state := range states {
		// Refresh players
		ui.drawPlayer(state.Hero)
		ui.drawPlayer(state.Enemy)

		// To avoid displaying too many logs, we pop the first entry of the logs out
		if len(logs) == maxLogEntries {
			logs = logs[1:]
		}
		logs = append(logs, state.BattleLog)
		ui.drawLogs(logs)

		if state.BattleLog.LogWaitTime == 0 {
			state.BattleLog.LogWaitTime = models.LogWaitTimeDefault
		}
		ui.screen.Show()
		time.Sleep(state.BattleLog.LogWaitTime)
	}
}

func (ui *UI) ChooseUpgrade(upgrades []models.Upgrade) int {
	selected := len(upgrades) / 2

	for {
		ui.printText(logStartx, ui.lastUsedYIndex+1, "", defaultStyle)
		ui.printText(logStartx, ui.lastUsedYIndex+2, upgradeText, titleStyle)

		ui.drawUpgrades(ui.lastUsedYIndex+3, selected, upgrades)

		ui.printText(logStartx, ui.lastUsedYIndex+5+len(upgrades), upgradeHelper, infoStyle)

		ui.screen.Show()

		ev := ui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				if selected > 0 {
					selected--
				}
			case tcell.KeyDown:
				if selected < len(upgrades)-1 {
					selected++
				}
			case tcell.KeyEnter:
				ui.lastUsedYIndex += 4
				return selected
			case tcell.KeyRune, tcell.KeyEscape:
				if ev.Rune() == 'q' {
					ui.lastUsedYIndex += 4
					<-ui.quit
				}
			}
		case *tcell.EventResize:
			ui.screen.Sync()
		}
	}

}

// Displays a prompt asking if the user wants to restart or quit
func (ui *UI) DrawEndGameScreen() bool {
	ui.printText(logStartx, ui.lastUsedYIndex, gameOverHelper, infoStyle)

	for {
		ev := ui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				switch {
				case ev.Rune() == 'q':
					return true
				case ev.Rune() == 'r':
					return false
				}
			case tcell.KeyEscape:
				return true
			}
		case *tcell.EventResize:
			ui.screen.Sync()
		}
	}
}

// Draws the quitting screen
func (ui *UI) DrawQuitScreen() {
	ui.screen.Clear()
	ui.printText(startScreenXindex, startScreenTitleYIndex, endGameText, titleStyle)
}
