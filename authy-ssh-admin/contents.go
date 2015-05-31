package main

import (
	"fmt"
	"github.com/authy/onetouch-ssh"
	"github.com/jroimartin/gocui"
)

var (
	contentsViewID = "contents"
)

// Contents is a list of users.
type Contents struct {
	gui *gocui.Gui
}

// NewContents creates an instance of the Contents
func NewContents(gui *gocui.Gui) *Contents {
	list := &Contents{
		gui: gui,
	}

	return list
}

func (contents *Contents) drawLayout() {
	g := contents.gui
	maxX, maxY := g.Size()

	if v, err := g.SetView(contentsViewID, 30, -1, maxX, maxY-2); err != nil {
		v.Frame = true
	}
}

func (contents *Contents) view() *gocui.View {
	v, err := contents.gui.View(contentsViewID)
	if err != nil {
		panic(err)
	}

	return v
}

func (contents *Contents) setupKeyBindings() {
}

// OnUserSelected implements UsersListListener interface.
func (contents *Contents) OnUserSelected(user *ssh.User) {
	v := contents.view()

	v.Clear()
	fmt.Fprintln(v, user.Username)
}
