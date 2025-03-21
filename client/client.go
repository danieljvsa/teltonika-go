package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	//exempCodec8 := "000000000000003608010000016B40D8EA300100000000000000000000000105021503010101425E0F01F10000601A014E0000000000000000010000C7CF"
	exempCodec82 := "000000000000002A08010000018F13F93C000008F66E6C01692CB600C8012C0A00320001848"
	protocol := flag.String("protocol", "tcp", "Protocol to use: tcp or udp")
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "8080", "Server port")
	message := flag.String("message", exempCodec82, "Message to send")
	flag.Parse()

	addr := fmt.Sprintf("%s:%s", *host, *port)

	rawBytes, err := hex.DecodeString(*message)
	if err != nil {
		fmt.Println("Error decoding hex:", err)
		return
	}

	if *protocol == "tcp" {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("Error connecting to TCP server:", err)
			os.Exit(1)
		}
		defer conn.Close()

		_, err = conn.Write([]byte(rawBytes))
		if err != nil {
			fmt.Println("Error sending message:", err)
			os.Exit(1)
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading response:", err)
			os.Exit(1)
		}

		fmt.Println("Server response:", string(buffer[:n]))
	} else if *protocol == "udp" {
		conn, err := net.Dial("udp", addr)
		if err != nil {
			fmt.Println("Error connecting to UDP server:", err)
			os.Exit(1)
		}
		defer conn.Close()

		_, err = conn.Write([]byte(rawBytes))
		if err != nil {
			fmt.Println("Error sending message:", err)
			os.Exit(1)
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading response:", err)
			os.Exit(1)
		}

		fmt.Println("Server response:", string(buffer[:n]))
	} else {
		fmt.Println("Invalid protocol. Use 'tcp' or 'udp'")
		os.Exit(1)
	}
}
