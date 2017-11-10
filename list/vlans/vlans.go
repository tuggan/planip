package vlans

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/tuggan/planip/database"

	"github.com/mitchellh/cli"
)

type cmd struct {
	Ui cli.Ui
}

func New(ui cli.Ui) *cmd {
	c := &cmd{Ui: ui}
	return c
}

func (c *cmd) Run(args []string) int {
	c.Ui.Output("List some vlans here please")
	c.Ui.Output(fmt.Sprintf("%+v", args))

	db, err := database.Open()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 2
	}

	vlans, err := db.GetVLANS()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 3
	}

	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "%s\t%s\t%s\n", "Name", "VLAN", "Site")
	for _, v := range vlans {
		fmt.Fprintf(w, "%s\t%d\t%s\n", v.Name, v.Vlan, v.SiteName)
	}
	w.Flush()
	c.Ui.Output(buf.String())

	return 0
}

func (c *cmd) Help() string {
	return `This will list all the vlans ever
	
	So just bear with me for a while plz`
}

func (c *cmd) Synopsis() string {
	return "List all vlans"
}
