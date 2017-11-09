package commands

import (
	"fmt"

	"github.com/mitchellh/cli"
)

type AddDeviceCommand struct {
	Ui cli.Ui
}

func (c *AddDeviceCommand) Run(args []string) int {
	c.Ui.Output("Add a device here please")
	c.Ui.Output(fmt.Sprintf("%+v", args))
	return 0
}

func (c *AddDeviceCommand) Help() string {
	return `This will add a device.
	
	Just bear with me for a while plz`
}

func (c *AddDeviceCommand) Synopsis() string {
	return "Add a device"
}
