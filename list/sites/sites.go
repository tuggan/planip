package sites

import (
	"bytes"
	"fmt"
	"text/tabwriter"

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
	nCom := 0
	c.Ui.Output("Listing sites :)")
	c.Ui.Output(fmt.Sprintf("%+v", args))
	if len(args) != nCom {
		c.Ui.Error(fmt.Sprintf("This command takes no argument, arguments given: %d", len(args)))
		for i := nCom; i < len(args); i++ {
			c.Ui.Error(fmt.Sprintf("Not needed: %s", args[i]))
		}
		return 1
	}

	db, err := database.Open()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 2
	}

	sites, err := db.GetSites()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 3
	}
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "ID", "Name", "Created", "Changed", "Comment")
	for _, v := range sites {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", v.ID, v.Name, v.Created, v.Changed, v.Comment.String)
		//c.Ui.Output(fmt.Sprintf("%d %s %s %s %s", v.ID, v.Name, v.Created, v.Changed, v.Comment.String))
	}
	w.Flush()
	c.Ui.Output(buf.String())

	err = db.Close()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed with message: %s", err))
		return 4
	}
	return 0
}

func (c *cmd) Help() string {
	return `This will list all sites.
	
	Just bear with me for a while plz`
}

func (c *cmd) Synopsis() string {
	return "List sites"
}
