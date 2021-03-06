package commands

import "gitlab.com/4thlabs/perfmon/internal/netw"

type BroadcastCmd struct {
	Listen    string `type:"string" help:"Listen address" default:"127.0.0.1:1234"`
	Remote    string `type:"string" help:"Remote address to target" default:"127.0.0.1"`
	NbThreads int    `type:"int" help:"Number of senders thread" default:15`
	Encrypt   bool   `type:"bool" help:"Packet Encryption" default=true`

	NbListeners int `type:"int" help:"Number of listeners" default:1000`
}

func (cmd *BroadcastCmd) Run() error {
	broadcaster := netw.NewBroadcaster()
	broadcaster.Start(cmd.Listen, cmd.Remote, cmd.NbListeners, cmd.NbThreads, cmd.Encrypt)
	return nil
}
