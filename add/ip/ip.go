package ip

import (
	"bytes"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/tuggan/planip/database"
)

type cmd struct {
	Ui       cli.Ui
	vlan     sql.NullInt64
	vlanName sql.NullString
	flags    *flag.FlagSet
}

func New(ui cli.Ui) *cmd {
	c := &cmd{Ui: ui}
	return c
}

func (c *cmd) Run(args []string) int {
	var err error

	c.flags = flag.NewFlagSet("vlan", flag.ExitOnError)

	vlan := c.flags.Int64("vlan", -1, "Set VLAN by VLAN?")
	vlanName := c.flags.String("vlan-name", "", "Set VLAN by name")
	c.flags.Parse(args)

	args = c.flags.Args()

	if vlan != nil && *vlan >= 0 && *vlan < 4096 {
		err = c.vlan.Scan(*vlan)
		if err != nil {
			c.Ui.Error(err.Error())
			return 1
		}
		fmt.Println(c.vlan.Int64)
	}
	if vlanName != nil {
		err = c.vlanName.Scan(*vlanName)
		if err != nil {
			c.Ui.Error(err.Error())
			return 2
		}
	}

	if len(args) < 1 {
		c.Ui.Error("This command needs atleast one argument")
		return 1
	}

	rest, nmask, err := parseNetmask(args[0])
	if err != nil {
		c.Ui.Error("Failed while parsing netmask")
		return 2
	}

	ip1, ip2, err := parseIPs(rest)

	if err != nil {
		c.Ui.Error(err.Error())
		return 3
	}

	db, err := database.Open()
	if err != nil {
		c.Ui.Error("Failed while opening database!")
		return 4
	}

	err = db.Init()
	if err != nil {
		c.Ui.Error("Failed whili initializing database!")
		return 5
	}

	network := database.IP{Netmask: nmask}

	if len(args) > 1 {
		err = network.SiteName.Scan(args[1])
		if err != nil {
			c.Ui.Error(err.Error())
		}
	}

	if ip1 != nil && ip2 != nil {
		err = addRange(db, ip1, ip2, network)
		if err != nil {
			c.Ui.Error(err.Error())
			return 6
		}
	} else if ip1 != nil {
		network.IP = ip1.String()
		err = db.AddIP(network)
		if err != nil {
			c.Ui.Error(err.Error())
			return 7
		}
	}

	return 0
}

func (c *cmd) Help() string {
	return `This will add a ip.
	
	Just bear with me for a while plz`
}

func (c *cmd) Synopsis() string {
	return "Add a site"
}

func parseNetmask(arg string) (string, sql.NullInt64, error) {
	split := strings.Split(arg, "/")
	mask := sql.NullInt64{Valid: false}

	if len(split) < 2 {
		return arg, mask, nil
	} else if len(split) > 2 {
		return arg, mask, errors.New("only one '/' allowed")
	}

	nmask, err := strconv.Atoi(split[1])
	if err != nil {
		return arg, mask, err
	}

	if nmask < 0 || nmask > 32 {
		return arg, mask, errors.New("netmasks are between 0 and 32")
	}

	mask.Int64 = int64(nmask)
	mask.Valid = true
	return split[0], mask, nil
}

func parseIPs(arg string) (net.IP, net.IP, error) {
	split := strings.Split(arg, "-")

	var ip1 net.IP
	var ip2 net.IP

	if len(split) < 2 {
		ip1 = net.ParseIP(arg)
		if ip1 != nil {
			return ip1, nil, nil
		}
		return nil, nil, errors.New("not a recognised IP adress")
	} else if len(split) > 2 {
		return nil, nil, errors.New("more than one '-' in range argument")
	}

	ip1 = net.ParseIP(split[0])
	if ip1 == nil {
		return nil, nil, errors.New("first ip not recognised as actual ip")
	}
	ip2 = net.ParseIP(split[1])
	if ip2 == nil {
		return nil, nil, errors.New("first ip not recognised as actual ip")
	}

	return ip1, ip2, nil
}

func inc(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func addRange(db *database.DBS, ip1, ip2 net.IP, data database.IP) error {
	ip := make(net.IP, len(ip1))
	for copy(ip, ip1); bytes.Compare(ip, ip2) <= 0; inc(ip) {
		data.IP = ip.String()
		err := db.AddIP(data)
		if err != nil {
			return err
		}
	}
	return nil
}
