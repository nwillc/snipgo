package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func getPreferences(filename string) (*Preferences, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s/%s", home, filename)
	jsonFile, err := os.Open(path)
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
