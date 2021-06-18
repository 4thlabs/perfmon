package commands

import (
	"github.com/alecthomas/kong"
	"gitlab.com/4thlabs/perfmon/internal/netw"
)

type ListenCmd struct {
	NbPorts int `type:"int" help:"Number of ports to listen" default:1000`
}

func (cmd *ListenCmd) Run(ctx *kong.Context) error {
	listener := netw.NewListener()
	listener.Start(cmd.NbPorts)

	return nil
}
