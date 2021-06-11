package udp

import (
	"fmt"
	"io"
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
	addr := net.UDPAddr{Port: port, IP: net.ParseIP(address)}
	conn, err := net.DialUDP("udp", nil, &addr)

	if err != nil {
		return err
	}

	client.addr = addr
	client.conn = conn

	return nil
}

func (client *Client) Start(recording *recording.Recording) {
	for i := 0; i < runtime.NumCPU()/2; i++ {

		go func(conn *net.UDPConn) {
			for {
				f, err := recording.ReadFrame()
				if err == io.EOF {
					recording.Reset()
					continue
				}
				_, err = conn.Write(f.Data)

				if err != nil {
					fmt.Println(err)
				}
			}
		}(client.conn)
	}
}

func (client *Client) Close() {
	client.conn.Close()
}
