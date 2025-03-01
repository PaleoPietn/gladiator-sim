package game

import model "stuff/models"

// NewGameState creates a new game state
func NewGameState() *model.GameState {
	return &model.GameState{
		CurrentEnemy:    1,
		UpgradeMode:     false,
		SelectedUpgrade: 0,
		BattleLog:       []string{},
		GameOver:        false,
	}
}

// ResetGameState resets the game state to starting values
func (h *GameHandler) ResetGameState(state *model.GameState) {
	state.CurrentEnemy = 1
	state.GameOver = false
	state.UpgradeMode = false
	state.BattleLog = []string{}
	state.SelectedUpgrade = 0

	ResetUpgradeTracker()
}
