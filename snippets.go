package main

import (
	"fmt"
	"strings"
)

const preferencesFile = ".snippets.json"

type Preferences struct {
	DefaultFile string `json:"defaultFile"`
}

type Snippet struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Body     string `json:"body"`
}

func (s Snippet) String() string { return fmt.Sprintf("%s: %s", s.Category, s.Title) }

type ByCategory []Snippet

func (s ByCategory) Len() int      { return len(s) }
func (s ByCategory) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByCategory) Less(i, j int) bool {
	return strings.ToLower(s[i].Category) < strings.ToLower(s[j].Category)
}
