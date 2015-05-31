package main

import (
	"fmt"
	"github.com/authy/onetouch-ssh"
	"github.com/jroimartin/gocui"
	"strconv"
	"strings"
)

var (
	contentsViewID = "contents"
)

// Contents is a list of users.
type Contents struct {
	gui *gocui.Gui

	formInputs        []string
	currentInputIndex int

	user *ssh.User
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
	contents.addFormLine("authy id", 3, 2)
	contents.addFormInput("country code", 6, 2)
	contents.addFormInput("phone number", 9, 2)
	contents.addFormInput("public keys", 12, maxY-16)
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
		v.Frame = false
		v.Editable = false
	}

	return v
}

func (contents *Contents) addFormInput(label string, y int, height int) {
	view := contents.addFormLine(label, y, height)
	view.Editable = true
	view.Frame = true
	view.Wrap = true

	contents.formInputs = append(contents.formInputs, label)
}

func (contents *Contents) setFormLineValue(label string, value string) {
	id := contents.idForInput(label)
	view, _ := contents.gui.View(id)

	view.Clear()
	fmt.Fprintf(view, value)
}

func (contents *Contents) getFormLineValue(label string) []string {
	id := contents.idForInput(label)
	view, _ := contents.gui.View(id)

	return strings.Split(view.Buffer(), "\n")
}

func (contents *Contents) setHelp() {
	setHelp(
		contents.gui,
		`enter: next field | ctrl-s: save user | ctrl-d discard changes`,
	)
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
		contents.gui.SetKeybinding(id, gocui.KeyCtrlS, gocui.ModNone, contents.finishEditing)
		contents.gui.SetKeybinding(id, gocui.KeyCtrlD, gocui.ModNone, contents.dischardChanges)
		contents.gui.SetKeybinding(id, gocui.KeyEnter, gocui.ModNone, contents.nextFormInput)
	}
}

func (contents *Contents) nextFormInput(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
	}

	nextPos := (contents.currentInputIndex + 1) % len(contents.formInputs)
	nextLabel := contents.formInputs[nextPos]

	contents.setFormInput(nextLabel)

	return nil
}

func (contents *Contents) finishEditing(g *gocui.Gui, v *gocui.View) error {
	contents.clearSelection()
	contents.updateUser()

	manager := ssh.NewUsersManager()
	err := manager.UpdateUser(contents.user)
	if err != nil {
		panic(err)
	}

	err = g.SetCurrentView(listViewID)
	if err != nil {
		panic(err)
	}
	return nil
}

func (contents *Contents) dischardChanges(g *gocui.Gui, v *gocui.View) error {
	contents.refresh()
	contents.clearSelection()
	err := g.SetCurrentView(listViewID)
	if err != nil {
		panic(err)
	}
	return nil
}

func (contents *Contents) currentFormInput() *gocui.View {
	label := contents.formInputs[contents.currentInputIndex]
	id := contents.idForInput(label)
	return contents.view(id)
}

func (contents *Contents) clearSelection() {
	contents.setFormInput("")
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
			//view.BgColor = gocui.ColorWhite
			//view.FgColor = gocui.ColorBlack

			contents.setHelp()
		} else {
			view.BgColor = gocui.ColorDefault
			view.FgColor = gocui.ColorDefault
		}
	}
}

func (contents *Contents) refresh() {
	user := contents.user
	contents.setFormLineValue("username", user.Username)
	contents.setFormLineValue("authy id", "<not set>")
	contents.setFormLineValue("country code", user.CountryCodeStr())
	contents.setFormLineValue("phone number", user.PhoneNumber)
	contents.setFormLineValue("public keys", strings.Join(user.PublicKeys, "\n"))
}

func (contents *Contents) updateUser() {
	user := contents.user
	countryCode, err := strconv.Atoi(contents.getFormLineValue("country code")[0])
	if err != nil {
		panic(err)
	} else {
		user.CountryCode = countryCode
	}
	user.PhoneNumber = contents.getFormLineValue("phone number")[0]
	user.PublicKeys = contents.getFormLineValue("public keys")

	user.Register()
}

// OnUserSelected implements UsersListListener interface.
func (contents *Contents) OnUserSelected(user *ssh.User) {
	contents.user = user
	contents.refresh()
}

// OnStartEditingUser implements UsersListListener interface.
func (contents *Contents) OnStartEditingUser(user *ssh.User) {
	contents.setFormInput("country code")
}

func editor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyArrowDown:
		_, y := v.Cursor()
		_, h := v.Size()
		lines := strings.Count(v.Buffer(), "\n")

		if y+1 < lines || lines < h {
			v.MoveCursor(0, 1, false)
		}
	case key == gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		x, y := v.Cursor()
		line, _ := v.Line(y)

		if x < len(line) {
			v.MoveCursor(1, 0, false)
		}
	}
}

func init() {
	gocui.Edit = editor
}
