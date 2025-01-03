package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	rg "github.com/go-vgo/robotgo"
	gh "github.com/robotn/gohook"
)

// executes a shortcut/combo
func makeShortcut(interval int, chars ...string) {
	for _, ch := range chars {
		switch ch {
		case "win":
			rg.KeyDown(rg.Cmd)
		case "alt":
			rg.KeyDown(rg.Alt)
		case "ctrl":
			rg.KeyDown(rg.Ctrl)
		case "enter":
			rg.KeyDown(rg.Enter)
		default:
			if len(ch) > 1 {
				rg.TypeStr(ch)
				return
			}
			rg.KeyDown(ch)
		}
	}
	time.Sleep(time.Millisecond * time.Duration(interval))
	for _, ch := range chars {
		switch ch {
		case "win":
			rg.KeyUp(rg.Cmd)
		case "alt":
			rg.KeyUp(rg.Alt)
		case "ctrl":
			rg.KeyUp(rg.Ctrl)
		case "enter":
			rg.KeyUp(rg.Enter)
		default:
			rg.KeyUp(ch)
		}
	}
}

type Combo struct {
	Keypress []string
	Interval int
	HoldTime int
}

type Macro struct {
	Key    string  // key to press (no win, ctrl, alt or F* allowed), lower letter
	Combos []Combo // array of combos, use shortened modifiers(win, ctrl, alt) and lower letters
}

func main() {
	var macros []Macro = []Macro{
		{
			Key: "m",
			Combos: []Combo{
				{Keypress: []string{"win", "f"}, Interval: 0, HoldTime: 200},
				{Keypress: []string{"go.dev/dl"}, Interval: 1000, HoldTime: 2000},
				{Keypress: []string{"enter"}, Interval: 10, HoldTime: 200},
			},
		},
	}

	data, err := os.ReadFile("macro.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = json.Unmarshal(data, &macros)
	if err != nil {
		log.Fatal(err.Error())
	}

	eventHook := gh.Start() // listen for events

	for e := range eventHook { // as eventHook is infinite stream, no external for loop needed
		if e.Kind == gh.KeyDown { // check if event is keydown
			for _, macro := range macros {
				if macro.Key == string(e.Keychar) {
					for _, combo := range macro.Combos {
						time.Sleep(time.Millisecond * time.Duration(combo.Interval))
						makeShortcut(combo.HoldTime, combo.Keypress...)
					}
				}
			}
			fmt.Printf("Key down: rawcode=%d, keychar=%s\n", e.Rawcode, string(e.Keychar)) // keylogger :>
		}
	}

}
