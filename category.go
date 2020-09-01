package main

import (
	"sort"
	"strings"
)

type Category struct {
	Name     string
	Snippets Snippets
}

type Categories []Category

// Categories implements sort.Interface
var _ sort.Interface = (*Categories)(nil)

func (c Categories) Len() int {
	return len(c)
}

func (c Categories) Less(i, j int) bool {
	return strings.ToLower(c[i].Name) < strings.ToLower(c[j].Name)
}

func (c Categories) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func SnippetsByCategory(snippets Snippets) Categories {
	catMap := make(map[string]Category)

	for _, snippet := range snippets {
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
