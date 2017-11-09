package add

import (
	"github.com/mitchellh/cli"
	"github.com/tuggan/planip/add/commands"
	"github.com/tuggan/planip/add/site"
)

type AddCommand struct {
	Ui cli.Ui
}

func (c *AddCommand) Run(args []string) int {
	add := cli.NewCLI("planip add", "")

	add.Args = args

	add.Commands = map[string]cli.CommandFactory{
		"site": func() (cli.Command, error) {
			return site.New(c.Ui), nil
		},
		"device": func() (cli.Command, error) {
			return &commands.AddDeviceCommand{Ui: c.Ui}, nil
		},
		"vlan": func() (cli.Command, error) {
			return &commands.AddVLANCommand{Ui: c.Ui}, nil
		},
	}

	if exitstatus, err := add.Run(); err != nil {
		c.Ui.Error(err.Error())
		return exitstatus
	} else {
		return exitstatus
	}
}

func (c *AddCommand) Help() string {
	return "Add things :)"
}

func (c *AddCommand) Synopsis() string {
	return "Add things, i hope :)"
}
