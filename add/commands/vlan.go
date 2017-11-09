package commands

import (
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tuggan/planip/database"
)

type AddVLANCommand struct {
	Ui cli.Ui
}

func (c *AddVLANCommand) Run(args []string) int {
	c.Ui.Output("Add a VLAN here please")
	c.Ui.Output(fmt.Sprintf("%+v", args))
	db, err := database.Open()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 1
	}
	//err = db.InitiateDatabase()
	//if err != nil {
	//	c.Ui.Error(fmt.Sprintf("Faid with message: %s", err))
	//	return 2
	//}

	err = db.InitiateSites()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Faid with message: %s", err))
		return 2
	}

	err = db.InitiateVLANS()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Faid with message: %s", err))
		return 3
	}

	err = db.Close()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Faid with message: %s", err))
		return 4
	}
	return 0
}

func (c *AddVLANCommand) Help() string {
	return `This will add a VLAN.
	
	Just bear with me for a while plz`
}

func (c *AddVLANCommand) Synopsis() string {
	return "Add a VLAN"
}
