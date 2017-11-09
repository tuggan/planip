package list

import (
	"github.com/mitchellh/cli"
	"github.com/tuggan/planip/list/commands"
	"github.com/tuggan/planip/list/sites"
)

type ListCommand struct {
	Ui cli.Ui
}

func (c *ListCommand) Run(args []string) int {
	list := cli.NewCLI("planip list", "")

	list.Args = args

	list.Commands = map[string]cli.CommandFactory{
		"sites": func() (cli.Command, error) {
			return sites.New(c.Ui), nil
		},
		"devices": func() (cli.Command, error) {
			return &commands.ListDevicesCommand{Ui: c.Ui}, nil
		},
		"vlans": func() (cli.Command, error) {
			return &commands.ListVlansCommand{Ui: c.Ui}, nil
		},
	}

	if exitstatus, err := list.Run(); err != nil {
		c.Ui.Error(err.Error())
		return exitstatus
	} else {
		return exitstatus
	}

}

func (c *ListCommand) Help() string {
	return "List things, yes :)"
}

func (c *ListCommand) Synopsis() string {
	return "List information, maybe"
}
