package commands

import "gitlab.com/4thlabs/perfmon/internal/netw"

type BroadcastCmd struct {
	Listen      string `type:"string" help:"Listen address" default:"127.0.0.1:1234"`
	NbListeners int    `type:"int" help:"Number of listeners" default:1000`
}

func (cmd *BroadcastCmd) Run() error {
	broadcaster := netw.NewBroadcaster()
	broadcaster.Start(cmd.Listen, cmd.NbListeners)
	return nil
}
