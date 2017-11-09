package site

import (
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tuggan/planip/database"
)

type cmd struct {
	Ui cli.Ui
}

func New(ui cli.Ui) *cmd {
	c := &cmd{Ui: ui}
	return c
}

func (c *cmd) Run(args []string) int {
	nCom := 1
	c.Ui.Output("Add a site here please")
	c.Ui.Output(fmt.Sprintf("%+v", args))
	if len(args) != nCom {
		c.Ui.Error(fmt.Sprintf("This command only takes one argument, arguments given: %d", len(args)))
		for i := nCom; i < len(args); i++ {
			c.Ui.Error(fmt.Sprintf("Not needed: %s", args[i]))
		}
		return 1
	}

	db, err := database.Open()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 1
	}

	err = db.Init()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 1
	}

	err = db.AddSite(args[0])
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 1
	}

	err = db.Close()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 1
	}
	return 0
}

func (c *cmd) Help() string {
	return `This will add a site.
	
	Just bear with me for a while plz`
}

func (c *cmd) Synopsis() string {
	return "Add a site"
}
