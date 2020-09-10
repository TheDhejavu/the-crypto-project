package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var balance string
	err = client.Call("PublicCryptoAPI.GetBalance", "14RwDN6Pj4zFUzdjiB8qUkVMC1QvRG5Cmr", &balance)
	if err != nil {
		log.Fatal("API error:", err)
	}
	fmt.Printf("BALANCE: %s", balance)
}
