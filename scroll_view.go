package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type ScrollView struct {
	name string
	x, y int
	w, h int
	body string
}

var _ gocui.Manager = (*ScrollView)(nil) // Implements gocui.Manager

func (s ScrollView) Layout(gui *gocui.Gui) error {
	v, err := gui.SetView(s.name, s.x, s.y, s.x+s.w, s.y+s.h)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = false
		fmt.Fprint(v, s.body)
	}
	return nil
}
