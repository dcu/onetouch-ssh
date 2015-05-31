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

	formInputs        []string
	currentInputIndex int
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

	contents.addFormLine("username", 0, 2)
	contents.addFormInput("country code", 3, 2)
	contents.addFormInput("phone number", 6, 2)
	contents.addFormInput("public keys", 9, maxY-14)
}

func (contents *Contents) addFormLine(label string, y int, height int) *gocui.View {
	g := contents.gui

	maxX, _ := g.Size()
	labelID := contentsViewID + "-label-" + label

	columnWidth := 30
	if v, err := g.SetView(labelID, columnWidth+1, y, columnWidth+20, y+height); err != nil {
		v.Frame = false
		fmt.Fprintf(v, label+":")
	}

	inputID := contents.idForInput(label)
	v, err := g.SetView(inputID, columnWidth+20, y, maxX-1, y+height)
	if err != nil {
		v.Frame = true
		v.Editable = false
	}

	return v
}

func (contents *Contents) addFormInput(label string, y int, height int) {
	view := contents.addFormLine(label, y, height)
	view.Editable = true

	contents.formInputs = append(contents.formInputs, label)
}

func (contents *Contents) idForInput(label string) string {
	return contentsViewID + "-input-" + label
}

func (contents *Contents) view(id string) *gocui.View {
	v, err := contents.gui.View(id)
	if err != nil {
		panic(err)
	}

	return v
}

func (contents *Contents) setupKeyBindings() {
	for _, label := range contents.formInputs {
		id := contents.idForInput(label)
		contents.gui.SetKeybinding(id, gocui.KeyEnter, gocui.ModNone, contents.nextFormInput)
	}
}

func (contents *Contents) nextFormInput(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.EditDelete(true)
	}

	nextPos := (contents.currentInputIndex + 1) % len(contents.formInputs)
	nextLabel := contents.formInputs[nextPos]

	contents.setFormInput(nextLabel)

	return nil
}

func (contents *Contents) setFormInput(label string) {
	for index, e := range contents.formInputs {
		viewID := contents.idForInput(e)
		view, err := contents.gui.View(viewID)
		if err != nil {
			panic(err)
		}

		if label == e {
			contents.currentInputIndex = index
			contents.gui.SetCurrentView(viewID)
			view.BgColor = gocui.ColorWhite
			view.FgColor = gocui.ColorBlack
		} else {
			view.BgColor = gocui.ColorDefault
			view.FgColor = gocui.ColorDefault
		}
	}
}

// OnUserSelected implements UsersListListener interface.
func (contents *Contents) OnUserSelected(user *ssh.User) {
	usernameInputID := contents.idForInput("username")
	v := contents.view(usernameInputID)

	v.Clear()
	fmt.Fprintln(v, user.Username)
}

// OnStartEditingUser implements UsersListListener interface.
func (contents *Contents) OnStartEditingUser(user *ssh.User) {
	contents.setFormInput("country code")
}
