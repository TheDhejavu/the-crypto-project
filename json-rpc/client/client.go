package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"

	rpc "github.com/workspace/the-crypto-project/json-rpc"
)

func main() {
	client, err := jsonrpc.Dial("tcp", "localhost:5000")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args := rpc.Args{
		Address: "14RwDN6Pj4zFUzdjiB8qUkVMC1QvRG5Cmr",
	}
	var bs rpc.Blocks
	err = client.Call("API.GetBlockchainData", args, &bs)
	if err != nil {
		log.Fatal("API error:", err.Error())
	}
	for _, block := range bs {
		fmt.Printf("%x", block.PrevHash)
	}
}
