package game

import model "stuff/models"

// NewHero creates a new player character
func NewHero() *model.Player {
	return &model.Player{
		Name:      "Hero",
		Health:    130,
		MaxHealth: 130,
		AttackMin: 10,
		AttackMax: 15,
		Defense:   1,
		Wins:      0,
		IsHero:    true,
	}
}

// ResetHero resets the hero to starting stats
func (h *GameHandler) ResetHero(hero *model.Player) {
	hero.Health = 130
	hero.MaxHealth = 130
	hero.AttackMin = 10
	hero.AttackMax = 15
	hero.Defense = 1
	hero.Wins = 0
}

// HandleUpgrade applies an upgrade to the player
func (h *GameHandler) HandleUpgrade(hero *model.Player, upgrade model.Upgrade) {
	upgrade.Effect(hero)
}
