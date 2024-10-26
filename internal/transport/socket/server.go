package socket_server

import (
	"fmt"
	"net"
)

func NewServer() {
	p := make([]byte, 2048)

	addr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8080,
	}

	ser, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
	}
}
