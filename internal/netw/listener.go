package netw

import (
	"fmt"
	"log"
	"net"
)

type Listener struct {
}

func NewListener() *Listener {
	return &Listener{}
}

func (l *Listener) Start(nb int) {
	for i := 0; i < nb; i++ {
		conn, err := net.ListenUDP("udp", &net.UDPAddr{
			Port: beginPort + i,
		})

		if err != nil {
			log.Fatalln(err)
		}

		defer conn.Close()
	}

	fmt.Scanln()
}
