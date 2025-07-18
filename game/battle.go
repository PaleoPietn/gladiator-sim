package game

import (
	"fmt"
	"math/rand"
	"time"

	model "gladiator-sim/models"
)

const (
	LogFirstBattle = "🔥GLADIATOR BATTLE🔥"

	LogHeroDeath           = "💀 %s has fallen! GAME OVER 💀"
	LogHeroDeathFinalScore = "Final Score: %d victories"

	logLastFight = "You've defeated all champions! Now face THE IMMORTAL!"

	logVictory        = "🏆 %s is VICTORIUS! 🏆"
	logPrepareUpgrade = "Choose an upgrade to continue your journey!" // TODO: move this log to UI?

	logLegendaryVictory1 = "🎉 LEGENDARY VICTORY! You've defeated The Immortal! 🎉"
	logLegendaryVictory2 = "🏆 Your name will be remembered for eternity! 🏆"

	logCritical = " 󰓥 CRITICAL HIT!"
	logBlocked  = " 󰒘 BLOCKED!"

	logStrike = "%s strikes %s for %d damage!"
)

// GameHandler implements the ui.InputHandler interface
type GameHandler struct{}

// TurnDelay is the delay between battle turns
const TurnDelay = 800 * time.Millisecond

func randRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// CalculateDamage determines attack damage with critical hits and blocks
func calculateDamage(attacker, defender *model.Player) model.BattleResult {
	damage := randRange(attacker.AttackMin, attacker.AttackMax)

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

// StartBattle handles the battle loop between two players
func (eng *GameEngine) StartBattle() {
	battleStates := []*model.BattleState{
		eng.generateBattleState(LogFirstBattle, model.LogTypeTitle, model.LogWaitTimeShort),
		eng.generateBattleState(fmt.Sprintf("%s vs %s", eng.hero.Name, eng.enemy.Name), model.LogTypeTitle, model.LogWaitTimeShort),
	}

	if eng.enemy.Description != "" {
		battleStates = append(battleStates, eng.generateBattleState(eng.enemy.Description, model.LogTypeInfo, model.LogWaitTimeLong))
	}

	turn := 0

	for eng.hero.Health > 0 && eng.enemy.Health > 0 {
		// Determine attacker and defender based on turn
		attacker, defender := eng.hero, eng.enemy
		if turn%2 == 1 {
			attacker, defender = eng.enemy, eng.hero
		}

		result := calculateDamage(attacker, defender)

		msg := fmt.Sprintf(logStrike, result.Attacker.Name, result.Defender.Name, result.Damage)
		switch {
		case result.IsCritical:
			msg += logCritical
			battleStates = append(battleStates, eng.generateBattleState(msg, model.LogTypeCritical, model.LogWaitTimeDefault))
		case result.IsBlocked:
			msg += logBlocked
			battleStates = append(battleStates, eng.generateBattleState(msg, model.LogTypeBlock, model.LogWaitTimeDefault))
		default:
			battleStates = append(battleStates, eng.generateBattleState(msg, model.LogTypeInfo, model.LogWaitTimeDefault))
		}

		if result.IsGameOver {
			if defender.IsHero {
				// Hero lost
				battleStates = append(battleStates, eng.generateBattleState(fmt.Sprintf(LogHeroDeath, eng.hero.Name), model.LogTypeInfo, model.LogWaitTimeDefault))
				battleStates = append(battleStates, eng.generateBattleState(fmt.Sprintf(LogHeroDeathFinalScore, eng.hero.Wins), model.LogTypeInfo, model.LogWaitTimeDefault))
				eng.gameState.GameOver = true
			} else {
				// Hero won
				battleStates = append(battleStates, eng.generateBattleState(fmt.Sprintf(logVictory, eng.hero.Name), model.LogTypeTitle, model.LogWaitTimeShort))

				switch {
				// Hero won the last fight
				case eng.hero.Wins >= len(enemyTypes)+1:
					battleStates = append(battleStates, eng.generateBattleState(logLegendaryVictory1, model.LogTypeTitle, model.LogWaitTimeShort))
					battleStates = append(battleStates, eng.generateBattleState(logLegendaryVictory2, model.LogTypeTitle, model.LogWaitTimeShort))
					eng.gameState.GameOver = true

					// Prepare for final battle
				case eng.hero.Wins == len(enemyTypes):
					battleStates = append(battleStates, eng.generateBattleState(logLastFight, model.LogTypeTitle, model.LogWaitTimeShort))

					// Otherwise, prepare for next battle
				default:
					battleStates = append(battleStates, eng.generateBattleState(logPrepareUpgrade, model.LogTypeTitle, model.LogWaitTimeShort))
				}
			}
		}
		turn++
	}

	eng.ui.DrawBattle(battleStates)
}
