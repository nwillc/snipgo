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
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestCategorySort(t *testing.T) {
	tests := []struct {
		name       string
		categories Categories
		sorted     Categories
	}{
		{
			name:       "empty",
			categories: Categories{},
			sorted:     Categories{},
		},
		{
			name:       "Sorted already",
			categories: Categories{{Name: "A"}, {Name: "B"}, {Name: "C"}},
			sorted:     Categories{{Name: "A"}, {Name: "B"}, {Name: "C"}},
		},
		{
			name:       "Not sorted",
			categories: Categories{{Name: "B"}, {Name: "C"}, {Name: "A"}},
			sorted:     Categories{{Name: "A"}, {Name: "B"}, {Name: "C"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, len(test.categories), len(test.sorted))
			sort.Sort(test.categories)
			for i, c := range test.categories {
				assert.Equal(t, test.sorted[i].Name, c.Name)
			}
		})
	}
}

func TestSnippetsByCategory(t *testing.T) {
	tests := []struct {
		name       string
		snippets   Snippets
		categories Categories
	}{
		{
			name:       "Empty",
			snippets:   Snippets{},
			categories: Categories{},
		},
		{
			name: "Two Snippets One Category",
			snippets: Snippets{
				{
					Category: "A",
					Title:    "A",
				},
				{
					Category: "A",
					Title:    "B",
				},
			},
			categories: Categories{
				{
					Name:     "A",
					Snippets: []Snippet{{Category: "A", Title: "A"}, {Category: "A", Title: "B"}},
				},
			},
		},
		{
			name: "Two Snippets Two Categories",
			snippets: Snippets{
				{
					Category: "A",
					Title:    "A",
				},
				{
					Category: "B",
					Title:    "B",
				},
			},
			categories: Categories{
				{
					Name:     "A",
					Snippets: []Snippet{{Category: "A", Title: "A"}},
				},
				{
					Name:     "B",
					Snippets: []Snippet{{Category: "B", Title: "B"}},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := SnippetsByCategory(test.snippets)
			assert.Equal(t, len(test.categories), len(result))
			for i, category := range test.categories {
				assert.Equal(t, category.Name, result[i].Name)
				assert.Equal(t, len(category.Snippets), len(result[i].Snippets))
				if len(category.Snippets) == len(result[i].Snippets) {
					for j, snippet := range category.Snippets {
						assert.Equal(t, snippet.Title, result[i].Snippets[j].Title)
					}
				}
			}
		})
	}
}
