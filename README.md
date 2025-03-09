# Teltonika TCP/UDP Server

## Overview

This project presents a **TCP and UDP server** developed in Go, specifically designed for interfacing with **Teltonika GPS devices**. The server facilitates bidirectional communication by accepting incoming connections, parsing data packets in accordance with Teltonika's proprietary protocol, processing the extracted information, and formulating appropriate responses.

## Key Features

- Supports **TCP (8080)** and **UDP (9090)** network communication protocols.
- Integrates seamlessly with **Teltonika GPS trackers**.
- Implements **packet decoding and encoding** based on the [Teltonika Data Sending Protocols](https://wiki.teltonika-gps.com/view/Teltonika_Data_Sending_Protocols).
- Easily deployable via **Docker** and **Docker Compose**.
- Includes a **client utility** to send test messages to both TCP and UDP servers.

## Project Structure

```
my-go-server/
│── main.go               # Core server logic implemented in Go
│── client.go             # Client implementation for testing TCP/UDP communication
│── Dockerfile            # Configuration for Docker containerization
│── docker-compose.yml    # Definition file for Docker Compose
│── README.md             # Documentation and setup guide
│── go.mod                # Go module configuration
│── go.sum                # Dependency checksum tracking
```

---

## Getting Started

### Prerequisites

Before proceeding with the installation, ensure the following dependencies are installed on your system:

- [Go](https://go.dev/dl/) (latest stable version recommended)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation Steps

#### 1. Clone the Repository

```sh
git clone https://github.com/yourusername/my-go-server.git
cd my-go-server
```

#### 2. Install Required Dependencies

```sh
go mod tidy
```

#### 3. Launch the Server

**Run the server locally:**

```sh
go run main.go
```

**Alternatively, run the server in a Docker container:**

```sh
docker-compose up --build
```

---

## Testing the Server

### Verifying TCP Communication

To test TCP connectivity, use **netcat**, a Teltonika GPS tracker, or the provided client:

```sh
nc -v localhost 8080
```

Alternatively, use the test client:

```sh
go run client.go -protocol tcp -message "Hello TCP Server" -host localhost -port 8080
```

### Verifying UDP Communication

To test UDP functionality, execute the following command:

```sh
echo -n "Hello UDP" | nc -u -w1 localhost 9090
```

Or use the test client:

```sh
go run client.go -protocol udp -message "Hello UDP Server" -host localhost -port 9090
```

---

## Teltonika Data Processing

This server is engineered to process **Teltonika AVL (Automatic Vehicle Location) data packets**, as defined in the [Teltonika Protocol](https://wiki.teltonika-gps.com/view/Teltonika_Data_Sending_Protocols). It performs the following operations:

- **Decodes AVL packets** received from Teltonika tracking devices.
- **Extracts and processes** GPS and IO (input/output) telemetry data.
- **Encodes and transmits acknowledgments** back to the originating device to confirm receipt of data.

### Data Flow Example

1. **A Teltonika device establishes a connection** via TCP or UDP.
2. **The server captures and reads the AVL data packet**.
3. **The server decodes and processes GPS and telemetry data**.
4. **An acknowledgment response is transmitted back to the device**.

---

## Deployment Instructions

To deploy the server using Docker in a production environment, execute:

```sh
docker-compose up -d
```

This command initializes the server in detached mode, ensuring continuous operation in the background.

---

## Contributing to the Project

We welcome contributions! If you would like to improve the server's functionality, particularly in handling Teltonika-specific data structures, feel free to submit **issues** or **pull requests**.

---

## License

This project is distributed under the **MIT License**.

---

## References

- [Teltonika Data Sending Protocols](https://wiki.teltonika-gps.com/view/Teltonika_Data_Sending_Protocols)
- [Go net package documentation](https://pkg.go.dev/net)

