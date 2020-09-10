package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
)

type Args struct {
	Address string
}

func main() {
	client, err := jsonrpc.Dial("tcp", "localhost:5000")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args := Args{
		Address: "14RwDN6Pj4zFUzdjiB8qUkVMC1QvRG5Cmr",
	}
	var balance string
	err = client.Call("API.GetBalance", args, &balance)
	if err != nil {
		log.Fatal("API error:", err.Error())
	}
	fmt.Printf("BALANCE: %s", balance)
}
