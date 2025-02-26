package game

import (
	"math/rand"
	model "stuff/models"
)

// CreateEnemy generates progressively harder enemies
func (h *GameHandler) CreateEnemy(level int) *model.Player {
	// Base stats that scale with level
	baseHealth := 80 + (level * 8)
	baseAttackMin := 5 + level
	baseAttackMax := 10 + (level * 2)
	baseDefense := level / 2

	// Names for enemies
	names := []string{
		"Novice Gladiator", "Veteran Fighter", "Arena Champion",
		"Blood Reaper", "Skull Crusher", "Death Dealer",
		"Soul Harvester", "Bone Breaker", "Doom Bringer",
	}

	// Add some randomness to stats
	healthVariance := rand.Intn(21) - 10 // -10 to +10
	attackVariance := rand.Intn(3) - 1   // -1 to +1

	nameIndex := level - 1
	if nameIndex >= len(names) {
		nameIndex = len(names) - 1
	}

	return &model.Player{
		Name:      names[nameIndex],
		Health:    baseHealth + healthVariance,
		MaxHealth: baseHealth + healthVariance,
		AttackMin: baseAttackMin + attackVariance,
		AttackMax: baseAttackMax + attackVariance,
		Defense:   baseDefense,
		IsHero:    false,
	}
}
