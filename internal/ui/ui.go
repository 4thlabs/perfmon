package ui

import (
	"fmt"
	"log"
	"math"

	term "github.com/gizak/termui/v3"
	"gitlab.com/4thlabs/perfmon/internal/stats"
	ui "gitlab.com/4thlabs/perfmon/internal/ui/components"
)

func Init() {
	if err := term.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	monitorUI := ui.NewMonitorUI()

	defer term.Close()

	term.Render(monitorUI)

	go func() {
		for {
			s, err := stats.Get()
			if err == nil {
				monitorUI.Cpu.Data = []float64{float64(math.Round(s.Cpu.User)), float64(math.Round(s.Cpu.System)), float64(math.Round(s.Cpu.Idle))}
				monitorUI.Network.Rows = []string{
					fmt.Sprintf("[0] RxBytes %d", s.Network.RxBytes),
					fmt.Sprintf("[1] TxBytes %d", s.Network.TxBytes),
					fmt.Sprintf("[2] RxPacket %d", s.Network.RxPackets),
				}
				term.Render(monitorUI)
			} else {
				log.Println(err)
			}
		}
	}()

	for e := range term.PollEvents() {
		if e.Type == term.KeyboardEvent {
			break
		}
	}
}
