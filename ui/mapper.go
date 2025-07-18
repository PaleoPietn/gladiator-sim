package ui

import (
	"gladiator-sim/models"

	"github.com/gdamore/tcell/v2"
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

func mapLogType(logType models.LogType) tcell.Style {
	switch logType {
	case models.LogTypeTitle:
		return titleStyle
	case models.LogTypeHero:
		return heroStyle
	case models.LogTypeEnemy:
		return enemyStyle
	case models.LogTypeInfo:
		return infoStyle
	case models.LogTypeDefault:
		return defaultStyle
	case models.LogTypeSelected:
		return selectedStyle
	case models.LogTypeCritical:
		return criticalStyle
	case models.LogTypeBlock:
		return blockStyle
	default:
		return defaultStyle
	}
}
