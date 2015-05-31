package main

import (
	"github.com/authy/onetouch-ssh"
	"github.com/jroimartin/gocui"
)

// App is the application.
type App struct {
	gui       *gocui.Gui
	UsersList *UsersList
	Contents  *Contents
	started   bool
}

// NewApp Instantiates a new App
func NewApp() *App {
	app := &App{}
	app.initGUI()

	ssh.NewUsersManager().AddListener(app)
	app.configureViews()

	app.gui.SetLayout(app.drawLayout)
	if err := app.keyBindings(); err != nil {
		panic(err)
	}

	return app
}

func (app *App) initGUI() {
	app.gui = gocui.NewGui()
	if err := app.gui.Init(); err != nil {
		panic(err)
	}

	app.gui.SelBgColor = gocui.ColorGreen
	app.gui.SelFgColor = gocui.ColorBlack
	app.gui.ShowCursor = true
}

func (app *App) configureViews() {
	app.UsersList = NewUsersList(app.gui)
	app.Contents = NewContents(app.gui)

	app.UsersList.AddListener(app.Contents)

}

// Start starts the application
func (app *App) Start() {
	defer app.gui.Close()

	err := app.gui.MainLoop()
	if err != nil && err != gocui.Quit {
		panic(err)
	}
}

func (app *App) drawLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("cmd-line", -1, maxY-2, maxX, maxY); err != nil {
		v.Editable = true
		v.Frame = false
	}

	app.UsersList.drawLayout()
	app.Contents.drawLayout()

	if v, err := g.SetView("status", -1, maxY-3, maxX, maxY-1); err != nil {
		v.Frame = false
		v.BgColor = gocui.ColorWhite
		v.FgColor = gocui.ColorBlack
	}

	if !app.started {
		manager := ssh.NewUsersManager()

		if manager.HasUsers() {
			app.UsersList.focus()
		} else {
			app.UsersList.showAddUserView(app.gui, nil)
		}
		app.started = true
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

func (app *App) keyBindings() error {
	app.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	app.gui.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone, app.UsersList.showAddUserView)

	app.UsersList.setupKeyBindings()
	app.Contents.setupKeyBindings()

	return nil
}

// OnUserAdded reports when a user was added.
func (app *App) OnUserAdded(user *ssh.User) {
}
