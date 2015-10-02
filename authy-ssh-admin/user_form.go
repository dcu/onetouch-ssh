package main

import (
	"fmt"
	"github.com/dcu/onetouch-ssh"
	"github.com/jroimartin/gocui"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

var (
	userFormViewID = "userForm"
)

// UserForm is a list of users.
type UserForm struct {
	gui  *gocui.Gui
	list *UsersList

	formInputs        []string
	viewIDs           []string
	currentInputIndex int
	id                string

	user *ssh.User
}

// NewUserForm creates an instance of the UserForm
func NewUserForm(gui *gocui.Gui, list *UsersList, user *ssh.User) *UserForm {
	formID := fmt.Sprintf("%s%d%d", user.Username, user.AuthyID, rand.Intn(math.MaxInt16))
	form := &UserForm{
		id:   formID,
		gui:  gui,
		user: user,
		list: list,
	}

	return form
}

func (userForm *UserForm) show() {
	userForm.drawLayout()
	userForm.setupKeyBindings()
	userForm.refresh()
	userForm.focus()

	if len(userForm.user.Username) == 0 {
		userForm.setFormInput("username")
	} else {
		userForm.setFormInput("email")
	}
}

func (userForm *UserForm) drawLayout() {
	userForm.formInputs = []string{}
	userForm.viewIDs = []string{}

	g := userForm.gui
	maxX, maxY := g.Size()

	if v, err := g.SetView(userFormViewID, -1, -1, maxX, maxY-4); err != nil {
		v.Frame = true
	}

	userForm.addFormLine("authy id", 0, 2)
	if len(userForm.user.Username) == 0 {
		userForm.addFormInput("username", 3, 2)
	} else {
		userForm.addFormLine("username", 3, 2)
	}
	userForm.addFormInput("email", 6, 2)
	userForm.addFormInput("country code", 9, 2)
	userForm.addFormInput("phone number", 12, 2)
	userForm.addFormInput("public keys", 15, maxY-18)
}

func (userForm *UserForm) addFormLine(label string, y int, height int) *gocui.View {
	g := userForm.gui

	maxX, _ := g.Size()
	labelID := userFormViewID + "-label-" + label + "-" + userForm.id
	width := 20

	if v, err := g.SetView(labelID, 0, y, width, y+height); err != nil {
		v.Frame = false
		fmt.Fprintf(v, label+":")

		userForm.viewIDs = append(userForm.viewIDs, labelID)
	}

	inputID := userForm.idForInput(label)
	v, err := g.SetView(inputID, width+1, y, maxX-3, y+height)
	if err != nil {
		ssh.Logger.Printf("Error adding view %s: %s", inputID, err)
	} else {
		ssh.Logger.Printf("Adding view %s", inputID)
	}
	if v == nil && err != nil {
		panic(err)
	}

	userForm.viewIDs = append(userForm.viewIDs, inputID)
	v.Frame = false
	v.Editable = false

	return v
}

func (userForm *UserForm) addFormInput(label string, y int, height int) {
	view := userForm.addFormLine(label, y, height)
	view.Editable = true
	view.Frame = true
	view.Wrap = false

	userForm.formInputs = append(userForm.formInputs, label)
}

func (userForm *UserForm) setFormLineValue(label string, value string) {
	id := userForm.idForInput(label)
	view, _ := userForm.gui.View(id)

	view.Clear()
	fmt.Fprintf(view, value)
}

func (userForm *UserForm) getFormLineValue(label string) []string {
	id := userForm.idForInput(label)
	view, err := userForm.gui.View(id)
	if err != nil {
		ssh.Logger.Printf("Error getting form line value for %s(%s): %s. Views: %#v", label, id, err, userForm.viewIDs)
		return []string{""}
	}

	return strings.Split(view.Buffer(), "\n")
}

func (userForm *UserForm) setHelp() {
	setHelp(
		userForm.gui,
		`enter: next field | ctrl-s: save user | ctrl-d go back to list`,
	)
}

func (userForm *UserForm) mainView() *gocui.View {
	v, err := userForm.gui.View(userFormViewID)
	if err != nil {
		panic(err)
	}

	return v
}

func (userForm *UserForm) focus() {
	userForm.gui.SetCurrentView(userFormViewID)
}

func (userForm *UserForm) idForInput(label string) string {
	return fmt.Sprintf("%s-input-%s-%s", userFormViewID, label, userForm.id)
}

func (userForm *UserForm) view(id string) *gocui.View {
	v, err := userForm.gui.View(id)
	if err != nil {
		panic(err)
	}

	return v
}

func (userForm *UserForm) closeView() {
	for _, viewID := range userForm.viewIDs {
		err := userForm.gui.DeleteView(viewID)
		if err != nil {
			ssh.Logger.Printf("Error closing view %s: %s", viewID, err)
		} else {
			ssh.Logger.Printf("Closing view %s", viewID)
		}
	}
	userForm.gui.DeleteView(userFormViewID)
	userForm.viewIDs = []string{}

	userForm.list.focus()
}

func (userForm *UserForm) setupKeyBindings() {
	for _, label := range userForm.formInputs {
		id := userForm.idForInput(label)
		userForm.gui.SetKeybinding(id, gocui.KeyCtrlS, gocui.ModNone, userForm.finishEditing)
		userForm.gui.SetKeybinding(id, gocui.KeyCtrlD, gocui.ModNone, userForm.discardChanges)
		userForm.gui.SetKeybinding(id, gocui.KeyEnter, gocui.ModNone, userForm.nextFormInput)
		userForm.gui.SetKeybinding(id, gocui.KeyTab, gocui.ModNone, userForm.nextFormInput)
	}
}

func (userForm *UserForm) nextFormInput(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
	}

	nextPos := (userForm.currentInputIndex + 1) % len(userForm.formInputs)
	nextLabel := userForm.formInputs[nextPos]

	userForm.setFormInput(nextLabel)

	return nil
}

func (userForm *UserForm) finishEditing(g *gocui.Gui, v *gocui.View) error {
	userForm.clearSelection()
	userForm.updateUser()

	userForm.refresh()
	manager := ssh.NewUsersManager()
	err := manager.UpdateUser(userForm.user)
	if err != nil {
		panic(err)
	}

	return nil
}

func (userForm *UserForm) discardChanges(g *gocui.Gui, v *gocui.View) error {
	userForm.closeView()
	err := g.SetCurrentView(listViewID)
	if err != nil {
		panic(err)
	}
	return nil
}

func (userForm *UserForm) currentFormInput() *gocui.View {
	label := userForm.formInputs[userForm.currentInputIndex]
	id := userForm.idForInput(label)
	return userForm.view(id)
}

func (userForm *UserForm) clearSelection() {
	userForm.setFormInput("")
}

func (userForm *UserForm) setFormInput(label string) {
	for index, e := range userForm.formInputs {
		viewID := userForm.idForInput(e)
		view, err := userForm.gui.View(viewID)
		if err != nil {
			continue
		}

		if label == e {
			userForm.currentInputIndex = index
			userForm.gui.SetCurrentView(viewID)
			//view.BgColor = gocui.ColorWhite
			//view.FgColor = gocui.ColorBlack

			userForm.setHelp()
		} else {
			view.BgColor = gocui.ColorDefault
			view.FgColor = gocui.ColorDefault
		}
	}
}

func (userForm *UserForm) refresh() {
	user := userForm.user

	userForm.setFormLineValue("username", user.Username)
	userForm.setFormLineValue("authy id", user.AuthyIDStr())
	userForm.setFormLineValue("email", user.Email)
	userForm.setFormLineValue("country code", user.CountryCodeStr())
	userForm.setFormLineValue("phone number", user.PhoneNumber)
	userForm.setFormLineValue("public keys", strings.Join(user.PublicKeys, "\n"))
}

func (userForm *UserForm) updateUser() {
	user := userForm.user
	if user == nil {
		return
	}

	countryCode, err := strconv.Atoi(userForm.getFormLineValue("country code")[0])
	if err != nil {
		panic(err)
	} else {
		user.CountryCode = countryCode
	}

	if len(user.Username) == 0 {
		user.Username = userForm.getFormLineValue("username")[0]
		ssh.Logger.Printf("Setting username to: %s", user.Username)
	}
	user.Email = userForm.getFormLineValue("email")[0]
	user.PhoneNumber = userForm.getFormLineValue("phone number")[0]
	user.PublicKeys = userForm.getFormLineValue("public keys")

	user.Register()
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
