package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const preferences = ".snippets.json"

type Preferences struct {
	DefaultFile string `json:"defaultFile"`
}

type Snippet struct {
	Category string `json:"category"`
	Title string `json:"title"`
	Body string `json:"body"`
}

func main() {
	fmt.Println("Start")
	home, err := os.UserHomeDir()
	prefs := fmt.Sprintf("%s/%s", home, preferences)
	jsonFile, err := os.Open(prefs)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Printf("Opened %s\n", prefs)
	defer jsonFile.Close()

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
	fmt.Printf("Opened %s\n", userPrefs.DefaultFile)
	defer snippetFile.Close()

	byteValue, _ = ioutil.ReadAll(snippetFile)
	var snippets []Snippet
	err = json.Unmarshal(byteValue, &snippets)
	if err != nil {
		fmt.Println(err)
		os.Exit(-2)
	}

	for _, snippet := range snippets {
		fmt.Printf("%s: %s\n%s\n", snippet.Category, snippet.Title, snippet.Body)
	}
}
