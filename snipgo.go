package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func main() {
	fmt.Println("Start")
	home, err := os.UserHomeDir()
	prefs := fmt.Sprintf("%s/%s", home, preferencesFile)
	jsonFile, err := os.Open(prefs)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	defer jsonFile.Close()
	fmt.Printf("Opened %s\n", prefs)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var userPrefs Preferences

	err = json.Unmarshal(byteValue, &userPrefs)
	if err != nil {
		fmt.Println(err)
		os.Exit(-2)
	}

	snippetFile, err := os.Open(userPrefs.DefaultFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer snippetFile.Close()
	fmt.Printf("Opened %s\n", userPrefs.DefaultFile)

	byteValue, _ = ioutil.ReadAll(snippetFile)
	var snippets []Snippet
	err = json.Unmarshal(byteValue, &snippets)
	if err != nil {
		fmt.Println(err)
		os.Exit(-2)
	}

	sort.Sort(ByCategoryTitle(snippets))
	for _, snippet := range snippets {
		fmt.Println(snippet.String())
	}
}
