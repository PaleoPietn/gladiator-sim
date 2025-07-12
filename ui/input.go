package ui

import (
	model "gladiator-sim/models"

	"github.com/gdamore/tcell/v2"
)

// InputHandler defines the interface for handling game input
type InputHandler interface {
	HandleUpgrade(hero *model.Player, upgrade model.Upgrade)
	CreateEnemy(level int) *model.Player
	ResetHero(hero *model.Player)
	ResetGameState(state *model.GameState)
	StartBattle(hero, enemy *model.Player, screen tcell.Screen, state *model.GameState, quit, done chan bool)
}

// StartInputHandler initializes the input handling goroutine
func StartInputHandler(
	screen tcell.Screen,
	hero *model.Player,
	enemy *model.Player,
	gameState *model.GameState,
	handler InputHandler,
	quit chan bool,
	done chan bool) {

	go func() {
		for {
			ev := screen.PollEvent()
			shouldQuit := HandleInput(ev, screen, hero, enemy, gameState, handler, quit, done)
			if shouldQuit {
				return
			}
		}
	}()
}

// HandleInput processes a single input event
func HandleInput(
	ev tcell.Event,
	screen tcell.Screen,
	hero *model.Player,
	enemy *model.Player,
	gameState *model.GameState,
	handler InputHandler,
	quit chan bool,
	done chan bool) bool {

	switch ev := ev.(type) {
	case *tcell.EventKey:
		if gameState.UpgradeMode {
			return handleUpgradeInput(ev, screen, hero, enemy, gameState, handler, quit, done)
		} else {
			return handleRegularInput(ev, screen, hero, gameState, handler, quit, done)
		}
	case *tcell.EventResize:
		screen.Sync()
		DrawUI(screen, hero, enemy, gameState)
		return false
	}
	return false
}

// handleUpgradeInput processes input during upgrade selection
func handleUpgradeInput(ev *tcell.EventKey,
	screen tcell.Screen,
	hero *model.Player,
	enemy *model.Player,
	gameState *model.GameState,
	handler InputHandler,
	quit chan bool,
	done chan bool) bool {

	switch ev.Key() {
	case tcell.KeyUp:
		gameState.SelectedUpgrade = (gameState.SelectedUpgrade - 1 + len(gameState.Upgrades)) % len(gameState.Upgrades)
		DrawUI(screen, hero, enemy, gameState)
		return false
	case tcell.KeyDown:
		gameState.SelectedUpgrade = (gameState.SelectedUpgrade + 1) % len(gameState.Upgrades)
		DrawUI(screen, hero, enemy, gameState)
		return false
	case tcell.KeyEnter:
		handler.HandleUpgrade(hero, gameState.Upgrades[gameState.SelectedUpgrade])

		// Prepare for next battle
		gameState.CurrentEnemy++
		gameState.UpgradeMode = false
		newEnemy := handler.CreateEnemy(gameState.CurrentEnemy)

		gameState.AddToBattleLog(
			"Upgrade chosen: " + gameState.Upgrades[gameState.SelectedUpgrade].Name)
		gameState.AddToBattleLog(
			"Preparing for battle against " + newEnemy.Name + "...")

		handler.StartBattle(hero, newEnemy, screen, gameState, quit, done)
		return false
	}
	return false
}

// handleRegularInput processes input during normal gameplay
func handleRegularInput(ev *tcell.EventKey, screen tcell.Screen, hero *model.Player, gameState *model.GameState,
	handler InputHandler, quit chan bool, done chan bool) bool {

	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'q':
			quit <- true
			return true
		case 'r':
			// Only allow restart after game over
			if gameState.GameOver {
				handler.ResetHero(hero)
				handler.ResetGameState(gameState)

				newEnemy := handler.CreateEnemy(gameState.CurrentEnemy)

				gameState.AddToBattleLog("Starting a new adventure...")
				handler.StartBattle(hero, newEnemy, screen, gameState, quit, done)
			}
			return false
		}
	} else if ev.Key() == tcell.KeyEscape {
		quit <- true
		return true
	}
	return false
}
