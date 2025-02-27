package game

import (
	"math/rand"
	model "stuff/models"
	"time"
)

// All possible upgrades in the game
var allUpgrades = []model.Upgrade{
	{
		Name:        "Full Heal",
		Description: "Restore all health points",
		Effect: func(p *model.Player) {
			p.Health = p.MaxHealth
		},
	},
	{
		Name:        "Strength Training",
		Description: "Increase minimum and maximum damage by 5",
		Effect: func(p *model.Player) {
			p.AttackMin += 5
			p.AttackMax += 5
		},
	},
	{
		Name:        "Defensive Stance",
		Description: "Gain 5 defense points (reduces incoming damage)",
		Effect: func(p *model.Player) {
			p.Defense += 5
		},
	},
	{
		Name:        "Vitality",
		Description: "Increase maximum health by 25",
		Effect: func(p *model.Player) {
			p.MaxHealth += 25
			p.Health += 25
		},
	},
	{
		Name:        "Critical Eye",
		Description: "Increase critical hit chance by 10%",
		Effect: func(p *model.Player) {
			p.CritChance += 10
		},
	},
	{
		Name:        "Vampiric Strike",
		Description: "Heal for 15% of damage dealt",
		Effect: func(p *model.Player) {
			p.LifeSteal += 15
		},
	},
	{
		Name:        "Berserker",
		Description: "Gain +10 max damage but -10 health",
		Effect: func(p *model.Player) {
			p.AttackMax += 10
			p.Health -= 10
			if p.Health <= 0 {
				p.Health = 1 // Prevent death from upgrade
			}
		},
	},
	{
		Name:        "Precision",
		Description: "Increase minimum damage by 8",
		Effect: func(p *model.Player) {
			p.AttackMin += 8
			if p.AttackMax < p.AttackMin {
				p.AttackMax = p.AttackMin
			}
		},
	},
	{
		Name:        "Block Master",
		Description: "Increase block chance by 10%",
		Effect: func(p *model.Player) {
			p.BlockChance += 10
		},
	},
	{
		Name:        "Executioner",
		Description: "Critical hits deal 50% more damage",
		Effect: func(p *model.Player) {
			p.CritDamage += 50
		},
	},
	{
		Name:        "Second Wind",
		Description: "Heal 5% of max health each turn",
		Effect: func(p *model.Player) {
			p.Regeneration += 5
		},
	},
	{
		Name:        "Balanced Training",
		Description: "Gain +3 to all stats (ATK, DEF, HP)",
		Effect: func(p *model.Player) {
			p.AttackMin += 3
			p.AttackMax += 3
			p.Defense += 3
			p.MaxHealth += 5
			p.Health += 5
		},
	},
}

// CreateUpgrades selects 3 random upgrades from the pool
func CreateUpgrades() []model.Upgrade {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Make a copy of all upgrades to avoid modifying the original
	availableUpgrades := make([]model.Upgrade, len(allUpgrades))
	copy(availableUpgrades, allUpgrades)

	// Shuffle the upgrades using the new random generator
	random.Shuffle(len(availableUpgrades), func(i, j int) {
		availableUpgrades[i], availableUpgrades[j] = availableUpgrades[j], availableUpgrades[i]
	})

	// Select the first 3 (or fewer if we don't have enough)
	numUpgrades := min(len(availableUpgrades), 3)

	return availableUpgrades[:numUpgrades]
}
