package main

import (
	"fmt"
	"github.com/authy/onetouch-ssh"
	"github.com/jroimartin/gocui"
)

// UsersList is a list of users.
type UsersList struct {
	App *App
}

// NewUsersList creates an instance of the UsersList
func (app *App) NewUsersList() *UsersList {
	list := &UsersList{
		App: app,
	}

	return list
}

func (list *UsersList) drawLayout() {
	g := list.App.gui
	_, maxY := g.Size()

	if v, err := g.SetView("users-list", -1, -1, 30, maxY-2); err != nil {
		v.Highlight = true
		v.Editable = false
		v.Wrap = false

		manager := ssh.NewUsersManager()
		for _, user := range manager.Users() {
			fmt.Fprintln(v, user.Username)
		}
	}
}

func (list *UsersList) view() *gocui.View {
	v, err := list.App.gui.View("users-list")
	if err != nil {
		panic(err)
	}

	return v
}

func (list *UsersList) usernameToAdd() string {
	v, err := list.App.gui.View("add-user-input")
	if err != nil {
		panic(err)
	}

	username, _ := v.Line(-1)

	return username
}

func (list *UsersList) setupKeyBindings() {
	app := list.App
	app.gui.SetKeybinding("add-user-input", gocui.KeyEnter, gocui.ModNone, list.validateAndAddUser)
}

func (list *UsersList) showAddUserView(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("add-user", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		v.Editable = false
		v.Frame = true
	}

	if v, err := g.SetView("add-user-label", maxX/2-30+1, maxY/2, maxX/2-20+1, maxY/2+2); err != nil {
		v.Frame = false
		fmt.Fprintln(v, "username:")
	}

	if v, err := g.SetView("add-user-input", maxX/2-20+1, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		v.Frame = false
		v.Editable = true
		if err := g.SetCurrentView("add-user-input"); err != nil {
			return err
		}
	}

	return nil
}

func (list *UsersList) validateAndAddUser(g *gocui.Gui, v *gocui.View) error {
	username := list.usernameToAdd()

	g.DeleteView("add-user")
	g.DeleteView("add-user-label")
	g.DeleteView("add-user-input")

	v, err := g.View("users-list")
	if err != nil {
		panic(err)
	}
	g.SetCurrentView("users-list")

	manager := ssh.NewUsersManager()
	user := ssh.NewUser(username)

	err = manager.AddUser(user)
	if err == nil {
		fmt.Fprintln(v, username)
		return nil
	}

	return err
}
