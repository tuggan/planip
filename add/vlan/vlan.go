package vlan

import (
	"fmt"
	"strconv"

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
	c.Ui.Output("Add VLAN here")
	c.Ui.Output(fmt.Sprintf("%+v", args))
	if len(args) < 2 {
		return 1
	}

	vlanID, err := strconv.Atoi(args[0])
	if err != nil {
		c.Ui.Error("First argument is the vlan id and has to be a number!")
	}

	if vlanID < 0 || vlanID > 4095 {
		c.Ui.Error("VLAN ID has to be a number between 0 and 4095")
	}

	sitename := args[1]

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

	err = db.AddVLAN(database.Vlan{Vlan: int16(vlanID), SiteName: sitename})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Added vlan %d to site %s", vlanID, sitename))

	return 0
}

func (c *cmd) Help() string {
	return `This will add a VLAN.
	
	Just bear with me for a while plz`
}

func (c *cmd) Synopsis() string {
	return "Add a VLAN"
}
