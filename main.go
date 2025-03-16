package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Start TCP server in a goroutine
	go startTCPServer()

	// Start UDP server in a goroutine
	go startUDPServer()

	// Keep the main function running
	select {}
}

func startTCPServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("TCP server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("TCP Accept error:", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("TCP Read error:", err)
		return
	}
	fmt.Printf("Received TCP message: %s\n", string(buf[:n]))
	RouterDecoder(buf[:n])
	conn.Write([]byte("TCP message received\n"))

}

func startUDPServer() {
	addr := net.UDPAddr{
		Port: 9090,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("UDP server listening on port 9090")

	for {
		buf := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("UDP Read error:", err)
			continue
		}
		fmt.Printf("Received UDP message: %s from %s\n", string(buf[:n]), clientAddr)
		conn.WriteToUDP([]byte("UDP message received\n"), clientAddr)
	}
}
