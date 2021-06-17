package commands

import (
	"github.com/alecthomas/kong"
	"gitlab.com/4thlabs/perfmon/internal/netw"
)

type StreamCmd struct {
	Remote string `type:"string" help:"Remote address" default:"127.0.0.1:1234"`
	Pps    int    `type:"int" default:300 help:"Packet per Second"`
	File   string `arg name:"file" help:"Path to the recording file" type:"path"`
}

func (cmd *StreamCmd) Run(ctx *kong.Context) error {
	streamer := netw.NewStreamer()
	defer streamer.Close()

	err := streamer.Start(cmd.File, cmd.Remote, cmd.Pps)

	return err
}
