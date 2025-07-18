package game

import (
	"gladiator-sim/models"
	"gladiator-sim/ui"
	"time"
)

// TODO replace ui with interface
type GameEngine struct {
	hero      *models.Player
	enemy     *models.Player
	gameState *models.GameState
	ui        *ui.UI
	quit      chan struct{}
}

func NewEngine(ui *ui.UI, quit chan struct{}) *GameEngine {
	heroName := ui.StartScreen()

	go func() {
		<-quit

		ui.DrawQuitScreen()
		time.Sleep(time.Second)

	}()

	return &GameEngine{
		hero: NewHero(heroName),
		ui:   ui,
		quit: quit,
	}
}

func (eng *GameEngine) StartGameLoop() {
	for {
		eng.resetGame()
		for !eng.gameState.GameOver {
			eng.StartBattle()

			if eng.gameState.GameOver {
				break
			}

			// Upgrades
			upgrades := eng.generateUpgrades()
			idxChosen := eng.ui.ChooseUpgrade(upgrades)
			eng.applyUpgrade(upgrades[idxChosen])

			// Prep new Enemy
			eng.gameState.CurrentEnemy++
			eng.enemy = CreateEnemy(eng.gameState.CurrentEnemy)
		}
		if eng.ui.DrawEndGameScreen() {
			eng.quit <- struct{}{}
			return
		}
	}
}

// Resets hero, enemy and gameState (but not hero name)
func (eng *GameEngine) resetGame() {
	eng.resetHero()
	eng.enemy = CreateEnemy(1)
	eng.gameState = NewGameState()
}
