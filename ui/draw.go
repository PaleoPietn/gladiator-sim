package ui

import (
	"fmt"
	"strings"

	model "gladiator-sim/models"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

var (
	// Define styles
	defaultStyle  = tcell.StyleDefault
	titleStyle    = defaultStyle.Bold(true).Foreground(tcell.ColorYellow)
	heroStyle     = defaultStyle.Foreground(tcell.ColorGreen)
	enemyStyle    = defaultStyle.Foreground(tcell.ColorRed)
	infoStyle     = defaultStyle.Foreground(tcell.ColorWhite)
	selectedStyle = defaultStyle.Background(tcell.ColorGainsboro).Foreground(tcell.ColorWhite)
	criticalStyle = defaultStyle.Foreground(tcell.ColorYellow)
	blockStyle    = defaultStyle.Foreground(tcell.ColorTeal)
)

const (
	// UI
	healthBarWidth = 20
	logStartY      = 8
	maxLogEntries  = 20
	heroXIndex     = 2
	enemyXIndex    = 52
	playerYIndex   = 3

	// Texts
	titleText        = "ROGUELIKE GLADIATOR ARENA"
	battleLogText    = "BATTLE LOG:"
	upgradeText      = "CHOOSE YOUR UPGRADE:"
	prefixUnselected = "   "
	prefixSelected   = ">> "
	upgradeHelper    = "Use UP/DOWN arrows to select, ENTER to confirm"
	gameOverHelper   = "Game Over! Press 'q' to exit or 'r' to start a new run."
	quitHelper       = "Press 'q' to quit."
)

// drawHealthBar creates a visual health bar
func drawHealthBar(current, max int, width int) string {
	if max <= 0 {
		return "[ERROR]" // TODO: return an actual error?
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
	screen.Clear()
	defer screen.Show()

	// Draw title and stats
	printText(screen, 2, 1, titleText, titleStyle)

	// Draw players (hero & enemy) stats with health bars
	drawPlayer(screen, hero)
	drawPlayer(screen, enemy)

	// Draw battle log with scrolling
	printText(screen, 2, 11, battleLogText, titleStyle)

	// TODO: refactor BattleLog to be its own type that includes the string + the status (crit, block, victory...)
	// so we can avoid strings.Contains for styling
	displayLog := gameState.BattleLog
	// Pop the top logs if the length of the battle is too long to display everything
	if len(gameState.BattleLog) > maxLogEntries {
		displayLog = gameState.BattleLog[len(gameState.BattleLog)-maxLogEntries:]
	}

	// Style the log if special (crit, block, victory)
	for i, line := range displayLog {
		style := defaultStyle
		switch {
		case strings.Contains(line, model.CriticalHit):
			printText(screen, 2, logStartY+i, line, criticalStyle)
		case strings.Contains(line, model.Blocked):
			printText(screen, 2, logStartY+i, line, blockStyle)
		case strings.Contains(line, model.Victorious):
			printText(screen, 2, logStartY+i, line, titleStyle)
		default:
			printText(screen, 2, logStartY+i, line, style)
		}
	}

	// Show controls or upgrade options
	controlsY := logStartY + len(displayLog) + 2

	switch {
	case gameState.UpgradeMode:
		printText(screen, 2, controlsY, upgradeText, titleStyle)
		for i, upgrade := range gameState.Upgrades {
			style := infoStyle

			prefix := prefixUnselected
			if i == gameState.SelectedUpgrade {
				style = selectedStyle
				prefix = prefixSelected
			}
			printText(screen, 2, controlsY+i+1, fmt.Sprintf("%s%d. %s - %s", prefix, i+1, upgrade.Name, upgrade.Description), style)
		}
		printText(screen, 2, controlsY+len(gameState.Upgrades)+2, upgradeHelper, infoStyle)
	case gameState.GameOver:
		printText(screen, 2, controlsY, gameOverHelper, infoStyle)
	default:
		printText(screen, 2, controlsY, quitHelper, infoStyle)
	}
}

func generateBuffsString(player *model.Player) string {
	// TODO: is the check necessary? Maybe we want to print buffs for enemies too?
	if !player.IsHero {
		return ""
	}

	// TODO: create an array of buffs that were selected to print them here
	buffs := ""
	if player.Regeneration > 0 {
		buffs += " 🌿"
	}

	return buffs
}

// Helper function to draw text with proper handling of wide characters
// TODO: refactor to pass only screen + an object that contains the other characters
func printText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	posX := x
	for _, c := range text {
		// Check if character is an emoji or other wide character
		width := runewidth.RuneWidth(c)
		screen.SetContent(posX, y, c, nil, style)
		posX += width
	}
}

func drawPlayer(screen tcell.Screen, player *model.Player) {
	style := heroStyle
	xIndex := heroXIndex

	startYIndex := playerYIndex
	if !player.IsHero {
		xIndex = enemyXIndex
		style = enemyStyle
	}

	healthBar := drawHealthBar(player.Health, player.MaxHealth, healthBarWidth)
	printText(screen, xIndex, startYIndex, fmt.Sprintf("%s %s", player.Name, generateBuffsString(player)), style)
	startYIndex++
	printText(screen, xIndex, startYIndex, fmt.Sprintf("%s %s", formatLifeCount(player.Health, player.MaxHealth), healthBar), style)
	startYIndex++
	switch {
	case player.IsHero:
		printText(screen, xIndex, startYIndex, fmt.Sprintf("ATK: %d-%d | DEF: %d | Wins: %d", player.AttackMin, player.AttackMax, player.Defense, player.Wins), style)
	default:
		printText(screen, xIndex, startYIndex, fmt.Sprintf("ATK: %d-%d | DEF: %d", player.AttackMin, player.AttackMax, player.Defense), style)
	}
}

func formatLifeCount(health, maxHealth int) string {
	minWidth := 7
	return fmt.Sprintf("HP: %*s", minWidth, fmt.Sprintf("%d/%d", health, maxHealth))
}
