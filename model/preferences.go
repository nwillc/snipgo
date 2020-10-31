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

package model

import (
	"fmt"
	"github.com/nwillc/snipgo/services"
	"os"
)

const preferencesFile = ".snippets.json" // default preferences file

// Preferences represents a users preferences.
type Preferences struct {
	DefaultFile string `json:"defaultFile"` // users default snippets file
}

// ReadPreferences reads the file indicated by filename and returns a *Preferences structure.
// If the filename is empty it will look for the default preferences file. If no file can be be
// found or the file is malformed an error is returned.
func ReadPreferences(ctx *services.Context, filename string) (*Preferences, error) {
	if filename == "" {
		home, err := ctx.Os.UserHomeDir()
		if err != nil {
			panic("Could not get home directory")
		}
		filename = fmt.Sprintf("%s/%s", home, preferencesFile)
	}
	jsonFile, err := ctx.Os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ctx.IoUtil.ReadAll(jsonFile)
	var preferences Preferences
	err = ctx.JSON.Unmarshal(byteValue, &preferences)
	if err != nil {
		return nil, err
	}
	return &preferences, nil
}

// Write the preferences as JSON to filename.
func (p *Preferences) Write(ctx *services.Context, filename string) error {
	jsonString, err := ctx.JSON.Marshal(p)
	if err != nil {
		return err
	}
	err = ctx.IoUtil.WriteFile(filename, jsonString, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
