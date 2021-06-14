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
	RxBytes, TxBytes     int64
	RxPackets, TxPackets int64
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

	var rxTotalBefore, txTotalBefore, rxPacketsTotalBefore, txPacketsTotalBefore int64
	for _, stat := range beforeIO {
		rxTotalBefore += int64(stat.RxBytes)
		txTotalBefore += int64(stat.TxBytes)
		rxPacketsTotalBefore += int64(stat.RxPackets)
		txPacketsTotalBefore += int64(stat.TxPackets)
	}

	var rxTotalAfter, txTotalAfter, rxPacketsTotalAfter, txPacketsTotalAfter int64
	for _, stat := range afterIO {
		rxTotalAfter += int64(stat.RxBytes)
		txTotalAfter += int64(stat.TxBytes)
		rxPacketsTotalAfter += int64(stat.RxPackets)
		txPacketsTotalAfter += int64(stat.TxPackets)
	}

	return Metrics{
		Cpu: Cpu{
			User:   float64(after.User-before.User) / total * 100,
			System: float64(after.System-before.System) / total * 100,
			Idle:   float64(after.Idle-before.Idle) / total * 100,
		},
		Network: Network{
			RxBytes:   int64((rxTotalAfter - rxTotalBefore) / 1024 / 1024),
			TxBytes:   int64((txTotalAfter - txTotalBefore) / 1024 / 1024),
			RxPackets: int64(rxPacketsTotalAfter - rxPacketsTotalBefore),
			TxPackets: int64(txPacketsTotalAfter - txPacketsTotalBefore),
		},
	}, nil
}
