package netw

import (
	"encoding/hex"
	"log"
	"net"
	"strconv"
)

type Broadcaster2 struct {
}

func NewBroadcaster2() *Broadcaster {
	return &Broadcaster{}
}

func makeSenderPool2(listeners int, threads int, ip string, encrypt bool) ([]chan []byte, error) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	channels := make([]chan []byte, threads)

	connNb := listeners / threads

	for i := 0; i < threads; i++ {
		channels[i] = make(chan []byte)
		dest := &net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: beginPort + i,
		}

		conn, err := net.DialUDP("udp", nil, dest)

		if err != nil {
			log.Panicln(err)
			return nil, err
		}

		go func(in <-chan []byte, conn *net.UDPConn, count int, idx int) {
			defer conn.Close()

			for {
				packet := <-in
				for i := 0; i < count; i++ {

					dest := &net.UDPAddr{
						IP:   net.ParseIP(ip),
						Port: (beginPort * idx) + i,
					}

					if encrypt {
						// var e error
						// packet, e = EncryptGCM(key, packet)

						// if e != nil {
						// 	log.Println(e)
						// }
						Encryt(key, packet)
					}

					_, err := conn.WriteToUDP(packet, dest)
					if err != nil {
						log.Println(err)
					}
					//log.Printf("Send packet with length %d", n)
				}
			}
		}(channels[i], conn, connNb, i)
	}

	return channels, nil
}

func (b *Broadcaster2) Start(address string, remote string, listeners int, threads int, encryt bool) error {
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

	senders, err := makeSenderPool2(listeners, threads, remote, encryt)
	if err != nil {
		return err
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
