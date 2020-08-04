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

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Highlight = true
	g.SelFgColor = gocui.ColorRed

	maxX, maxY := g.Size()

	titles := &ScrollView{
		name: "titles",
		x:    1,
		y:    1,
		w:    maxX / 2,
		h:    maxY / 2,
		body: "",
	}

	g.SetManager(titles)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	distinctTitles := make(map[string]bool)

	for _, snippet := range snippets {
		distinctTitles[snippet.Category] = true
	}

	for cat := range distinctTitles {
		titles.body += cat
		titles.body += "\n"
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
