// ui/draw.go
package ui

import (
	"fmt"
	"strings"

	model "stuff/models"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// DrawHealthBar creates a visual health bar
func DrawHealthBar(current, max int, width int) string {
	if max <= 0 {
		return "[ERROR]"
	}

	filledWidth := int(float64(current) / float64(max) * float64(width))
	if filledWidth < 0 {
		filledWidth = 0
	}

	bar := "["
	bar += strings.Repeat("█", filledWidth)
	bar += strings.Repeat("░", width-filledWidth)
	bar += "]"

	return bar
}

// DrawUI renders the game interface to the screen
func DrawUI(screen tcell.Screen, hero *model.Player, enemy *model.Player, gameState *model.GameState) {
	buffs := ""
	if hero.Regeneration > 0 {
		buffs += " 🌿"
	}

	screen.Clear()

	// Helper function to draw text with proper handling of wide characters
	printText := func(x, y int, text string, style tcell.Style) {
		posX := x
		for _, c := range text {
			// Check if character is an emoji or other wide character
			width := runewidth.RuneWidth(c)
			screen.SetContent(posX, y, c, nil, style)
			posX += width // Advance by the actual width
		}
	}

	// Define styles
	defaultStyle := tcell.StyleDefault
	titleStyle := defaultStyle.Bold(true).Foreground(tcell.ColorYellow)
	heroStyle := defaultStyle.Foreground(tcell.ColorGreen)
	enemyStyle := defaultStyle.Foreground(tcell.ColorRed)
	infoStyle := defaultStyle.Foreground(tcell.ColorWhite)
	selectedStyle := defaultStyle.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite)
	criticalStyle := defaultStyle.Foreground(tcell.ColorYellow)
	blockStyle := defaultStyle.Foreground(tcell.ColorTeal)

	// Draw title and stats
	printText(2, 1, "ROGUELIKE GLADIATOR ARENA", titleStyle)

	// Draw player stats with health bars
	healthBarWidth := 20
	heroHealthBar := DrawHealthBar(hero.Health, hero.MaxHealth, healthBarWidth)
	enemyHealthBar := DrawHealthBar(enemy.Health, enemy.MaxHealth, healthBarWidth)

	printText(2, 3, fmt.Sprintf("%s %s", hero.Name, buffs), heroStyle)
	printText(2, 4, fmt.Sprintf("HP: %d/%d %s", hero.Health, hero.MaxHealth, heroHealthBar), heroStyle)
	printText(2, 5, fmt.Sprintf("ATK: %d-%d | DEF: %d | Wins: %d",
		hero.AttackMin, hero.AttackMax, hero.Defense, hero.Wins), heroStyle)

	printText(2, 7, fmt.Sprintf("%s", enemy.Name), enemyStyle)
	printText(2, 8, fmt.Sprintf("HP: %d/%d %s", enemy.Health, enemy.MaxHealth, enemyHealthBar), enemyStyle)
	printText(2, 9, fmt.Sprintf("ATK: %d-%d | DEF: %d",
		enemy.AttackMin, enemy.AttackMax, enemy.Defense), enemyStyle)

	// Draw battle log with scrolling
	printText(2, 11, "BATTLE LOG:", titleStyle)

	startY := 12
	displayLog := gameState.BattleLog
	if len(gameState.BattleLog) > model.MaxLogEntries {
		displayLog = gameState.BattleLog[len(gameState.BattleLog)-model.MaxLogEntries:]
	}

	// Process battle log for better formatting
	for i, line := range displayLog {

		// Use different styles for special messages
		style := defaultStyle
		if strings.Contains(line, "CRITICAL HIT") {
			printText(2, startY+i, line, criticalStyle)
		} else if strings.Contains(line, "BLOCKED") {
			printText(2, startY+i, line, blockStyle)
		} else if strings.Contains(line, "VICTORIOUS") {
			printText(2, startY+i, line, titleStyle)
		} else {
			printText(2, startY+i, line, style)
		}
	}

	// Show controls or upgrade options
	controlsY := startY + len(displayLog) + 2

	if gameState.UpgradeMode {
		printText(2, controlsY, "CHOOSE YOUR UPGRADE:", titleStyle)
		for i, upgrade := range gameState.Upgrades {
			style := infoStyle
			prefix := "   "
			if i == gameState.SelectedUpgrade {
				style = selectedStyle
				prefix = ">> "
			}
			printText(2, controlsY+i+1, fmt.Sprintf("%s%d. %s - %s",
				prefix, i+1, upgrade.Name, upgrade.Description), style)
		}
		printText(2, controlsY+len(gameState.Upgrades)+2,
			"Use UP/DOWN arrows to select, ENTER to confirm", infoStyle)
	} else if gameState.GameOver {
		printText(2, controlsY, "Game Over! Press 'q' to exit or 'r' to start a new run.", infoStyle)
	} else {
		printText(2, controlsY, "Press 'q' to quit.", infoStyle)
	}

	screen.Show()
}
