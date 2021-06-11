package ui

import (
	"image"
	"sync"

	term "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type MonitorUI struct {
	Grid    *term.Grid
	Cpu     *widgets.BarChart
	Network *widgets.List

	sync.Mutex
}

// System Monitoring UI
func NewMonitorUI() *MonitorUI {
	component := &MonitorUI{}

	bc := widgets.NewBarChart()
	bc.Data = []float64{0, 0, 0}
	bc.Labels = []string{"User", "System", "Idle"}
	bc.Title = "CPU Usage (%)"
	//bc.BarWidth = 4
	bc.BarGap = 4

	bc.BarColors = []term.Color{term.ColorRed, term.ColorGreen, term.ColorBlue}
	bc.LabelStyles = []term.Style{term.NewStyle(term.ColorWhite)}
	bc.NumStyles = []term.Style{term.NewStyle(term.ColorWhite)}

	net := widgets.NewList()
	net.Title = "Network"

	grid := term.NewGrid()
	termWidth, termHeight := term.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		term.NewRow(1.0/2.0,
			term.NewCol(1.0/4.0, bc),
			term.NewCol(2.0/4.0, net),
		),
	)

	component.Grid = grid
	component.Cpu = bc
	component.Network = net

	return component
}

func (ui *MonitorUI) GetRect() image.Rectangle {
	return ui.Grid.GetRect()
}

func (ui *MonitorUI) SetRect(x1 int, y1 int, x2 int, y2 int) {
	ui.Grid.SetRect(x1, y1, x2, y2)
}

func (ui *MonitorUI) Draw(buffer *term.Buffer) {
	ui.Grid.Draw(buffer)
}
