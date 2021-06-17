package main

import (
	"fmt"
	"io"
	"log"

	"gitlab.com/4thlabs/perfmon/internal/netw"
	"gitlab.com/4thlabs/perfmon/internal/recording"
	"gitlab.com/4thlabs/perfmon/internal/ui"
)

func ReadRecording() {
	r, err := recording.Open("6cb3b905-8d80-3e5d-8178-1967915f8849_new_1")
	if err != nil {
		panic(err)
	}

	defer r.Close()

	for i := 0; i < 10000000; i++ {
		f, err := r.ReadFrame()
		if err != nil {
			if err == io.EOF {
				r.Reset()
				continue
			} else {
				panic(err)
			}
		}

		fmt.Printf("Packet length: %d\n", f.Length)
	}
}

func main() {
	r, err := recording.Open("6cb3b905-8d80-3e5d-8178-1967915f8849_new_1")
	if err != nil {
		log.Panicln(err)
	}

	err = r.LoadInMemory(10000)
	if err != nil {
		log.Panicln(err)
	}

	defer r.Close()

	client := netw.NewClient()
	server := netw.NewServer()

	if err := client.Connect("127.0.0.1", 1237); err != nil {
		panic(err)
	}

	if err := server.Connect("127.0.0.1", 1237); err != nil {
		panic(err)
	}

	defer client.Close()
	defer server.Close()

	//	server.Start()
	client.Start(r)

	// for {
	// 	stats, err := stats.Get()
	// 	if err == nil {
	// 		fmt.Printf("User CPU %f \n", stats.Cpu.User)
	// 		fmt.Printf("System CPU %f \n", stats.Cpu.System)
	// 		fmt.Printf("Bytes Sent %d \n", stats.Network.RxBytes)
	// 		fmt.Printf("Bytes Received %d \n", stats.Network.TxBytes)

	// 		fmt.Printf("Packets Sent %d \n", stats.Network.RxPackets)
	// 		fmt.Printf("Packets Received %d \n", stats.Network.TxPackets)
	// 	}
	// }

	ui.Init()
}
