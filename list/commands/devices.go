package commands

import (
	"fmt"

	"github.com/mitchellh/cli"
)

type ListDevicesCommand struct {
	Ui cli.Ui
}

func (c *ListDevicesCommand) Run(args []string) int {
	c.Ui.Output("List some devices here please")
	c.Ui.Output(fmt.Sprintf("%+v", args))
	return 0
}

func (c *ListDevicesCommand) Help() string {
	return `This will list all the devices ever
	
	So just bear with me for a while plz`
}

func (c *ListDevicesCommand) Synopsis() string {
	return "List all devices"
}
