package main

import (
	"fmt"
	"gladiator-sim/game"
	"gladiator-sim/ui"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Start the UI
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("Error creating screen:", err)
		return
	}
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		fmt.Println("Error initializing screen:", err)
		return
	}
	screen.Clear()

	// Start the game
	playerName := ui.ShowStartScreen(screen)

	gameHandler := &game.GameHandler{}
	hero := game.NewHero(playerName)
	gameState := game.NewGameState()

	enemy := gameHandler.CreateEnemy(gameState.CurrentEnemy)

	quit := make(chan bool)
	done := make(chan bool)

	ui.StartInputHandler(screen, hero, enemy, gameState, gameHandler, quit, done)

	gameHandler.StartBattle(hero, enemy, screen, gameState, quit, done)

	<-quit
}
