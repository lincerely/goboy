package iotui

import (
	"log"

	"github.com/Humpheh/goboy/pkg/bits"
	"github.com/Humpheh/goboy/pkg/gb"
	"github.com/gdamore/tcell"
)

func NewTuiIOBinding(gameboy *gb.Gameboy, disableVsync bool) *TuiIOBinding {
	monitor := TuiIOBinding{}
	monitor.Init(gameboy, disableVsync)
	return &monitor
}

type TuiIOBinding struct {
	Gameboy   *gb.Gameboy
	Screen    tcell.Screen
	EventQ    chan tcell.Event
	isRunning bool
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

	//listen to event
	t.EventQ = make(chan tcell.Event)
	go func() {
		for {
			ev := t.Screen.PollEvent()
			t.EventQ <- ev
		}
	}()

	t.isRunning = true
}

// RenderScreen renders a frame of the game.
func (t *TuiIOBinding) RenderScreen() {
	for y := 0; y < 144; y++ {
		for x := 0; x < 160; x++ {
			col := t.Gameboy.PreparedData[x][y]
			color := tcell.NewRGBColor(int32(col[0]), int32(col[1]), int32(col[2]))
			style := tcell.StyleDefault.Background(color)
			t.Screen.SetContent(x*3, y, ' ', nil, style)
			t.Screen.SetContent(x*3+1, y, ' ', nil, style)
			t.Screen.SetContent(x*3+2, y, ' ', nil, style)
		}
	}
	t.Screen.Show()
}

// Destroy the IOBinding instance.
func (t *TuiIOBinding) Destroy() {
	t.Screen.Fini()
}

// ProcessInput processes input.
func (t *TuiIOBinding) ProcessInput() {
	if !t.Gameboy.IsGameLoaded() || t.Gameboy.ExecutionPaused {
		return
	}

	var inputMask byte = 0xFF

	//readinputs
	func() {
		for {
			select {
			case ev := <-t.EventQ:
				switch ev := ev.(type) {
				case *tcell.EventResize:
					t.Screen.Sync()
				case *tcell.EventKey:
					switch ev.Key() {
					case tcell.KeyRune:
						switch ev.Rune() {
						case 'z': //A Button
							inputMask = bits.Reset(inputMask, 0)
						case 'x': //B Button
							inputMask = bits.Reset(inputMask, 1)
						}
					case tcell.KeyBackspace: //Select
						inputMask = bits.Reset(inputMask, 2)
					case tcell.KeyEnter: //Start
						inputMask = bits.Reset(inputMask, 3)
					case tcell.KeyRight:
						inputMask = bits.Reset(inputMask, 4)
					case tcell.KeyLeft:
						inputMask = bits.Reset(inputMask, 5)
					case tcell.KeyUp:
						inputMask = bits.Reset(inputMask, 6)
					case tcell.KeyDown:
						inputMask = bits.Reset(inputMask, 7)
					case tcell.KeyCtrlC:
						t.isRunning = false
					}
				}
			default:
				return
			}
		}
	}()

	for i := 0; i < 8; i++ {
		offset := byte(i)
		pressed := bits.Val(inputMask, offset) == 0
		released := !pressed && (bits.Val(t.Gameboy.InputMask, offset) == 0)

		if pressed {
			t.Gameboy.InputMask = bits.Reset(t.Gameboy.InputMask, offset)
			t.Gameboy.RequestJoypadInterrupt() // Joypad interrupt
		} else if released {
			t.Gameboy.InputMask = bits.Set(t.Gameboy.InputMask, offset)
		}
	}
}

// SetTitle sets the title of the window.
func (t *TuiIOBinding) SetTitle(fps int) {
	//nop
}

// IsRunning returns if the monitor is still running.
func (t *TuiIOBinding) IsRunning() bool {
	return t.isRunning
}
