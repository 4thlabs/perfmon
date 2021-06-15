package udp

import (
	"io"
	"log"
	"net"
	"runtime"

	"gitlab.com/4thlabs/perfmon/internal/recording"
)

type Client struct {
	addr net.UDPAddr
	conn *net.UDPConn
}

func NewClient() *Client {
	return &Client{}
}

func (client *Client) Connect(address string, port int) error {
	// addr := net.UDPAddr{Port: port, IP: net.ParseIP(address)}
	// conn, err := net.DialUDP("udp", nil, &addr)

	// if err != nil {
	// 	return err
	// }

	// client.addr = addr
	// client.conn = conn

	return nil
}

func (client *Client) Start(recording *recording.Recording) {
	//data := make([]byte, 1024)
	for i := 0; i < runtime.NumCPU(); i++ {

		go func(idx int) {
			addr := net.UDPAddr{Port: 1000 + idx, IP: net.ParseIP("127.0.0.1")}
			conn, err := net.DialUDP("udp", nil, &addr)

			if err != nil {
				log.Fatalln(err)
			}

			defer conn.Close()

			for {
				f, err := recording.ReadFrame()
				if err == io.EOF {
					recording.Reset()
					continue
				}
				_, err = conn.Write(f.Data)
				//writer.Flush()

				if err != nil {
					//fmt.Println(err)
				}
			}
		}(i)
	}
}

func (client *Client) Close() {
	//client.conn.Close()
}
