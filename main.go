package main

import (
	"fmt"
	"stuff/game"
	"stuff/ui"

	"github.com/gdamore/tcell/v2"
)

func main() {

	// Initialize screen
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

	// Starting Screen and Character Creation Screen
	playerName := ui.ShowStartScreen(screen)

	// Create game handler
	gameHandler := &game.GameHandler{}

	// Create hero
	hero := game.NewHero(playerName)

	// Initialize game state
	gameState := game.NewGameState()

	// Create first enemy
	enemy := gameHandler.CreateEnemy(gameState.CurrentEnemy)

	// Channels for game control
	quit := make(chan bool)
	done := make(chan bool)

	// Start the input handler
	ui.StartInputHandler(screen, hero, enemy, gameState, gameHandler, quit, done)

	// Start the first battle
	gameHandler.StartBattle(hero, enemy, screen, gameState, quit, done)

	// Wait for quit signal
	<-quit
}
