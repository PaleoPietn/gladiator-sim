package game

import (
	"fmt"
	"math/rand"
	"time"

	model "stuff/models"
	"stuff/ui"

	"github.com/gdamore/tcell/v2"
)

// GameHandler implements the ui.InputHandler interface
type GameHandler struct{}

// TurnDelay is the delay between battle turns
const TurnDelay = 800 * time.Millisecond

// RandRange returns a random number between min and max (inclusive)
func RandRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// CalculateDamage determines attack damage with critical hits and blocks
func CalculateDamage(attacker, defender *model.Player) model.BattleResult {
	damage := RandRange(attacker.AttackMin, attacker.AttackMax)
	isCritical := rand.Intn(model.CriticalChance) == 0
	isBlocked := rand.Intn(model.BlockChance) == 0

	if isCritical {
		damage *= 2
	}
	if isBlocked {
		damage /= 2
	}

	// Apply defender's defense to reduce damage
	damage -= defender.Defense
	if damage < 1 {
		damage = 1 // Minimum damage is 1
	}

	defender.Health -= damage
	isGameOver := defender.Health <= 0

	// Ensure health doesn't go below zero for display purposes
	if defender.Health < 0 {
		defender.Health = 0
	}

	var winnerName string
	if isGameOver {
		winnerName = attacker.Name
		attacker.Wins++
	}

	return model.BattleResult{
		Attacker:   attacker,
		Defender:   defender,
		Damage:     damage,
		IsCritical: isCritical,
		IsBlocked:  isBlocked,
		IsGameOver: isGameOver,
		WinnerName: winnerName,
	}
}

// FormatBattleMessage creates a descriptive message for the battle log
func FormatBattleMessage(result model.BattleResult) string {
	msg := fmt.Sprintf("%s strikes %s for %d damage! (%s HP: %d/%d)",
		result.Attacker.Name, result.Defender.Name, result.Damage,
		result.Defender.Name, result.Defender.Health, result.Defender.MaxHealth)

	if result.IsCritical {
		msg += " 💥 CRITICAL HIT!"
	}
	if result.IsBlocked {
		msg += " 🛡️  BLOCKED!"
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

	// Call the UI package to draw the initial battle screen
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

				// Calculate damage and update health
				result := CalculateDamage(attacker, defender)
				gameState.AddToBattleLog(FormatBattleMessage(result))

				// Check for game over
				if result.IsGameOver {
					gameState.AddToBattleLog("")

					if defender.IsHero {
						// Hero lost
						gameState.AddToBattleLog(fmt.Sprintf("💀 %s has fallen! GAME OVER 💀", hero.Name))
						gameState.AddToBattleLog(fmt.Sprintf("Final Score: %d victories", hero.Wins))
						gameState.GameOver = true
					} else {
						// Hero won
						gameState.AddToBattleLog(fmt.Sprintf("🏆 %s is VICTORIOUS! 🏆", hero.Name))
						gameState.AddToBattleLog("Choose an upgrade to continue your journey!")
						gameState.UpgradeMode = true
					}

					// Update the UI with the final battle state
					ui.DrawUI(screen, hero, enemy, gameState)
					done <- true
					return
				}

				turn++
				// Update the UI after each turn
				ui.DrawUI(screen, hero, enemy, gameState)
				time.Sleep(TurnDelay)
			}
		}
	}()
}
