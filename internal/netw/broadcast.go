package netw

import (
	"encoding/hex"
	"log"
	"net"
	"strconv"
)

type Broadcaster struct {
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{}
}

type ConnectionPool []*net.UDPConn

const beginPort = int(1400)

func makeConnectionPool(ip string, nb int) (ConnectionPool, error) {
	pool := make(ConnectionPool, nb)

	for i := 0; i < nb; i++ {
		dest := &net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: beginPort + i,
		}

		conn, err := net.DialUDP("udp", nil, dest)

		if err != nil {
			log.Panicln(err)
			return nil, err
		}

		pool[i] = conn
	}

	return pool, nil
}

func makeSenderPool(conns ConnectionPool, threads int) ([]chan []byte, error) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	channels := make([]chan []byte, threads)

	connNb := len(conns) / threads

	for i := 0; i < threads; i++ {
		channels[i] = make(chan []byte)

		go func(in <-chan []byte, conns ConnectionPool, idx int) {
			for {
				packet := <-in
				for _, c := range conns {

					Encryt(key, packet)
					Hmac(key, packet)

					_, err := c.Write(packet)
					if err != nil {
						log.Println(err)
					}
					//log.Printf("Send packet with length %d", n)
				}
			}
		}(channels[i], conns[connNb*i:connNb*(i+1)], i)
	}

	return channels, nil
}

func (b *Broadcaster) Start(address string, remote string, listeners int, threads int) error {
	host, sport, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(sport)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	})

	if err != nil {
		return err
	}

	defer conn.Close()

	buffer := make([]byte, 2048)

	conns, err := makeConnectionPool(remote, listeners)
	if err != nil {
		return err
	}

	senders, err := makeSenderPool(conns, threads)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		defer conn.Close()
	}

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(err)
		}

		for _, sender := range senders {
			sender <- buffer[0:n]
		}

		//log.Printf("Received packet with length %d", n)
	}

	//return nil
}
