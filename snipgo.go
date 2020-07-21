package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func main() {
	/*
		g, err := gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			log.Panicln(err)
		}
		defer g.Close()
		g.SetManagerFunc(layout)

		if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
			log.Panicln(err)
		}

		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
	*/
	preferences, err := getPreferences(preferencesFile)
	if err != nil {
		panic("Could not get preferences")
	}
	fmt.Printf("Preferences: %+v\n", preferences)
	//fmt.Println("Start")
	//home, err := os.UserHomeDir()
	//prefs := fmt.Sprintf("%s/%s", home, preferencesFile)
	//jsonFile, err := os.Open(prefs)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}
	//
	//defer jsonFile.Close()
	//

	//fmt.Printf("Opened %s\n", prefs)
	//
	//byteValue, _ := ioutil.ReadAll(jsonFile)
	//
	//var userPrefs Preferences
	//
	//err = json.Unmarshal(byteValue, &userPrefs)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(-2)
	//}
	//
	//snippetFile, err := os.Open(userPrefs.DefaultFile)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}
	//defer snippetFile.Close()
	//fmt.Printf("Opened %s\n", userPrefs.DefaultFile)
	//
	//byteValue, _ = ioutil.ReadAll(snippetFile)
	//var snippets []Snippet
	//err = json.Unmarshal(byteValue, &snippets)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(-2)
	//}
	//
	//sort.Sort(ByCategoryTitle(snippets))
	//for _, snippet := range snippets {
	//	fmt.Println(snippet)
	//}
}

func layout(g *gocui.Gui) error {
	maxX, _ := g.Size()
	if _, err := g.SetView("hello", 1, 1, maxX/2, 4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
