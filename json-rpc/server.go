package rpc

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	"github.com/workspace/the-crypto-project/cmd/utils"
)

type PublicCryptoAPI struct {
	RPCEnabled bool
	cmd        *utils.CommandLine
}

func (api *PublicCryptoAPI) GetBalance(address *string, balance *string) error {
	*balance = api.cmd.GetBalance(*address)
	return nil
}

func StartServer(rpcEnabled bool, rpcPort int, rpcAddr string) {

	publicAPI := &PublicCryptoAPI{
		rpcEnabled,
		&utils.CommandLine{},
	}
	err := rpc.Register(publicAPI)
	checkError("Error registering API", err)

	rpc.HandleHTTP()

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError("Listener error:", err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError("Error serving:", err)

	http.Serve(listener, nil)
	log.Printf("Serving rpc on port %d", 1234)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
	}

}

func checkError(message string, err error) {
	if err != nil {
		fmt.Println(message, err.Error())
		os.Exit(1)
	}
}
