package game

import (
	model "gladiator-sim/models"
	"time"
)

// NewGameState creates a new game state
func NewGameState() *model.GameState {
	return &model.GameState{
		CurrentEnemy: 1,
		GameOver:     false,
	}
}

// ResetGameState resets the game state to starting values
func (h *GameHandler) ResetGameState(state *model.GameState) {
	state.CurrentEnemy = 1
	state.GameOver = false

	ResetUpgradeTracker()
}

// Generates a BattleState based on the log details passed *current hero & enemy state
func (eng *GameEngine) generateBattleState(logMessage string, logType model.LogType, logWaitTime time.Duration) *model.BattleState {
	return &model.BattleState{
		Hero: &model.Player{
			Name:         eng.hero.Name,
			Health:       eng.hero.Health,
			MaxHealth:    eng.hero.MaxHealth,
			AttackMin:    eng.hero.AttackMin,
			AttackMax:    eng.hero.AttackMax,
			Defense:      eng.hero.Defense,
			Wins:         eng.hero.Wins,
			IsHero:       true,
			Regeneration: eng.hero.Regeneration,
		},
		Enemy: &model.Player{
			Name:      eng.enemy.Name,
			Health:    eng.enemy.Health,
			MaxHealth: eng.enemy.MaxHealth,
			AttackMin: eng.enemy.AttackMin,
			AttackMax: eng.enemy.AttackMax,
			Defense:   eng.enemy.Defense,
			Wins:      eng.enemy.Wins,
			IsHero:    false, // not mandatory but cleaner/avoids confusion
		},
		BattleLog: &model.BattleLog{
			LogMessage:  logMessage,
			LogType:     logType,
			LogWaitTime: logWaitTime,
		},
	}
}
