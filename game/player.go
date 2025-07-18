package game

import model "gladiator-sim/models"

// Default stats for Player Character
var heroDefaultStats = &model.Player{
	Wins:         0,
	IsHero:       true,
	Health:       130,
	MaxHealth:    130,
	AttackMin:    10,
	AttackMax:    15,
	Defense:      1,
	CritChance:   10,
	BlockChance:  10,
	LifeSteal:    0,
	CritDamage:   0,
	Regeneration: 0,
}

// NewHero creates a new player character
func NewHero(playerName string) *model.Player {
	hero := heroDefaultStats
	hero.Name = playerName
	return hero
}

// ResetHero resets the hero to starting stats
func (eng *GameEngine) resetHero() {
	// keep old name
	name := eng.hero.Name

	// reset hero stats
	eng.hero = heroDefaultStats

	// set old name
	eng.hero.Name = name
}

// applyUpgrade applies an upgrade to the player
func (eng *GameEngine) applyUpgrade(upgrade model.Upgrade) {
	upgrade.Effect(eng.hero)
}
