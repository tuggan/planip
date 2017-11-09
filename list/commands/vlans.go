package commands

import (
	"fmt"

	"github.com/mitchellh/cli"
)

type ListVlansCommand struct {
	Ui cli.Ui
}

func (c *ListVlansCommand) Run(args []string) int {
	c.Ui.Output("List some vlans here please")
	c.Ui.Output(fmt.Sprintf("%+v", args))
	return 0
}

func (c *ListVlansCommand) Help() string {
	return `This will list all the vlans ever
	
	So just bear with me for a while plz`
}

func (c *ListVlansCommand) Synopsis() string {
	return "List all vlans"
}
