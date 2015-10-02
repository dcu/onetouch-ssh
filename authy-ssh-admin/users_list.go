package main

import (
	"errors"
	"fmt"
	"github.com/dcu/onetouch-ssh"
	"github.com/jroimartin/gocui"
	"regexp"
	"strings"
)

var (
	listViewID = "users-list"

	addUserBgViewID    = "add-user-background"
	addUserViewID      = "add-user"
	addUserViewLabelID = "add-user-label"
	addUserViewInputID = "add-user-input"

	columns        = []string{"ID", "Username", "Email", "Phone Number", "Configured", "Protected"}
	splitColumnsRx = regexp.MustCompile(`\s\s+`)
)

// UsersList is a list of users.
type UsersList struct {
	gui *gocui.Gui
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
	}

	list.refresh()
}

func (list *UsersList) refresh() {
	maxX, _ := list.gui.Size()
	columnSize := maxX / len(columns)

	v := list.view()
	v.Clear()
	manager := ssh.NewUsersManager()
	for _, user := range manager.Users() {
		for _, columnName := range columns {
			value := user.ValueForColumn(columnName)
			toFill := columnSize - len(value)
			fmt.Fprintf(v, value+strings.Repeat(" ", toFill))
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
	list.gui.SetKeybinding(addUserViewInputID, gocui.KeyCtrlU, gocui.ModNone, clearView)

	list.gui.SetKeybinding(listViewID, gocui.KeyEnter, gocui.ModNone, list.editCurrentUser)
	list.gui.SetKeybinding(listViewID, gocui.KeyArrowDown, gocui.ModNone, list.cursorDown)
	list.gui.SetKeybinding(listViewID, gocui.KeyArrowUp, gocui.ModNone, list.cursorUp)
}

func (list *UsersList) showAddUserView(g *gocui.Gui, v *gocui.View) error {
	user := ssh.NewUser("")
	form := NewUserForm(g, list, user)
	form.show()

	return nil
}

func (list *UsersList) selectedUsername() string {
	v := list.view()

	_, cy := v.Cursor()
	selected, err := v.Line(cy)
	if err != nil {
		selected = ""
	} else {
		// NOTE: this is a weak way to find the selected username
		fields := splitColumnsRx.Split(selected, 3)
		if len(fields) > 2 {
			selected = fields[1]
		}
	}

	return selected
}

func (list *UsersList) focus() {
	list.gui.SetCurrentView(listViewID)
	list.setHelp()
}

func (list *UsersList) editCurrentUser(g *gocui.Gui, v *gocui.View) error {
	username := list.selectedUsername()
	if username == "" {
		return nil
	}

	manager := ssh.NewUsersManager()
	user := manager.LoadUser(username)

	if user == nil {
		return errors.New("User " + username + " not found.")
	}

	form := NewUserForm(g, list, user)
	form.show()

	return nil
}

func (list *UsersList) cursorUp(g *gocui.Gui, v *gocui.View) error {
	cursorUp(g, v)
	return nil
}

func (list *UsersList) cursorDown(g *gocui.Gui, v *gocui.View) error {
	cursorDown(g, v)
	return nil
}
