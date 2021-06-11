package stats

import (
	"time"

	"gitlab.com/4thlabs/perfmon/internal/stats/cpu"
	"gitlab.com/4thlabs/perfmon/internal/stats/network"
)

type Cpu struct {
	User   float64
	System float64
	Idle   float64
}

type Network struct {
	Rx, Tx int64
}

type Metrics struct {
	Cpu     Cpu
	Network Network
}

func Get() (Metrics, error) {
	before, err := cpu.Get()
	beforeIO, errIO := network.Get()

	if err != nil {
		return Metrics{}, err
	}

	if errIO != nil {
		return Metrics{}, errIO
	}

	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	afterIO, errIO := network.Get()

	if err != nil {
		return Metrics{}, err
	}

	if errIO != nil {
		return Metrics{}, err
	}

	total := float64(after.Total - before.Total)

	rxTotalBefore := 0
	txTotalBefore := 0
	for _, stat := range beforeIO {
		rxTotalBefore += int(stat.RxBytes)
		txTotalBefore += int(stat.TxBytes)
	}

	rxTotalAfter := 0
	txTotalAfter := 0
	for _, stat := range afterIO {
		rxTotalAfter += int(stat.RxBytes)
		txTotalAfter += int(stat.TxBytes)
	}

	return Metrics{
		Cpu: Cpu{
			User:   float64(after.User-before.User) / total * 100,
			System: float64(after.System-before.System) / total * 100,
			Idle:   float64(after.Idle-before.Idle) / total * 100,
		},
		Network: Network{
			Rx: int64(rxTotalAfter),
			Tx: int64(txTotalAfter),
		},
	}, nil
}
