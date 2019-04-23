package gbio

import "github.com/Humpheh/goboy/pkg/gb"

// IOBinding provides an interface for display and input bindings.
type IOBinding interface {
	// Init the IOBinding
	Init(gameboy *gb.Gameboy, disableVsync bool)
	// RenderScreen renders a frame of the game.
	RenderScreen()
	// Destroy the IOBinding instance.
	Destroy()
	// ProcessInput processes input.
	ProcessInput()
	// SetTitle sets the title of the window.
	SetTitle(fps int)
	// IsRunning returns if the monitor is still running.
	IsRunning() bool
}
