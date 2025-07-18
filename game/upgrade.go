package game

import (
	"fmt"
	model "gladiator-sim/models"
	"math/rand"
	"slices"
)

// UpgradeType defines an upgrade with its properties and limitations
type UpgradeType struct {
	Name        string
	Description string
	Effect      func(p *model.Player)
	MaxLevel    int // Maximum times this upgrade can be chosen
	Rarity      int // Higher rarity means less common (1-3)
	IsAvailable func(p *model.Player) bool
}

// All possible upgrades in the game
var allUpgrades = []UpgradeType{
	{
		Name:        "Full Heal",
		Description: "Restore all health points",
		Effect: func(p *model.Player) {
			p.Health = p.MaxHealth
		},
		MaxLevel: 999, // unlimited uses
		Rarity:   1,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	{
		Name:        "Strength Training",
		Description: "Increase minimum and maximum damage by 8",
		Effect: func(p *model.Player) {
			p.AttackMin += 8
			p.AttackMax += 8
		},
		MaxLevel: 3,
		Rarity:   1,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	{
		Name:        "Advanced Strength Training",
		Description: "Increase minimum and maximum damage by 15",
		Effect: func(p *model.Player) {
			p.AttackMin += 15
			p.AttackMax += 15
		},
		MaxLevel: 1,
		Rarity:   3,
		IsAvailable: func(p *model.Player) bool {
			return GetUpgradeLevel("Strength Training") >= 2
		},
	},
	{
		Name:        "Defensive Stance",
		Description: "Gain 5 defense points",
		Effect: func(p *model.Player) {
			p.Defense += 5
		},
		MaxLevel: 5,
		Rarity:   1,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	{
		Name:        "Iron Skin",
		Description: "Gain 12 defense points and +10% block chance",
		Effect: func(p *model.Player) {
			p.Defense += 12
			p.BlockChance += 10
		},
		MaxLevel: 1,
		Rarity:   3,
		IsAvailable: func(p *model.Player) bool {
			return GetUpgradeLevel("Defense Stance") >= 2
		},
	},
	{

		Name:        "Vitality",
		Description: "Increase maximum health by 40",
		Effect: func(p *model.Player) {
			p.MaxHealth += 40
			p.Health += 40
		},
		MaxLevel: 4,
		Rarity:   1,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	{
		Name:        "Critical Eye",
		Description: "Increase critical hit chance by 12%",
		Effect: func(p *model.Player) {
			p.CritChance += 12
		},
		MaxLevel: 5,
		Rarity:   1,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	{
		Name:        "Vampiric Strike",
		Description: "Heal for 20% of damage dealt",
		Effect: func(p *model.Player) {
			p.LifeSteal += 20
		},
		MaxLevel: 3,
		Rarity:   2,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	{
		Name:        "Blood Frenzy",
		Description: "Increase lifesteal by 15% and gain +10 attack",
		Effect: func(p *model.Player) {
			p.LifeSteal += 15
			p.AttackMin += 10
			p.AttackMax += 10
		},
		MaxLevel: 2,
		Rarity:   2,
		IsAvailable: func(p *model.Player) bool {
			return GetUpgradeLevel("Vampiric Strike") >= 1
		},
	},
	{
		Name:        "Berserker",
		Description: "Gain +25 max damage but -15 health",
		Effect: func(p *model.Player) {
			p.AttackMax += 25
			p.Health -= 15
			if p.Health <= 0 {
				p.Health = 1 // Prevent death from upgrade
			}
		},
		MaxLevel: 3,
		Rarity:   2,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	{
		Name:        "Precision",
		Description: "Increase minimum damage to 85% of maximum damage",
		Effect: func(p *model.Player) {
			p.AttackMin = int(float64(p.AttackMax) * 0.85)
		},
		IsAvailable: func(p *model.Player) bool {
			return p.AttackMin < int(float64(p.AttackMax)*0.7) // Only if there's a significant difference
		},
		MaxLevel: 5,
		Rarity:   2,
	},
	{
		Name:        "Block Master",
		Description: "Increase block chance by 15%",
		Effect: func(p *model.Player) {
			p.BlockChance += 15
		},
		MaxLevel: 5,
		Rarity:   1,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	{
		Name:        "Executioner",
		Description: "Critical hits deal 75% more damage",
		Effect: func(p *model.Player) {
			p.CritDamage += 75
		},
		MaxLevel: 3,
		Rarity:   2,
		IsAvailable: func(p *model.Player) bool {
			return p.CritChance > 10
		},
	},
	{
		Name:        "Deathblow",
		Description: "Critical hits deal 100% more damage and +5% crit chance",
		Effect: func(p *model.Player) {
			p.CritDamage += 100
			p.CritChance += 5
		},
		MaxLevel: 1,
		Rarity:   3,
		IsAvailable: func(p *model.Player) bool {
			return GetUpgradeLevel("Executioner") >= 1 && p.CritChance >= 20
		},
	},
	{
		Name:        "Second Wind",
		Description: "Heal 15% of max health each turn",
		Effect: func(p *model.Player) {
			p.Regeneration += 10
		},
		MaxLevel: 2,
		Rarity:   3,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	{
		Name:        "Battle Meditation",
		Description: "Heals for 30% after each kill",
		Effect: func(p *model.Player) {
			p.LifeOnKill += 30
		},
		MaxLevel: 1,
		Rarity:   2,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	{
		Name:        "Balanced Training",
		Description: "Gain a bit of all stats",
		Effect: func(p *model.Player) {
			p.AttackMin += 5
			p.AttackMax += 5
			p.Defense += 4
			p.MaxHealth += 20
			p.Health += 20
		},
		MaxLevel: 5,
		Rarity:   1,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	{
		Name:        "I'm Feeling Lucky",
		Description: "Gain random stat boost to random stat",
		Effect: func(p *model.Player) {
			randomStat := rand.Intn(9)
			randomAmount := rand.Intn(10) + 1
			switch randomStat {
			case 0:
				p.AttackMin += randomAmount
			case 1:
				p.AttackMax += randomAmount
			case 2:
				p.Defense += randomAmount
			case 3:
				p.MaxHealth += randomAmount
				p.Health += randomAmount
			case 4:
				p.CritChance += randomAmount
			case 5:
				p.BlockChance += randomAmount
			case 6:
				p.LifeSteal += randomAmount
			case 7:
				p.CritDamage += randomAmount
			case 8:
				p.Regeneration += randomAmount
			}
		},
		MaxLevel: 999,
		Rarity:   3,
		IsAvailable: func(p *model.Player) bool {
			return true
		},
	},
	// Other Ideas
	// Adrenalin rush -> more damage at low HP
	// ALLES ODER NIX -> 25% chance to deal triple damage, 25% chance to deal no damage
	// Adaptive armor -> gain 1 defense after each hit (max 10), resets after each battle
	// Battle Trance -> after 3 turns, gain 50% crit chance for 1 turn
	// Battle Trance 2 -> after 3 turns, gain +10 attack until end of battle
	// Epic Die Move -> once per game, revieve with 50% health after dying
	// Spiked Armor -> reflect 10% of damage taken back to attacker
	// Spiked Shiled -> deal damage on block
	// First Strike -> 25% chance to attack twice on first turn
}

// generateUpgrades generates a list of possible upgrades for the player to choose from
func (eng *GameEngine) generateUpgrades() []model.Upgrade {
	availableUpgrades := []UpgradeType{}

	for _, upgrade := range allUpgrades {
		currentLevel := GetUpgradeLevel(upgrade.Name)

		if currentLevel < upgrade.MaxLevel && upgrade.IsAvailable(eng.hero) {
			availableUpgrades = append(availableUpgrades, upgrade)
		}
	}

	selectedUpgrades := selectUpgradesByRarity(availableUpgrades, 3)

	result := []model.Upgrade{}
	for _, upgrade := range selectedUpgrades {
		// Create a copy of the upgrade to avoid closure issues
		upgradeName := upgrade.Name
		upgradeEffect := upgrade.Effect

		result = append(result, model.Upgrade{
			Name:        upgrade.Name,
			Description: upgrade.Description + getUpgradeLevelText(upgradeName),
			Effect: func(p *model.Player) {
				upgradeEffect(p)
				IncrementUpgradeLevel(upgradeName)
			},
		})
	}

	return result
}

// getUpgradeLevelText returns text showing current/max level of an upgrade
func getUpgradeLevelText(upgradeName string) string {
	for _, upgrade := range allUpgrades {
		if upgrade.Name == upgradeName {
			currentLevel := GetUpgradeLevel(upgradeName)
			if upgrade.MaxLevel < 999 {
				return fmt.Sprintf(" (%d/%d)", currentLevel, upgrade.MaxLevel)
			}
			return ""
		}
	}
	return ""
}

// selectUpgradesByRarity selects n upgrades with weighted randomness based on rarity
func selectUpgradesByRarity(upgrades []UpgradeType, n int) []UpgradeType {
	if len(upgrades) <= n {
		return upgrades
	}

	rand.Shuffle(len(upgrades), func(i, j int) {
		upgrades[i], upgrades[j] = upgrades[j], upgrades[i]
	})

	weights := make([]int, len(upgrades))
	totalWeight := 0

	for i, upgrade := range upgrades {
		// Rarity 1: weight 6, Rarity 2: weight 4, Rarity 3: weight 2
		weight := 8 - upgrade.Rarity*2
		weights[i] = weight
		totalWeight += weight
	}

	selected := []UpgradeType{}
	for len(selected) < n && len(upgrades) > 0 {
		r := rand.Intn(totalWeight)

		cumulativeWeight := 0
		for i, weight := range weights {
			cumulativeWeight += weight
			if r < cumulativeWeight {
				selected = append(selected, upgrades[i])

				totalWeight -= weights[i]
				upgrades = slices.Delete(upgrades, i, i+1)
				weights = slices.Delete(weights, i, i+1)
				break
			}
		}
	}

	return selected
}

// Global map to track upgrade levels
var upgradeTracker = make(map[string]int)

// ResetUpgradeTracker resets all upgrade levels to 0
func ResetUpgradeTracker() {
	upgradeTracker = make(map[string]int)
}

// GetUpgradeLevel returns the current level of an upgrade
func GetUpgradeLevel(upgradeName string) int {
	return upgradeTracker[upgradeName]
}

// IncrementUpgradeLevel increases the level of an upgrade by 1
func IncrementUpgradeLevel(upgradeName string) {
	upgradeTracker[upgradeName]++
}
