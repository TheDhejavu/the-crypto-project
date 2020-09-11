![logo](https://github.com/TheDhejavu/the-crypto-project/blob/master/public/cover.png)
# The Crypto Project

This is a blockchain project that implements some of the major feature of popular cryptocurrency project like Bitcoin and ethereum using go programming language. This an experimental project for learning purposes and it contains detailed overview of how blockchain works, most importantly how this project works. You can checkout references to some major materials that plays important roles in making this project happen.

# Flow Diagram
![flow diagram](https://github.com/TheDhejavu/the-crypto-project/blob/master/public/flowdiagram.jpg)

## Prerequisite
- [Golang](https://golang.org/)
- [libp2p-go ](https://docs.libp2p.io/)
- [BadgerDB](https://github.com/dgraph-io/badger)

## The Codebase
- Blockchain
- Consensus, Blocks & Proof Of Work (POW)
- Wallet
- Transactions
- Uspent Transaction Output (UTXO)
- Merkle Tree
- Networking (P2P/Distributed System)

### Repository Contents

| Folder     | Contents                                         |
|:-----------|:-------------------------------------------------|
| `./p2p`    | Scripts for the crypto project nNetwork Layer    |
| `./Binaries`| An on-demand Folder for executable E.G Wallet   |
| `./cmd`    | CLI Scripts for interacting with the blockhain   |
| `./wallet` | Wallet Source code.                              |
| `./examples`| Code samples for API wrappers written in Go, Python and NodeJS                                    |

### Blockchain
Blockchain can be defined as a database that stores blocks, with every next block being linked to the previous one in a cryptographically secure way so that it’s not possible to change anything in previous blocks, it is a Decentralized system of Nodes that works in a co-ordinated way.

The cryto-project blockchain has different actions which are made accessible via the in-built CLI and the REST API

### Consensus, Blocks & Proof Of Work (POW)

###  Wallet
The wallet system, comparable to a bank account, contains a pair of public and private cryptographic keys. The keys can be used to track ownership, receive or spend cryptocurrencies. A public key allows for other wallets to make payments to the wallet's address, whereas a private key enables the spending of cryptocurrency from that address.
#### NB: you can't spend your digital currency without your private key and once your private key is compromise, moving your money to a new wallet address is the best thing to do.


The wallet system is independent of the blockchain network and it is built ontop of the `demon` Command line(the network default CLI) and also there is a dedicated executable file in the `binaries` folder coupled with basic commands for performing different actions

##### Download https://github.com/TheDhejavu/the-crypto-project/tree/master/binaries/wallet.exe

#### Via Standalone Binaries
#### Commands 

Generate new wallet

    ./wallet new

Print all local wallet

    ./wallet print

Print wallet by Address

    ./wallet print --address ADDRESS

### Transactions


### Uspent Transaction Output (UTXO) Model

This concept became really popular due to the bitcoin blockhain and this can be defined as an output of a blockchain transaction that has not been spent

They are available to be used in new transactions (as long as you can unlock them with your private key), which makes them useful. UTXOs is used as inputs when a user tries to send X amount of token to Y person given that the amount of UTXOs that the user can unlock is enough to be used as an input. Calculating a wallet address balance can be gotten by accumulating all the unspent transaction outputs that are locked to the particular address
 #### Why do we need this ?
 
 Blockchain data are quite verbose, it can range from hundrends to billions of data and computing user wallet balance from a blockchian of that size is computationally expensive in which UTXOs came in as a resucue to reduce overhead. UTXOs ain't all that clever but it's a progress, Ethereum introduced a better way to compute user balance which i think is way better than UTXOs.

 ### How it works
 UTXos are stored on BadgerDB and specific commands were provided to handle this but Note, UTXos are created from the blockchain

### Merkle Tree

### Networking (P2P/Distributed System)


### Project Setup


### Interfacing with Blockchain Nodes (JSON-RPC)


Get Balance

| Type     | Content                                        |
|:-----------|:-------------------------------------------------|
| Go    | GetBalance(address string) string {}  |
| JSON-RPC | {"method":"API.GetBalance","params":["14RwDN6Pj4zFUzdjiB8qUkVMC1QvRG5Cmr"] |

Example 

    curl -X POST -H "Content-Type: application/json" -d '{"id": 1, "method": "API.GetBalance", "params": [{"Address":"1EWXfMkVj3dAytVuUEHUdoAKdEfAH99rxa"}]}' http://localhost:5000/_jsonrpc


Create Wallet

| Type     | Content                                        |
|:-----------|:-------------------------------------------------|
| Go    | CreateWalet() string {}  |
| JSON-RPC | {"method":"API.GetBalance","params":["14RwDN6Pj4zFUzdjiB8qUkVMC1QvRG5Cmr"] |

Example 

    curl -X POST -H "Content-Type: application/json" -d '{"id": 1, "method": "API.CreateWallet", "params": []}' http://localhost:5000/_jsonrpc


Get Blockchain Data

| Type     | Content                                        |
|:-----------|:-------------------------------------------------|
| Go    | GetBlockchainData() string {}  |
| JSON-RPC | {"method": "API.GetBlockchainData", "params": [] |

Example 

    curl -X POST -H "Content-Type: application/json" -d '{"id": 1,"method": "API.GetBlockchainData", "params": []}' http://localhost:5000/_jsonrpc


Send coins

| Type     | Content                                        |
|:-----------|:-------------------------------------------------|
| Go    | Send(from string, to string, amount float64, mineNow bool )string  {} |
| JSON-RPC | {"method": "API.Send", "params": [] |

Example

   curl -X POST -H "Content-Type: application/json" -d '{"id": 1 ,"method": "API.Send", "params": [{"sendTo":"1D214Jcep7x7zPphLGsLdS1hHaxn", "sendFrom":"15ViKshPBH6SzKun1UwmHpbAKD2mKZNtBU", "amount": 50,"mine": true}]}' http://localhost:5000/_jsonrpc

### Demon CLI

This is the official command line for the crypto project, this commandline allows developers to interact with the blockchain network

##### CLI https://github.com/TheDhejavu/the-crypto-project/tree/master/cmd/chain

#### Commands 

Generate new wallet

    demon wallet new

List Addresses

    demon wallet listaddress

Get Balance

    demon wallet balance --address ADDRESS

Print blockchain

    demon printblockchain
    
Compute UTXOs

    demon computeutxos

Send

    demon send --sendFrom ADDRESS --sendTo ADDRESS --amount AMOUNT 

Start Node

The minerAddress, miner and ListenPort Flags are optional if this flags already exist in `.env` file

    demon startnode --ListenPort PORT --minerAddress MINER_ADDRESS --miner

#### Command Usage

    Usage:
        demon [flags]
        demon [command]

    Available Commands:
        computeutxos Re-build and Compute Unspent transaction outputs
        help         Help about any command
        init         Initialize the blockchain and create the genesis block
        print        Print the blocks in the blockchain
        send         Send x amount of token to address from local
        wallet address
        startnode    start a node
        wallet       Manage wallets

    Flags:
            --address string   Wallet address
        -h, --help             help for demon
            --rpc              Enable the HTTP-RPC server
            --rpcAddr string   HTTP-RPC server listening interface (default: localhost)
            --rpcPort int       HTTP-RPC server listening port (default: 1245)

    Use "demon [command] --help" for more information about a command.

## TODO

- Data visualization
- Smart Contract & Basic Virtual Machine (Maybe, 😃)
- Write Test

## References
- [Blockchain Basic, A Non-technical guide](https://www.goodreads.com/book/show/34137265-blockchain-basics)
- [MIT 6.824 Distributed Systems (Spring 2020)](https://www.youtube.com/playlist?list=PLrw6a1wE39_tb2fErI4-WkMbsvGQk9_UB)
- [Code a simple P2P blockchain in Go!](https://medium.com/@mycoralhealth/code-a-simple-p2p-blockchain-in-go-46662601f417)
- [Advanced Blockchain Concepts for Beginners](https://medium.com/@mycoralhealth/advanced-blockchain-concepts-for-beginners-32887202afad)
- [Tensor go programming ](https://www.youtube.com/playlist?list=PLJbE2Yu2zumCe9cO3SIyragJ8pLmVv0z9)
- [MerkleTree](https://brilliant.org/wiki/merkle-tree/)
- [BadgerDB](https://github.com/dgraph-io/badger)
- [Wallet](https://en.wikipedia.org/wiki/Cryptocurrency_wallet)
- [How Bitcoin Works under the Hood](https://www.youtube.com/watch?v=Lx9zgZCMqXE)
- [What is the difference between decentralized and distributed systems?](https://medium.com/distributed-economy/what-is-the-difference-between-decentralized-and-distributed-systems-f4190a5c6462#:~:text=A%20decentralized%20system%20is%20a%20subset%20of%20a%20distributed%20system.&text=Decentralized%20means%20that%20there%20is,where%20the%20decision%20is%20made.&text=Distributed%20means%20that%20the%20processing,and%20use%20complete%20system%20knowledge.)
- [Ethereum Geth](https://geth.ethereum.org/docs/)
- [Original Bitcoin Client](https://en.bitcoin.it/wiki/Original_Bitcoin_client/API_calls_list)