

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	protocol := flag.String("protocol", "tcp", "Protocol to use: tcp or udp")
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "8080", "Server port")
	message := flag.String("message", "Hello, Server!", "Message to send")
	flag.Parse()

	addr := fmt.Sprintf("%s:%s", *host, *port)

	if *protocol == "tcp" {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("Error connecting to TCP server:", err)
			os.Exit(1)
		}
		defer conn.Close()

		_, err = conn.Write([]byte(*message))
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

		_, err = conn.Write([]byte(*message))
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
