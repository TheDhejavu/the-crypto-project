package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Write([]byte("Hello server"))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Message sent")

	for {
		buffer := make([]byte, 1400)
		size, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Connection closed")
			return
		}

		data := buffer[:size]
		fmt.Printf("recieved message %s", data)
	}
}
