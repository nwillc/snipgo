package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/nwillc/snipgo/app"
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

	app, err := app.NewApp()
	if err != nil {
		log.Panicln(err)
	}
	defer app.Gui.Close()

	maxX, maxY := app.Gui.Size()

	categoryView := &ScrollView{
		name: "titles",
		x:    1,
		y:    1,
		w:    maxX / 2,
		h:    maxY / 2,
		body: "",
	}

	app.Gui.SetManager(categoryView)

	if err := app.Gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	distinctCategories := make(map[string]bool)
	categoryList := make([]string, 0)

	for _, snippet := range snippets {
		_, ok := distinctCategories[snippet.Category]
		if !ok {
			distinctCategories[snippet.Category] = true
			categoryList = append(categoryList, snippet.Category)
		}
	}

	sort.Strings(categoryList)
	for _, cat := range categoryList {
		categoryView.body += cat
		categoryView.body += "\n"
	}

	if err := app.Gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
