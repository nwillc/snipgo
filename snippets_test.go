package main

import (
	"testing"
)

func TestDummy(t *testing.T) {
	low := Snippet{
		Category: "A",
		Title:    "A",
	}

	high := Snippet{
		Category: "A",
		Title:    "B",
	}

	snippets := ByCategoryTitle([]Snippet{low, high})
	if snippets[0].Category > snippets[1].Category {
		t.Error("Out of order")
	}
	if snippets[0].Title > snippets[1].Title {
		t.Error("Out of order")
	}

	if !snippets.Less(0, 1) {
		t.Error("Less failed")
	}

	snippets.Swap(0, 1)
	if snippets.Less(0, 1) {
		t.Error("Swap failed")
	}
}
