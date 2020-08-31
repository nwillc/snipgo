package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const preferencesFile = ".Snippets.json"

type Preferences struct {
	DefaultFile string `json:"defaultFile"`
}

func readPreferences(filename string) (*Preferences, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var preferences Preferences
	err = json.Unmarshal(byteValue, &preferences)
	if err != nil {
		return nil, err
	}
	return &preferences, nil
}

func GetPreferences() (*Preferences, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Could not get home directory")
	}
	prefFileName := fmt.Sprintf("%s/%s", home, preferencesFile)
	preferences, err := readPreferences(prefFileName)
	if err != nil {
		panic("Could not read Snippets")
	}
	return preferences, nil
}
