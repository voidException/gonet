package main

import (
	"io"
	"log"
	"net"
	"os"
)

import (
	"cfg"
	. "db"
	"hub/protos"
)

//----------------------------------------------- HUB start
func main() {
	log.Println("Starting HUB")

	// start db
	StartDB()

	// data init
	startup_work()

	// Listen
	service := ":9090"
	config := cfg.Get()

	if config["hub_service"] != "" {
		service = config["hub_service"]
	}

	log.Println("Hub Service:", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

//----------------------------------------------- handle logical server
func handleClient(conn net.Conn) {
	defer conn.Close()

	header := make([]byte, 2)
	ch := make(chan []byte, 8192)

	go protos.HubAgent(ch, conn)

	for {
		// header
		n, err := io.ReadFull(conn, header)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			log.Println("error receving header:", err)
			break
		}

		// data
		size := int(header[0])<<8 | int(header[1])
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		if err != nil {
			log.Println("error receving msg:", err)
			break
		}
		ch <- data
	}

	close(ch)
}

func checkError(err error) {
	if err != nil {
		log.Println("Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
