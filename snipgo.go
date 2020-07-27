package main

import (
	"fmt"
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
}
