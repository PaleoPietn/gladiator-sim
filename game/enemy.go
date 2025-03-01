// game/enemy.go
package game

import (
	"math/rand"
	model "stuff/models"
)

// EnemyType defines a template for creating enemies with specific characteristics
type EnemyType struct {
	Name         string
	HealthMod    float64
	AttackMod    float64
	DefenseMod   float64
	CritChance   int
	BlockChance  int
	LifeSteal    int
	CritDamage   int
	Regeneration int
	Description  string
}

// Define enemy archetypes
var enemyTypes = []EnemyType{
	{
		Name:        "Novice Gladiator",
		HealthMod:   1.0,
		AttackMod:   0.9,
		DefenseMod:  0.8,
		CritChance:  5,
		BlockChance: 5,
		Description: "A fresh recruit to the arena, eager but inexperienced.",
	},
	{
		Name:        "Veteran Fighter",
		HealthMod:   1.1,
		AttackMod:   1.0,
		DefenseMod:  1.0,
		CritChance:  8,
		BlockChance: 10,
		Description: "Seasoned by countless battles, this warrior knows the arena well.",
	},
	{
		Name:        "Arena Champion",
		HealthMod:   1.4,
		AttackMod:   1.1,
		DefenseMod:  1.2,
		CritChance:  10,
		BlockChance: 15,
		CritDamage:  20,
		Description: "A celebrated fighter who has claimed many lives in the arena.",
	},
	{
		Name:        "Blood Reaper",
		HealthMod:   0.9,
		AttackMod:   1.3,
		DefenseMod:  0.7,
		CritChance:  15,
		BlockChance: 5,
		CritDamage:  30,
		Description: "Known for swift, devastating attacks that leave opponents bleeding.",
	},
	{
		Name:        "Skull Crusher",
		HealthMod:   1.2,
		AttackMod:   1.4,
		DefenseMod:  0.8,
		CritChance:  12,
		BlockChance: 8,
		CritDamage:  40,
		Description: "Wields a massive weapon that can shatter bone with a single blow.",
	},
	{
		Name:        "Death Dealer",
		HealthMod:   1.0,
		AttackMod:   1.5,
		DefenseMod:  0.6,
		CritChance:  20,
		BlockChance: 5,
		CritDamage:  50,
		Description: "An executioner who specializes in finishing opponents quickly.",
	},
	{
		Name:        "Soul Harvester",
		HealthMod:   0.8,
		AttackMod:   1.2,
		DefenseMod:  0.5,
		CritChance:  10,
		BlockChance: 5,
		LifeSteal:   15,
		Description: "Drains the life force from opponents to sustain itself.",
	},
	{
		Name:        "Bone Breaker",
		HealthMod:   1.3,
		AttackMod:   1.3,
		DefenseMod:  1.0,
		CritChance:  15,
		BlockChance: 10,
		CritDamage:  35,
		Description: "Targets joints and weak points, causing crippling injuries.",
	},
	{
		Name:        "Doom Bringer",
		HealthMod:   1.2,
		AttackMod:   1.4,
		DefenseMod:  1.1,
		CritChance:  15,
		BlockChance: 15,
		CritDamage:  40,
		LifeSteal:   10,
		Description: "A harbinger of death whose mere presence strikes fear into opponents.",
	},
	{
		Name:        "Shadow Assassin",
		HealthMod:   0.7,
		AttackMod:   1.6,
		DefenseMod:  0.4,
		CritChance:  25,
		BlockChance: 15,
		CritDamage:  60,
		Description: "Strikes from the darkness with lethal precision.",
	},
	{
		Name:         "Hans",
		HealthMod:    1.6,
		AttackMod:    0.9,
		DefenseMod:   1.5,
		CritChance:   5,
		BlockChance:  25,
		Regeneration: 2,
		Description:  "A walking fortress clad in impenetrable armor.",
	},
	{
		Name:        "Berserker",
		HealthMod:   1.1,
		AttackMod:   1.5,
		DefenseMod:  0.3,
		CritChance:  20,
		BlockChance: 0,
		CritDamage:  50,
		Description: "Fights with reckless abandon, caring nothing for defense.",
	},
	{
		Name:         "Blood Mage",
		HealthMod:    0.9,
		AttackMod:    1.3,
		DefenseMod:   0.6,
		CritChance:   15,
		BlockChance:  10,
		LifeSteal:    20,
		Regeneration: 3,
		Description:  "Wields forbidden magic that manipulates life essence.",
	},
	{
		Name:         "Undying One",
		HealthMod:    1.3,
		AttackMod:    1.0,
		DefenseMod:   0.8,
		CritChance:   10,
		BlockChance:  10,
		Regeneration: 5,
		Description:  "A fighter who refuses to fall, healing from even grievous wounds.",
	},
	{
		Name:        "Twin Blade",
		HealthMod:   0.8,
		AttackMod:   1.4,
		DefenseMod:  0.7,
		CritChance:  18,
		BlockChance: 18,
		CritDamage:  30,
		Description: "Wields a blade in each hand, attacking with blinding speed.",
	},
}

// CreateEnemy generates a themed enemy based on the current level
func (h *GameHandler) CreateEnemy(level int) *model.Player {
	// Check if this is the final boss level
	if level == len(enemyTypes)+1 {
		baseHealth := 80 + (level * 10)
		baseAttackMin := 5 + (level * 2)
		baseAttackMax := 10 + (level * 3)
		baseDefense := level

		health := int(float64(baseHealth) * finalBoss.HealthMod)
		attackMin := int(float64(baseAttackMin) * finalBoss.AttackMod)
		attackMax := int(float64(baseAttackMax) * finalBoss.AttackMod)
		defense := int(float64(baseDefense) * finalBoss.DefenseMod)

		return &model.Player{
			Name:         finalBoss.Name,
			Health:       health,
			MaxHealth:    health,
			AttackMin:    attackMin,
			AttackMax:    attackMax,
			Defense:      defense,
			CritChance:   finalBoss.CritChance,
			BlockChance:  finalBoss.BlockChance,
			LifeSteal:    finalBoss.LifeSteal,
			CritDamage:   finalBoss.CritDamage,
			Regeneration: finalBoss.Regeneration,
			Description:  finalBoss.Description,
			IsHero:       false,
		}
	}

	baseHealth := 80 + (level * 8)
	baseAttackMin := 5 + level
	baseAttackMax := 10 + (level * 2)
	baseDefense := level / 2

	// Select enemy type based on level
	enemyIndex := level - 1

	enemyType := enemyTypes[enemyIndex]

	health := int(float64(baseHealth) * enemyType.HealthMod)
	attackMin := int(float64(baseAttackMin) * enemyType.AttackMod)
	attackMax := int(float64(baseAttackMax) * enemyType.AttackMod)
	defense := int(float64(baseDefense) * enemyType.DefenseMod)

	// Add some randomness to stats
	healthVariance := rand.Intn(11) - 5 // -5 to +5
	attackVariance := rand.Intn(3) - 1  // -1 to +1

	return &model.Player{
		Name:         enemyType.Name,
		Health:       health + healthVariance,
		MaxHealth:    health + healthVariance,
		AttackMin:    attackMin + attackVariance,
		AttackMax:    attackMax + attackVariance,
		Defense:      defense,
		CritChance:   enemyType.CritChance,
		BlockChance:  enemyType.BlockChance,
		LifeSteal:    enemyType.LifeSteal,
		CritDamage:   enemyType.CritDamage,
		Regeneration: enemyType.Regeneration,
		Description:  enemyType.Description,
		IsHero:       false,
	}
}

// Final boss enemy type
var finalBoss = EnemyType{
	Name:         "The Immortal",
	HealthMod:    2.0,
	AttackMod:    1.8,
	DefenseMod:   1.5,
	CritChance:   20,
	BlockChance:  20,
	LifeSteal:    15,
	CritDamage:   50,
	Regeneration: 3,
	Description:  "The legendary undefeated champion of the arena. None have survived his wrath.",
}
