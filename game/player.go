package game

import model "stuff/models"

// NewHero creates a new player character
func NewHero() *model.Player {
	return &model.Player{
		Name:         "Hero",
		Health:       130,
		MaxHealth:    130,
		AttackMin:    10,
		AttackMax:    15,
		Defense:      1,
		Wins:         0,
		IsHero:       true,
		CritChance:   10,
		BlockChance:  10,
		LifeSteal:    0,
		CritDamage:   0,
		Regeneration: 0,
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
	hero.CritChance = 10
	hero.BlockChance = 10
	hero.LifeSteal = 0
	hero.CritDamage = 0
	hero.Regeneration = 0
}

// HandleUpgrade applies an upgrade to the player
func (h *GameHandler) HandleUpgrade(hero *model.Player, upgrade model.Upgrade) {
	upgrade.Effect(hero)
}
