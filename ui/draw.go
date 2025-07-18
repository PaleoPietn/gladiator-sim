package ui

import (
	"fmt"
	"strings"

	"gladiator-sim/models"
	model "gladiator-sim/models"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

const (
	// UI
	healthBarWidth = 20
	logStartY      = 8
	logStartx      = 2
	maxLogEntries  = 20
	heroXIndex     = 2
	enemyXIndex    = 52
	playerYIndex   = 3

	// Texts
	titleText        = "ROGUELIKE GLADIATOR ARENA"
	battleLogText    = "BATTLE LOG:"
	upgradeText      = "CHOOSE YOUR UPGRADE:"
	HPDisplay        = "HP: %*s"
	HPRatioDisplay   = "%d/%d"
	prefixUnselected = "   "
	prefixSelected   = ">> "
	upgradeHelper    = "Use UP/DOWN arrows to select, ENTER to confirm"
	gameOverHelper   = "Game Over! Press 'q' to exit or 'r' to start a new run."
	quitHelper       = "Press 'q' to quit."
	endGameText      = "Thank you for playing ROGUELIKE GLADIATOR ARENA, a game made by PaleoPietn Studios!"
)

// drawHealthBar creates a visual health bar. Max should be above 0.
func generateHealthBar(current, max int, width int) string {
	if max < 1 {
		max = 1
	}

	filledWidth := int(float64(current) / float64(max) * float64(width))
	if filledWidth < 0 {
		filledWidth = 0
	}

	return fmt.Sprintf("[%s%s]", strings.Repeat("█", filledWidth), strings.Repeat("░", width-filledWidth))
}

func generateBuffsString(player *model.Player) string {
	// TODO: create an array of buffs that were selected to print them here
	buffs := ""
	if player.Regeneration > 0 {
		buffs += " 🌿"
	}

	return buffs
}

// Helper function to draw text with proper handling of wide characters
func (ui *UI) printText(x, y int, text string, style tcell.Style) {
	posX := x
	for _, c := range text {
		// Check if character is an emoji or other wide character
		width := runewidth.RuneWidth(c)
		ui.screen.SetContent(posX, y, c, nil, style)
		posX += width
	}
}

func (ui *UI) drawUpgrades(startIndex, selected int, upgrades []models.Upgrade) {
	for i, upgrade := range upgrades {
		if i == selected {
			ui.printText(2, startIndex+i+1, fmt.Sprintf("%s%d. %s - %s", prefixSelected, i+1, upgrade.Name, upgrade.Description), selectedStyle)
			continue
		}
		ui.printText(2, startIndex+i+1, fmt.Sprintf("%s%d. %s - %s", prefixUnselected, i+1, upgrade.Name, upgrade.Description), infoStyle)
	}
}

// Draws a Player
func (ui *UI) drawPlayer(player *model.Player) {
	xIndex := heroXIndex
	style := heroStyle
	startYIndex := playerYIndex

	if !player.IsHero {
		xIndex = enemyXIndex
		style = enemyStyle
	}

	healthBar := generateHealthBar(player.Health, player.MaxHealth, healthBarWidth)
	ui.printText(xIndex, startYIndex, fmt.Sprintf("%s %s", player.Name, generateBuffsString(player)), style)
	startYIndex++
	ui.printText(xIndex, startYIndex, fmt.Sprintf("%s %s", formatLifeCount(player.Health, player.MaxHealth), healthBar), style)
	startYIndex++
	ui.printText(xIndex, startYIndex, fmt.Sprintf("ATK: %d-%d | DEF: %d | Wins: %d", player.AttackMin, player.AttackMax, player.Defense, player.Wins), style)
}

func formatLifeCount(health, maxHealth int) string {
	minWidth := 7
	return fmt.Sprintf(HPDisplay, minWidth, fmt.Sprintf(HPRatioDisplay, health, maxHealth))
}

// Draws the entire log entry
func (ui *UI) drawLogs(logs []*models.BattleLog) {
	for k, v := range logs {
		ui.printText(logStartx, logStartY+k, v.LogMessage, mapLogType(v.LogType))
	}
	ui.lastUsedYIndex = logStartY + len(logs)
}
