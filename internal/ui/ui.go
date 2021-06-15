package ui

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"gitlab.com/4thlabs/perfmon/internal/stats"
	ui "gitlab.com/4thlabs/perfmon/internal/ui/components"
)

const rootID = "root"
const redrawInterval = 250 * time.Millisecond

func Init() {
	t, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		log.Fatalln(err)
	}

	defer t.Close()

	c, err := container.New(t, container.ID(rootID))
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	comp, err := ui.NewMonitorUI(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	gridOpts, err := comp.Layout(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if err := c.Update(rootID, gridOpts...); err != nil {
		log.Fatalln(err)
	}

	go ui.Periodic(ctx, 1*time.Second, func() error {
		stats, err := stats.Get()
		if err != nil {
			return err
		}

		comp.Network.Reset()
		comp.Network.Write(fmt.Sprintf("[0] RxBytes: %d\n", stats.Network.RxBytes))
		comp.Network.Write(fmt.Sprintf("[1] TxBytes: %d\n", stats.Network.TxBytes))
		comp.Network.Write(fmt.Sprintf("[2] RxPackets: %d\n", stats.Network.RxPackets))
		comp.Network.Write(fmt.Sprintf("[3] TxPackets: %d", stats.Network.TxPackets))

		return comp.Cpu.Values([]int{int(stats.Cpu.User), int(stats.Cpu.System), int(stats.Cpu.Idle)}, 100)
	})

	quit := func(k *terminalapi.Keyboard) {
		if k.Key == keyboard.KeyEsc || k.Key == keyboard.KeyCtrlC {
			cancel()
		}
	}
	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quit), termdash.RedrawInterval(redrawInterval)); err != nil {
		log.Fatalln(err)
	}
}
