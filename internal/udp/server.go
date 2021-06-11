package udp

import (
	"fmt"
	"net"
	"runtime"
)

type Server struct {
	addr net.UDPAddr
	conn net.PacketConn
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Connect(address string, port int) error {
	addr := net.UDPAddr{Port: port, IP: net.ParseIP(address)}
	conn, err := net.ListenPacket("udp", fmt.Sprintf("%s:%d", address, port))

	if err != nil {
		return err
	}

	server.addr = addr
	server.conn = conn

	return nil
}

func (server *Server) Start() {
	for i := 0; i <= runtime.NumCPU()/2; i++ {
		go func(conn net.PacketConn, addr net.UDPAddr) {
			buffer := make([]byte, 2048)
			for {
				//fmt.Println("Reading")
				_, _, err := conn.ReadFrom(buffer)
				if err != nil {
					fmt.Println(err)
				}
			}
		}(server.conn, server.addr)
	}
}

func (server *Server) Close() {
	server.conn.Close()
}
