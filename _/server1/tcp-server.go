package main

import (
	"fmt"
	"log"
	"net"
)

func listenConnection(conn net.Conn) {
	for {
		buffer := make([]byte, 1400)
		size, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Connection closed")
			return
		}

		data := buffer[:size]
		fmt.Printf("recieved message %s", data)

		_, err = conn.Write(data)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
func main() {
	fmt.Println("Listening to localhost:5000")
	listener, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		log.Fatalln(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("New connection")
		go listenConnection(conn )
	}
}
