package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const preferencesFile = ".Snippets.json"

type Preferences struct {
	DefaultFile string `json:"defaultFile"`
}

func ReadPreferences(filename string) (*Preferences, error) {
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
