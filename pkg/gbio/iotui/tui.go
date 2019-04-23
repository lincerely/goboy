package iotui

import (
	"log"

	"github.com/Humpheh/goboy/pkg/gb"
	"github.com/gdamore/tcell"
)

func NewTuiIOBinding(gameboy *gb.Gameboy, disableVsync bool) *TuiIOBinding {
	monitor := TuiIOBinding{}
	monitor.Init(gameboy, disableVsync)
	return &monitor
}

type TuiIOBinding struct {
	Gameboy *gb.Gameboy
	Screen  tcell.Screen
}

// Init the IOBinding
func (t *TuiIOBinding) Init(gameboy *gb.Gameboy, disableVsync bool) {
	t.Gameboy = gameboy

	var err error
	if t.Screen, err = tcell.NewScreen(); err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}

	if err = t.Screen.Init(); err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}

	t.Screen.Clear()
}

// RenderScreen renders a frame of the game.
func (t *TuiIOBinding) RenderScreen() {

	//	r, g, b := gb.GetPaletteColour(3)
	//	bgColor := tcell.NewRGBColor(r, g, b)
	//	t.Screen.Fill(' ', tcell.StyleDefault.Background(bgColor))

	for y := 0; y < 144; y++ {
		for x := 0; x < 160; x++ {
			col := t.Gameboy.PreparedData[x][y]
			color := tcell.NewRGBColor(int32(col[0]), int32(col[1]), int32(col[2]))
			style := tcell.StyleDefault.Background(color)
			t.Screen.SetCell(x, y, style, ' ')
		}
	}

	t.Screen.Show()
}

// Destroy the IOBinding instance.
func (t *TuiIOBinding) Destroy() {
	t.Screen.Fini()
	t.Screen = nil
}

// ProcessInput processes input.
func (t *TuiIOBinding) ProcessInput() {

}

// SetTitle sets the title of the window.
func (t *TuiIOBinding) SetTitle(fps int) {
	//nop
}

// IsRunning returns if the monitor is still running.
func (t *TuiIOBinding) IsRunning() bool {
	return t.Screen != nil
}
