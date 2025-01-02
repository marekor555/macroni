package main

import (
	"fmt"
	"time"

	rg "github.com/go-vgo/robotgo"
	gh "github.com/robotn/gohook"
)

// executes a shortcut/combo
func makeShortcut(chars ...string) {
	for _, ch := range chars {
		switch ch {
		case "win":
			rg.KeyDown(rg.Cmd)
			break
		case "alt":
			rg.KeyDown(rg.Alt)
			break
		case "ctrl":
			rg.KeyDown(rg.Ctrl)
		default:
			rg.KeyDown(ch)
		}
	}
	time.Sleep(time.Millisecond * 250)
	for _, ch := range chars {
		switch ch {
		case "win":
			rg.KeyUp(rg.Cmd)
			break
		case "alt":
			rg.KeyUp(rg.Alt)
			break
		case "ctrl":
			rg.KeyUp(rg.Ctrl)
		default:
			rg.KeyUp(ch)
		}
	}
}

type Macro struct {
	Key             rune       // key to press (no win, ctrl, alt or F* allowed), lower letter
	ComboIntervalMs int        // interval in ms
	Combos          [][]string // array of combos, use shortened modifiers(win, ctrl, alt) and lower letters
}

func main() {
	// TODO: add reading from a file
	var macros []Macro = []Macro{
		{
			Key:             'o',
			ComboIntervalMs: 1000,
			Combos: [][]string{
				{"win", "f"},
				{"ctrl", "f"},
			},
		},
	}
	eventHook := gh.Start() // listen for events

	for e := range eventHook { // as eventHook is infinite stream, no external for loop needed
		if e.Kind == gh.KeyDown { // check if event is keydown
			for _, macro := range macros {
				if macro.Key == e.Keychar {
					for _, combo := range macro.Combos {
						makeShortcut(combo...)
						time.Sleep(time.Millisecond * time.Duration(macro.ComboIntervalMs))
					}
				}
			}
			fmt.Printf("Key down: rawcode=%d, keychar=%s\n", e.Rawcode, string(e.Keychar)) // keylogger :>
		}
	}

}
