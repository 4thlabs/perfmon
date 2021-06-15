package ui

import (
	"context"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/text"
)

type MonitorUI struct {
	Cpu     *barchart.BarChart
	Network *text.Text
}

// System Monitoring UI
func NewMonitorUI(ctx context.Context) (*MonitorUI, error) {

	cpu, err := newBarChart(ctx)
	if err != nil {
		return nil, err
	}

	return &MonitorUI{
		Cpu: cpu,
	}, nil
}

func (component *MonitorUI) Layout(ctx context.Context) ([]container.Option, error) {
	builder := grid.New()

	builder.Add(
		grid.RowHeightPerc(50,
			grid.ColWidthPerc(30,
				grid.Widget(component.Cpu,
					container.Border(linestyle.Light),
					container.BorderTitle("Cpu"),
				),
			),
			grid.ColWidthPerc(70),
		),
		grid.RowHeightPerc(50),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return gridOpts, nil
}

func newBarChart(ctx context.Context) (*barchart.BarChart, error) {
	bc, err := barchart.New(
		barchart.BarColors([]cell.Color{
			cell.ColorNumber(33),
			cell.ColorNumber(39),
			cell.ColorNumber(45),
		}),
		barchart.ValueColors([]cell.Color{
			cell.ColorWhite,
			cell.ColorWhite,
			cell.ColorWhite,
		}),
		barchart.Labels([]string{"User", "System", "Idle"}),
		barchart.ShowValues(),
	)

	if err != nil {
		return nil, err
	}

	return bc, nil
}

func Periodic(ctx context.Context, interval time.Duration, fn func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := fn(); err != nil {
				panic(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
