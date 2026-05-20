package server

import (
	"fmt"
	"net"
	"os"
)

func StartUDP(port int, logChannel chan<- string) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("ERROR serving UDP server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Server UDP listening on port %d\n", port)
	buffer := make([]byte, 2048)

	for {
		bytes, remoteAddr, _ := conn.ReadFromUDP(buffer)
		msg := string(buffer[:bytes])
		logEntry := fmt.Sprintf("[%s]: %s", remoteAddr, msg)
		logChannel <- logEntry
	}
}

func StartTCP(port int, logChannel chan<- string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("ERROR serving UDP server: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Server TCP listening on port %d\n", port)
	for {
		conn, _ := listener.Accept()
		go func(c net.Conn) {
			defer c.Close()
			buffer := make([]byte, 2048)
			bytes, _ := c.Read(buffer)
			logChannel <- string(buffer[:bytes])
		}(conn)
	}
}
