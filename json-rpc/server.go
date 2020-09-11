package rpc

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	"github.com/workspace/the-crypto-project/cmd/utils"
)

type API struct {
	RPCEnabled bool
	cmd        *utils.CommandLine
}

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

func (api *API) GetBalance(args Args, balance *string) error {
	*balance = api.cmd.GetBalance(args.Address)
	return nil
}

func (api *API) CreateWallet(args Args, address *string) error {
	*address = api.cmd.CreateWallet()
	return nil
}

func (api *API) GetBlockchainData(args SendArgs, data *Blocks) error {
	*data = api.cmd.GetBlockchainData()
	return nil
}

func (api *API) Send(args SendArgs, data *string) error {
	*data = api.cmd.Send(args.sendFrom, args.sendTo, args.amount, args.mine)
	return nil
}

func StartServer(rpcEnabled bool, rpcPort int, rpcAddr string) {

	publicAPI := &API{
		rpcEnabled,
		&utils.CommandLine{},
	}
	err := rpc.Register(publicAPI)
	checkError("Error registering API", err)
	rpc.HandleHTTP()

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":5000")
	checkError("Listener error:", err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError("Error serving:", err)

	// sample test endpoint
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		io.WriteString(res, "RPC SERVER LIVE!")
	})
	log.Printf("Serving rpc on port %d", 5000)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
		http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Path == "/_jsonrpc" {
				serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(200)
				err := rpc.ServeRequest(serverCodec)
				if err != nil {
					log.Printf("Error while serving JSON request: %v", err)
					http.Error(w, "Error while serving JSON request, details have been logged.", 500)
					return
				}
			}
		}))
	}

}

func checkError(message string, err error) {
	if err != nil {
		fmt.Println(message, err.Error())
		os.Exit(1)
	}
}
