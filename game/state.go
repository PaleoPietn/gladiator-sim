package game

import "stuff/models"

// NewGameState creates a new game state
func NewGameState() *model.GameState {
	return &model.GameState{
		CurrentEnemy:    1,
		UpgradeMode:     false,
		Upgrades:        CreateUpgrades(),
		SelectedUpgrade: 0,
		BattleLog:       []string{},
		GameOver:        false,
	}
}

// ResetGameState resets the game state for a new run
func (h *GameHandler) ResetGameState(gs *model.GameState) {
	gs.CurrentEnemy = 1
	gs.UpgradeMode = false
	gs.SelectedUpgrade = 0
	gs.BattleLog = []string{}
	gs.GameOver = false
}
