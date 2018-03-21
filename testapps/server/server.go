package server

import (
	"net"
	"fmt"
	"time"
)

func RunServerBehavior(listenAddress string) {
	addr, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	time.Sleep(time.Second * 2)
	fmt.Print("Listening... ")

	conn, err := listener.AcceptTCP()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	data := make([]byte, 256)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("received: %s\n", string(data[:n]))

	conn.Write([]byte("HO"))

	}
