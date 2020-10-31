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
	"sort"
	"strings"
)

// Snippet defines a snippet
type Snippet struct {
	Category string `json:"category"` // category of the snippet
	Title    string `json:"title"`    // title of the snippet
	Body     string `json:"body"`     // body of the snippet
}

// Snippet implements fmt.Stringer
var _ fmt.Stringer = (*Snippet)(nil)

// Snippets is a collection of Snippet
type Snippets []Snippet

// Snippets implements sort.Interface
var _ sort.Interface = (*Snippets)(nil)

// ReadSnippets accepts a filename, reads from the file and returns the Snippets found in the file.
func ReadSnippets(ctx *services.Context, filename string) (Snippets, error) {
	if filename == "" {
		preferences, err := ReadPreferences(ctx, "")
		if err != nil {
			return nil, err
		}
		filename = preferences.DefaultFile
	}
	snippetFile, err := ctx.Os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer snippetFile.Close()
	byteValue, err := ctx.IoUtil.ReadAll(snippetFile)
	if err != nil {
		return nil, err
	}
	var snippets []Snippet
	if err = ctx.Json.Unmarshal(byteValue, &snippets); err != nil {
		return nil, err
	}
	return snippets, nil
}

// WriteSnippets writes Snippets to the file named.
func (s *Snippets) WriteSnippets(ctx *services.Context, filename string) error {
	if filename == "" {
		preferences, err := ReadPreferences(ctx, "")
		if err != nil {
			return err
		}
		filename = preferences.DefaultFile
	}
	jsonString, err := ctx.Json.Marshal(s)
	if err != nil {
		return err
	}
	err = ctx.IoUtil.WriteFile(filename, jsonString, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Implement fmt.Stringer
func (s Snippet) String() string { return fmt.Sprintf("%s: %s", s.Category, s.Title) }

// Implement sort interface

// Len returns the length of the Snippets.
func (s Snippets) Len() int { return len(s) }

// Swap swaps two Snippet in the Snippets.
func (s Snippets) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less compares two Snippet.
func (s Snippets) Less(i, j int) bool {
	return (strings.ToLower(s[i].Category) < strings.ToLower(s[j].Category)) ||
		((strings.ToLower(s[i].Category) == strings.ToLower(s[j].Category)) &&
			(strings.ToLower(s[i].Title) < strings.ToLower(s[j].Title)))
}

// ByCategory splits Snippets into their Category.
func (s Snippets) ByCategory() Categories {
	catMap := make(map[string]Category)

	for _, snippet := range s {
		category, ok := catMap[snippet.Category]
		if !ok {
			category = Category{Name: snippet.Category}
		}
		category.Snippets = append(category.Snippets, snippet)
		catMap[category.Name] = category
	}

	var c []Category
	for _, v := range catMap {
		c = append(c, v)
	}
	sort.Sort(Categories(c))
	return c
}
