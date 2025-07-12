package game

import (
	"fmt"
	"math/rand"
	"time"

	model "gladiator-sim/models"
	"gladiator-sim/ui"

	"github.com/gdamore/tcell/v2"
)

// GameHandler implements the ui.InputHandler interface
type GameHandler struct{}

// TurnDelay is the delay between battle turns
const TurnDelay = 800 * time.Millisecond

func RandRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// CalculateDamage determines attack damage with critical hits and blocks
func CalculateDamage(attacker, defender *model.Player) model.BattleResult {
	damage := RandRange(attacker.AttackMin, attacker.AttackMax)

	critChance := model.CriticalChance
	if attacker.CritChance > 0 {
		critChance = attacker.CritChance
	}

	blockChance := model.BlockChance
	if defender.BlockChance > 0 {
		blockChance = defender.BlockChance
	}

	isCritical := rand.Intn(100) < critChance
	isBlocked := rand.Intn(100) < blockChance

	if isCritical {
		critMultiplier := 2.0
		if attacker.CritDamage > 0 {
			critMultiplier = 2.0 + (float64(attacker.CritDamage) / 100.0)
		}
		damage = int(float64(damage) * critMultiplier)
	}

	if isBlocked {
		damage /= 2
	}

	damage -= defender.Defense
	if damage < 1 {
		damage = 1
	}

	defender.Health -= damage

	if attacker.LifeSteal > 0 {
		healAmount := int(float64(damage) * float64(attacker.LifeSteal) / 100.0)
		if healAmount > 0 {
			attacker.Health += healAmount
			if attacker.Health > attacker.MaxHealth {
				attacker.Health = attacker.MaxHealth
			}
		}
	}

	regen := 0
	if defender.Regeneration > 0 {
		healAmount := int(float64(defender.MaxHealth) * float64(defender.Regeneration) / 100.0)
		regen = healAmount
		if healAmount > 0 {
			defender.Health += healAmount
			if defender.Health > defender.MaxHealth {
				defender.Health = defender.MaxHealth
			}
		}
	}

	isGameOver := defender.Health <= 0

	if defender.Health < 0 {
		defender.Health = 0
	}

	var winnerName string
	if isGameOver {
		winnerName = attacker.Name
		attacker.Wins++
		if attacker.LifeOnKill > 0 {
			attacker.Health += attacker.LifeOnKill
			if attacker.Health > attacker.MaxHealth {
				attacker.Health = attacker.MaxHealth
			}
		}
	}

	return model.BattleResult{
		Attacker:     attacker,
		Defender:     defender,
		Damage:       damage,
		IsCritical:   isCritical,
		IsBlocked:    isBlocked,
		IsGameOver:   isGameOver,
		WinnerName:   winnerName,
		Regeneration: regen,
	}
}

// FormatBattleMessage creates a descriptive message for the battle log
func FormatBattleMessage(result model.BattleResult) string {
	msg := fmt.Sprintf("%s strikes %s for %d damage!",
		result.Attacker.Name, result.Defender.Name, result.Damage)

	if result.IsCritical {
		msg += fmt.Sprintf(" 󰓥 %s!", model.CriticalHit)
	}
	if result.IsBlocked {
		msg += fmt.Sprintf(" 󰒘 %s!", model.Blocked)
	}
	return msg
}

// StartBattle handles the battle loop between two players
func (h *GameHandler) StartBattle(hero, enemy *model.Player, screen tcell.Screen, gameState *model.GameState, quit chan bool, done chan bool) {
	gameState.BattleLog = []string{
		"🔥GLADIATOR BATTLE🔥",
		fmt.Sprintf("%s vs %s", hero.Name, enemy.Name),
		"",
	}

	if enemy.Description != "" {
		gameState.BattleLog = append(gameState.BattleLog, enemy.Description)
	}

	gameState.BattleLog = append(gameState.BattleLog, "")

	ui.DrawUI(screen, hero, enemy, gameState)
	time.Sleep(TurnDelay)

	go func() {
		turn := 0

		for hero.Health > 0 && enemy.Health > 0 {
			select {
			case <-quit:
				return // Exit if user presses 'q'
			default:
				// Determine attacker and defender based on turn
				attacker, defender := hero, enemy
				if turn%2 == 1 {
					attacker, defender = enemy, hero
				}

				result := CalculateDamage(attacker, defender)
				gameState.AddToBattleLog(FormatBattleMessage(result))

				if result.IsGameOver {
					gameState.AddToBattleLog("")

					if defender.IsHero {
						// Hero lost
						gameState.AddToBattleLog(fmt.Sprintf("💀 %s has fallen! GAME OVER 💀", hero.Name))
						gameState.AddToBattleLog(fmt.Sprintf("Final Score: %d victories", hero.Wins))
						gameState.GameOver = true
					} else {
						// Hero won
						gameState.AddToBattleLog(fmt.Sprintf("🏆 %s is %s! 🏆", hero.Name, model.Victorious))

						if hero.Wins >= len(enemyTypes)+1 {
							gameState.AddToBattleLog("🎉 LEGENDARY VICTORY! You've defeated The Immortal! 🎉")
							gameState.AddToBattleLog("🏆 Your name will be remembered for eternity! 🏆")
							gameState.GameOver = true

							// Prepare for final battle
						} else if hero.Wins == len(enemyTypes) {
							gameState.AddToBattleLog("You've defeated all champions! Now face THE IMMORTAL!")
							gameState.UpgradeMode = true
							gameState.Upgrades = CreateUpgrades(hero)

							// Otherwise, prepare for next battle
						} else {
							gameState.AddToBattleLog("Choose an upgrade to continue your journey!")
							gameState.UpgradeMode = true
							gameState.Upgrades = CreateUpgrades(hero)
						}
					}

					ui.DrawUI(screen, hero, enemy, gameState)
					done <- true
					return
				}

				turn++
				ui.DrawUI(screen, hero, enemy, gameState)
				time.Sleep(TurnDelay)
			}
		}
	}()
}
