package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"os"
	"sort"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Could not get home directory")
	}
	prefs := fmt.Sprintf("%s/%s", home, preferencesFile)
	preferences, err := ReadPreferences(prefs)
	if err != nil {
		panic("Could not get preferences")
	}
	fmt.Printf("Snippets at: %s\n", preferences.DefaultFile)
	snippets, err := ReadSnippets(preferences.DefaultFile)
	if err != nil {
		panic("Could not read snippets")
	}
	fmt.Printf("Read %d Snippets\n", len(snippets))
	sort.Sort(ByCategoryTitle(snippets))
	for _, snippet := range snippets {
		fmt.Println(snippet)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Highlight = true
	g.SelFgColor = gocui.ColorRed

	// help := NewHelpWidget("help", 1, 1, "Help me I'm trapped in goland")
	status := NewStatusbarWidget("status", 1, 7, 50)
	butdown := NewButtonWidget("butdown", 52, 7, "DOWN", statusDown(status))
	butup := NewButtonWidget("butup", 58, 7, "UP", statusUp(status))
	g.SetManager(status, butdown, butup)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
