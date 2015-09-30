package main

import (
	"fmt"
	"github.com/dcu/onetouch-ssh"
	"github.com/jroimartin/gocui"
	"strings"
)

var (
	listViewID = "users-list"

	addUserBgViewID    = "add-user-background"
	addUserViewID      = "add-user"
	addUserViewLabelID = "add-user-label"
	addUserViewInputID = "add-user-input"

	columns = []string{"ID", "Username", "Email", "Phone Number", "Configured", "Protected"}
)

// UsersList is a list of users.
type UsersList struct {
	gui       *gocui.Gui
	listeners []UsersListListener
}

// UsersListListener is a interface for listing users list events.
type UsersListListener interface {
	OnUserSelected(user *ssh.User)
	OnStartEditingUser(user *ssh.User)
}

// NewUsersList creates an instance of the UsersList
func NewUsersList(g *gocui.Gui) *UsersList {
	list := &UsersList{
		gui: g,
	}

	return list
}

func (list *UsersList) setHelp() {
	setHelp(list.gui, `ctrl-a: add user | enter: edit user | up/down: select user | ctrl-c close app`)
}

// AddListener adds a listener for list's events.
func (list *UsersList) AddListener(listener UsersListListener) {
	list.listeners = append(list.listeners, listener)
}

func (list *UsersList) drawLayout() {
	maxX, maxY := list.gui.Size()
	columnSize := maxX / len(columns)

	if v, err := list.gui.SetView(listViewID+"-title", 2, -1, maxX-2, 1); err != nil {
		v.Highlight = false
		v.Editable = false
		v.Wrap = false

		for _, columnName := range columns {
			toFill := columnSize - len(columnName)
			fmt.Fprintf(v, columnName+strings.Repeat(" ", toFill))
		}
		fmt.Fprintln(v, "")
	}
	if v, err := list.gui.SetView(listViewID, 2, 1, maxX-2, maxY-2); err != nil {
		v.Highlight = true
		v.Editable = false
		v.Wrap = false

		manager := ssh.NewUsersManager()
		for _, user := range manager.Users() {
			for _, columnName := range columns {
				value := user.ValueForColumn(columnName)
				toFill := columnSize - len(value)
				fmt.Fprintf(v, value+strings.Repeat(" ", toFill))
			}
		}
		fmt.Fprintln(v, "")
	}
}

func (list *UsersList) view() *gocui.View {
	v, err := list.gui.View(listViewID)
	if err != nil {
		panic(err)
	}

	return v
}

func (list *UsersList) usernameToAdd() string {
	v, err := list.gui.View(addUserViewInputID)
	if err != nil {
		panic(err)
	}

	username, err := v.Line(0)
	if err != nil {
		return ""
	}

	return username
}

func (list *UsersList) setupKeyBindings() {
	list.gui.SetKeybinding(addUserViewInputID, gocui.KeyEnter, gocui.ModNone, list.validateAndAddUser)
	list.gui.SetKeybinding(addUserViewInputID, gocui.KeyCtrlU, gocui.ModNone, clearView)
	list.gui.SetKeybinding(addUserViewInputID, gocui.KeyCtrlD, gocui.ModNone, list.cancelAddUser)

	list.gui.SetKeybinding(listViewID, gocui.KeyEnter, gocui.ModNone, list.editCurrentUser)
	list.gui.SetKeybinding(listViewID, gocui.KeyArrowDown, gocui.ModNone, list.cursorDown)
	list.gui.SetKeybinding(listViewID, gocui.KeyArrowUp, gocui.ModNone, list.cursorUp)
}

func (list *UsersList) showAddUserView(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(addUserBgViewID, -1, -1, maxX, maxY); err != nil {
		v.BgColor = gocui.ColorWhite
	}

	if v, err := g.SetView(addUserViewID, maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		v.Editable = false
		v.Frame = true
	}

	if v, err := g.SetView(addUserViewLabelID, maxX/2-30+1, maxY/2, maxX/2-20+1, maxY/2+2); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "username:")
	}

	if v, err := g.SetView(addUserViewInputID, maxX/2-20+1, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		v.Frame = false
		v.Editable = true
		if err := g.SetCurrentView(addUserViewInputID); err != nil {
			return err
		}
	}

	return nil
}

func (list *UsersList) cancelAddUser(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView(addUserViewID)
	g.DeleteView(addUserViewLabelID)
	g.DeleteView(addUserViewInputID)
	g.DeleteView(addUserBgViewID)

	g.SetCurrentView(listViewID)
	return nil
}

func (list *UsersList) validateAndAddUser(g *gocui.Gui, v *gocui.View) error {
	username := list.usernameToAdd()
	list.cancelAddUser(g, v)

	if username == "" {
		return nil
	}

	v = list.view()

	manager := ssh.NewUsersManager()
	user := ssh.NewUser(username)

	err := manager.AddUser(user)
	if err == nil {
		fmt.Fprintln(v, username)
		return nil
	}

	panic(err)

	return nil
}

func (list *UsersList) selectedUsername() string {
	v := list.view()

	_, cy := v.Cursor()
	selected, err := v.Line(cy)
	if err != nil {
		selected = ""
	}

	return selected
}

func (list *UsersList) focus() {
	v := list.view()
	list.gui.SetCurrentView(listViewID)

	list.selectCurrentUser(list.gui, v)
}

func (list *UsersList) selectCurrentUser(g *gocui.Gui, v *gocui.View) error {
	list.setHelp()
	username := list.selectedUsername()
	if username == "" {
		return nil
	}

	manager := ssh.NewUsersManager()
	user := manager.LoadUser(username)

	for _, listener := range list.listeners {
		listener.OnUserSelected(user)
	}

	return nil
}

func (list *UsersList) editCurrentUser(g *gocui.Gui, v *gocui.View) error {
	username := list.selectedUsername()

	manager := ssh.NewUsersManager()
	user := manager.LoadUser(username)

	for _, listener := range list.listeners {
		listener.OnStartEditingUser(user)
	}

	return nil
}

func (list *UsersList) cursorUp(g *gocui.Gui, v *gocui.View) error {
	cursorUp(g, v)
	list.selectCurrentUser(g, v)
	return nil
}

func (list *UsersList) cursorDown(g *gocui.Gui, v *gocui.View) error {
	cursorDown(g, v)
	list.selectCurrentUser(g, v)
	return nil
}
