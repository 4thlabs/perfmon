package main

import (
	"github.com/alecthomas/kong"
	"gitlab.com/4thlabs/perfmon/internal/commands"
)

var cli struct {
	Monitor   commands.MonitorCmd   `cmd help:"Performance Monitoring"`
	Stream    commands.StreamCmd    `cmd help:"Start streaming server"`
	Broadcast commands.BroadcastCmd `cmd help:"Listen to incoming streams and broadcast them to X listeners"`
	Listen    commands.ListenCmd    `cmd help:"Listen to un number of UDP Ports"`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
