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
			name:       "Empty",
			categories: Categories{},
			sorted:     Categories{},
		},
		{
			name:       "Sorted",
			categories: Categories{{Name: "A"}, {Name: "B"}, {Name: "C"}},
			sorted:     Categories{{Name: "A"}, {Name: "B"}, {Name: "C"}},
		},
		{
			name:       "Unsorted",
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

func TestToSnippets(t *testing.T) {
	categories := Categories{
		{
			Name:     "A",
			Snippets: []Snippet{{"A", "A", "A"}},
		},
		{
			Name:     "B",
			Snippets: []Snippet{{"B", "B", "B"}, {"B", "B2", "B2"}},
		},
	}
	snippets := categories.ToSnippets()
	assert.Len(t, *snippets, 3)
}
