package main

import (
	"fmt"
	"gladiator-sim/game"
	"gladiator-sim/ui"

	"github.com/gdamore/tcell/v2"
)

func main() {
	quit := make(chan struct{})

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

	// Start the engine and register quit inputs in UI
	ui := ui.NewUI(screen, quit)
	gameEngine := game.NewEngine(ui, quit)

	// This should take no argument and just use its internals
	gameEngine.StartGameLoop()
}
