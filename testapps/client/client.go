package client

import (
	"net"
	"fmt"
)

func RunClientBehavior(address string ) {
	rAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write([]byte("JO"))
	if err != nil {
		fmt.Println(err)
	}
	data := make([]byte, 256)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("received: %s\n", string(data[:n]))
	conn.Close()
}
