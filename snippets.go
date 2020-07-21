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

// Implement Stringer
func (s Snippet) String() string { return fmt.Sprintf("%s: %s", s.Category, s.Title) }

type ByCategoryTitle []Snippet

func (s ByCategoryTitle) Len() int      { return len(s) }
func (s ByCategoryTitle) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByCategoryTitle) Less(i, j int) bool {
	return (strings.ToLower(s[i].Category) < strings.ToLower(s[j].Category)) ||
		((strings.ToLower(s[i].Category) == strings.ToLower(s[j].Category)) &&
			(strings.ToLower(s[i].Title) < strings.ToLower(s[j].Title)))
}
