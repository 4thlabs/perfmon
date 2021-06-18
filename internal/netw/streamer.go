package netw

import (
	"log"
	"net"
	"strconv"
	"time"

	"gitlab.com/4thlabs/perfmon/internal/recording"
)

const NbFrames = 10000

type Streamer struct {
	conn      *net.UDPConn
	recording *recording.Recording
}

func NewStreamer() *Streamer {
	return &Streamer{}
}

func (streamer *Streamer) LoadRecording(file string) error {
	r, err := recording.Open(file)
	if err != nil {
		return err
	}
	streamer.recording = r

	err = streamer.recording.LoadInMemory(NbFrames)
	if err != nil {
		return err
	}

	return nil
}

func (streamer *Streamer) Start(file string, address string, pps int) error {
	err := streamer.LoadRecording(file)
	if err != nil {
		return err
	}

	host, sport, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(sport)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	})

	if err != nil {
		return err
	}

	streamer.conn = conn

	sleepTime := (time.Second.Nanoseconds() / int64(pps))
	ticker := time.NewTicker(time.Duration(sleepTime))
	done := make(chan bool)

	go func() {
		for {
			<-ticker.C
			frame := streamer.recording.GetInMemoryFrame()
			_, err := streamer.conn.Write(frame.Data)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()

	<-done
	return nil
}

func (streamer *Streamer) Close() {
	if streamer.conn != nil {
		streamer.conn.Close()
	}
}
