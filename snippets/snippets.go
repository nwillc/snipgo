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

package snippets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Snippet struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Body     string `json:"body"`
}

// Snippet implements fmt.Stringer
var _ fmt.Stringer = (*Snippet)(nil)

type Snippets []Snippet

// Snippets implements sort.Interface
var _ sort.Interface = (*Snippets)(nil)

func ReadSnippets(filename string) (Snippets, error) {
	snippetFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer snippetFile.Close()
	byteValue, err := ioutil.ReadAll(snippetFile)
	if err != nil {
		return nil, err
	}
	var snippets []Snippet
	if err = json.Unmarshal(byteValue, &snippets); err != nil {
		return nil, err
	}
	return snippets, nil
}

// Implement fmt.Stringer
func (s Snippet) String() string { return fmt.Sprintf("%s: %s", s.Category, s.Title) }

// Implement sort interface
func (s Snippets) Len() int      { return len(s) }
func (s Snippets) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Snippets) Less(i, j int) bool {
	return (strings.ToLower(s[i].Category) < strings.ToLower(s[j].Category)) ||
		((strings.ToLower(s[i].Category) == strings.ToLower(s[j].Category)) &&
			(strings.ToLower(s[i].Title) < strings.ToLower(s[j].Title)))
}

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
