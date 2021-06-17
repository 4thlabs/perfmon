package commands

import (
	"github.com/alecthomas/kong"
	"gitlab.com/4thlabs/perfmon/internal/ui"
)

type MonitorCmd struct {
}

func (cmd *MonitorCmd) Run(ctx *kong.Context) error {
	ui.Init()
	return nil
}
