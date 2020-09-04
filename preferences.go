/*
 * Copyright (c) 2020, nwillc@gmail.com
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const preferencesFile = ".snippets.json"

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

func UserPreferences() (*Preferences, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Could not get home directory")
	}
	prefFileName := fmt.Sprintf("%s/%s", home, preferencesFile)
	preferences, err := readPreferences(prefFileName)
	if err != nil {
		panic(fmt.Sprintf("Could not read preferences: %s", prefFileName))
	}
	return preferences, nil
}
