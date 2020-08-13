package app

import "github.com/jroimartin/gocui"

type App struct {
	Gui *gocui.Gui
}

func NewApp() (*App, error) {
	app := &App{}
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	app.Gui = g
	return app, nil
}

func (app *App) Run() error {
	err := app.Gui.MainLoop()
	return err
}
