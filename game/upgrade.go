package game

import model "stuff/models"

// CreateUpgrades defines the available upgrades after each battle
func CreateUpgrades() []model.Upgrade {
	return []model.Upgrade{
		{
			Name:        "Full Heal",
			Description: "Restore all health points",
			Effect: func(p *model.Player) {
				p.Health = p.MaxHealth
			},
		},
		{
			Name:        "Strength Training",
			Description: "Increase minimum and maximum damage by 3",
			Effect: func(p *model.Player) {
				p.AttackMin += 3
				p.AttackMax += 3
			},
		},
		{
			Name:        "Defensive Stance",
			Description: "Gain 3 defense point (reduces incoming damage)",
			Effect: func(p *model.Player) {
				p.Defense += 3
			},
		},
	}
}
