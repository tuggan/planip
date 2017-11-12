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
		c.Ui.Error("Atleast two arguments needed!")
		return 1
	}

	vlanID, err := strconv.Atoi(args[0])
	if err != nil {
		c.Ui.Error("First argument is the vlan id and has to be a number!")
	}

	if vlanID < 0 || vlanID > 4095 {
		c.Ui.Error(fmt.Sprintf("VLAN ID %d is not between 0 and 4095", vlanID))
		return 2
	}

	sitename := args[1]
	var name string
	if len(args) > 2 {
		name = args[2]
	}

	db, err := database.Open()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 3
	}

	err = db.Init()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 4
	}

	q := db.NewVLANQuery()

	if err = q.SetName(name); err != nil {
		c.Ui.Error(err.Error())
		return 6
	}

	if err = q.SetVLAN(vlanID); err != nil {
		c.Ui.Error(err.Error())
		return 7
	}

	if err = q.SetSiteName(sitename); err != nil {
		c.Ui.Error(err.Error())
		return 8
	}

	if err = q.PopulateSiteID(db); err != nil {
		c.Ui.Error("Failed while setting site id with message: " + err.Error())
	}

	if err = q.Add(db); err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 5
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
