package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Snippet struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Body     string `json:"body"`
}

func ReadSnippets(filename string) ([]Snippet, error) {
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

// Implement Stringer
func (s Snippet) String() string { return fmt.Sprintf("%s: %s", s.Category, s.Title) }

type ByCategoryTitle []Snippet

// Implement sort interface
func (s ByCategoryTitle) Len() int      { return len(s) }
func (s ByCategoryTitle) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByCategoryTitle) Less(i, j int) bool {
	return (strings.ToLower(s[i].Category) < strings.ToLower(s[j].Category)) ||
		((strings.ToLower(s[i].Category) == strings.ToLower(s[j].Category)) &&
			(strings.ToLower(s[i].Title) < strings.ToLower(s[j].Title)))
}
