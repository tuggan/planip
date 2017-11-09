package main

import (
	"log"
	"os"

	"github.com/tuggan/planip/add"
	"github.com/tuggan/planip/list"

	//_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/cli"
)

type ExitCommand struct {
	Ui cli.Ui
}

func (c *ExitCommand) Run(_ []string) int {
	c.Ui.Output("Will somehow exit application here, not sure how yet...")
	return 0
}

func (c *ExitCommand) Help() string {
	return "Exit applicaton"
}

func (c *ExitCommand) Synopsis() string {
	return "Exit the application"
}

func main() {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c := cli.NewCLI("planip", "0.0.1")

	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"exit": func() (cli.Command, error) {
			return &ExitCommand{
				Ui: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorRed,
				},
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &list.ListCommand{Ui: ui}, nil
		},
		"add": func() (cli.Command, error) {
			return &add.AddCommand{Ui: ui}, nil
		},
	}

	exitStatus, err := c.Run()

	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
