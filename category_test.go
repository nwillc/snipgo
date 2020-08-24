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
