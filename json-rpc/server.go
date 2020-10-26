package rpc

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/workspace/the-crypto-project/cmd/utils"
	blockchain "github.com/workspace/the-crypto-project/core"
	appUtils "github.com/workspace/the-crypto-project/util/utils"
)

var (
	port = "5000"
)

type API struct {
	RPCEnabled bool
	cmd        *utils.CommandLine
}

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (conn *HttpConn) Read(p []byte) (n int, err error) {
	return conn.in.Read(p)
}
func (conn *HttpConn) Write(d []byte) (n int, err error) {
	return conn.out.Write(d)
}
func (conn *HttpConn) Close() error {
	return nil
}

func (api *API) CreateWallet(args Args, address *string) error {
	*address = api.cmd.CreateWallet()
	return nil
}

func (api *API) GetBalance(args Args, balance *utils.BalanceResponse) error {
	*balance = api.cmd.GetBalance(args.Address)
	return nil
}

func (api *API) GetBlockchain(args Args, data *Blocks) error {
	*data = api.cmd.GetBlockchain()
	return nil
}

func (api *API) GetBlockByHeight(args BlockArgs, data *blockchain.Block) error {
	*data = api.cmd.GetBlockByHeight(args.Height)
	return nil
}

func (api *API) Send(args SendArgs, data *utils.SendResponse) error {
	*data = api.cmd.Send(args.SendFrom, args.SendTo, args.Amount, args.Mine)
	return nil
}

func StartServer(cli *utils.CommandLine, rpcEnabled bool, rpcPort string, rpcAddr string) {
	if rpcPort != "" {
		port = rpcPort
	}

	publicAPI := &API{
		rpcEnabled,
		cli,
	}
	defer cli.Blockchain.Database.Close()
	go appUtils.CloseDB(cli.Blockchain)

	err := rpc.Register(publicAPI)
	checkError("Error registering API", err)
	rpc.HandleHTTP()

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", rpcAddr, port))
	checkError("Listener error:", err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError("Error serving:", err)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		io.WriteString(res, "RPC SERVER LIVE!")
	})
	log.Infof("Serving rpc on port %s", port)

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
					log.Errorf("Error while serving JSON request: %v", err)
					http.Error(w, "Error while serving JSON request", 500)
					return
				}
			}
		}))
	}

}

func checkError(message string, err error) {
	if err != nil {
		log.Info(message, err.Error())
		os.Exit(1)
	}
}
