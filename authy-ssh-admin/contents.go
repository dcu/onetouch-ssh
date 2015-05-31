package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

var (
	contentsViewID = "contents"
)

// Contents is a list of users.
type Contents struct {
	App *App
}

// NewContents creates an instance of the Contents
func (app *App) NewContents() *Contents {
	list := &Contents{
		App: app,
	}

	return list
}

func (list *Contents) drawLayout() {
	g := list.App.gui
	maxX, maxY := g.Size()

	if v, err := g.SetView(contentsViewID, 30, -1, maxX, maxY-2); err != nil {
		v.Frame = true
		fmt.Fprintln(v, "Something")
	}
}

func (list *Contents) view() *gocui.View {
	v, err := list.App.gui.View(contentsViewID)
	if err != nil {
		panic(err)
	}

	return v
}

func (list *Contents) setupKeyBindings() {
}
