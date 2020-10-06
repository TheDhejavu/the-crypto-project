![logo](https://github.com/TheDhejavu/the-crypto-project/blob/master/public/cover.png)
# The Crypto Project

This is a blockchain project that implements some of the major feature of popular cryptocurrency project like Bitcoin and ethereum using go programming language. This an experimental project for learning purposes and it contains detailed overview of how blockchain works, **most importantly how this project works**. This project was inspired by [Go Ethereum](https://geth.ethereum.org/docs/) and [bitcoin](https://bitcoin.org)

# Flow Diagram
![flow diagram](https://github.com/TheDhejavu/the-crypto-project/blob/master/public/flowdiagram.jpg)

## Prerequisite
- [Golang](https://golang.org/)
- [libp2p-go ](https://docs.libp2p.io/)
- [BadgerDB](https://github.com/dgraph-io/badger)

## The Codebase
- Blockchain
- Nodes
- Consensus, Blocks & Proof Of Work (POW)
- Wallet
- Transactions
- Uspent Transaction Output (UTXO)
- Merkle Tree
- Networking (P2P/Distributed System)

### Repository Contents

| Folder     | Contents                                         |
|:-----------|:-------------------------------------------------|
| `./p2p`    | Scripts for the crypto project Network Layer    |
| `./Binaries`| An on-demand Folder for executable E.G Wallet   |
| `./cmd`    | CLI Scripts for interacting with the blockhain   |
| `./wallet` | Wallet Source code.                              |
| `./examples`| Code samples for API wrappers written in Go, Python, Typescript, Rust and JS E.T.C                                   |

### Blockchain
Blockchain can be defined as a database that stores blocks, with every next block being linked to the previous one in form of a linked list and a cryptographically secure way so that it’s not possible to change anything in previous blocks, it is a decentralized distributed system of Nodes/Peers that works in a co-ordinated way. What are Distributed systems? this is a computing paradigm whereby two or more nodes work with each other in a coordinated fashion in order to achieve a common outcome and it's modeled in such a way that end users see it as a single logical platform. Distributed system are interestingly complex due to the ability of nodes to coordinate themselves properly. 


### Nodes
Nodes can be defined as any kind of device(mostly computers), phones, laptops, large data centers that uses [graphics processing unit(GPU)](https://en.wikipedia.org/wiki/Graphics_processing_unit) , [Tensor Processing Unit (TPU)](https://en.wikipedia.org/wiki/Tensor_Processing_Unit) E.T.C for expensive and overhead computations. Nodes form the basic infrastructure of a blockchain network, without a node there is no network. All nodes on a blockchain are connected to each other and they constantly exchange the latest blockchain data with each other so all nodes stay up to date. The main purpose of nodes includes but not limited to: storage of blockchain data, verifying of new transactions and blocks, helping of new and existing nodes stay upto date E.t.c. 

This blockchain comprises of three types of nodes whihc includes:

#### Mining Nodes

This is the most important node in the blockchain network, they keep the network running, they foster the minting of new coins, they verify transactions , verifies and mine new blocks.

#### Full Nodes

This type of nodes fosters the verifications of new transactions, manage memory pool (unconfirmed transactions) for miners and also verifies new blocks. 

#### Ordinary Nodes
This type of nodes are part of the network to keep the network running, they mostly verify new blocks on the network

The-crypto-project only supports **1 fullnode** and **1 mining node** with **infinite number of ordinary node** due to underlining issues with the memorypool and mining implementation as explained here

### Consensus mechanism, Blocks & Proof Of Work (POW)
Consensus  mechanism means to reach agreements among network nodes or systems. It fosters consistency of information accross multiple Nodes. Most financial institution today are centralized with lot's of restrictions and regulations, blockchian helps remove that barrier and consensus mechanism is an essential part of the blockchain network  because it allows every nodes in the network to maintain an identical copy of the database. Otherwise, we might end up with conflicting information, undermining the entire purpose of the blockchain network.  Bitcoin was the first cryptocurrency to solve the problem of distributed consensus in a trustless network by using the idea behind [Hashcash](http://www.hashcash.org/). Hashcash is a proof-of-work algorithm, which has been used as a denial-of-service (Dos)counter measure technique in a number of systems. Proof of work fosters minting of new digital currency in blockchain network by allowing Nodes to perfrorm expensive computer calculation, also called **mining**, that needs to be performed in order to create a new group of trustless transactions that forms a **block** on a distributed ledger called **blockchain**. The key purpose of this is to prevent [double spending](https://en.wikipedia.org/wiki/Double-spending), [distributed denial-of-service attack (DDoS)](https://en.wikipedia.org/wiki/Denial-of-service_attack) E.T.C. There are different kinds of consensus mechanism algorithms which work on different principles E.G [Proof of Capacity (POC)](https://www.investopedia.com/terms/c/consensus-mechanism-cryptocurrency.asp) and  [proof of stake (POS)](https://www.investopedia.com/terms/p/proof-stake-pos.asp) but this project implements the Proof of work algorithm used in bitcoin & litecoin

###  Wallet
The wallet system, comparable to a bank account, contains a pair of public and private cryptographic keys. The keys can be used to track ownership, receive or spend cryptocurrencies. A public key allows for other wallets to make payments to the wallet's address, whereas a private key enables the spending of cryptocurrency from that address. 
#### NB: you can't spend your digital currency without your private key and once your private key is compromise, moving your money to a new wallet address is the best thing to do.


The wallet system is independent of the blockchain network and it is built ontop of the `demon` Command line(the network default CLI) and also there is a dedicated executable file in the `binaries` folder coupled with basic commands for performing different actions like generating new wallet, listing existing wallets. 

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

####  Memory pool
This is also know as the transaction pool, this is the waiting area for unconfirmed transactions. When a transaction is carried out by a user, it is sent out to all the avialaible **full nodes** in the network, this full nodes verifies the transaction before adding it to their memory pool while waiting for **mining nodes** to pick it up and includes it in the next block.

From What is the Bitcoin Mempool? A Beginner's Explanation (2020 Updated)
https://99bitcoins.com/bitcoin/mempool/

### Uspent Transaction Output (UTXO) Model

This concept became really popular due to the bitcoin blockhain and this can be defined as an output of a blockchain transaction that has not been spent

They are available to be used in new transactions (as long as you can unlock them with your private key), which makes them useful. UTXOs is used as inputs when a user tries to send X amount of token to Y person given that the amount of UTXOs that the user can unlock is enough to be used as an input. Calculating a wallet address balance can be gotten by accumulating all the unspent transaction outputs that are locked to the particular address

#### Why do we need this ?
 
 Blockchain data are quite verbose, it can range from hundrends to billions of data and computing user wallet balance from a blockchian of that size is computationally expensive in which UTXOs came in as a resucue to reduce overhead. UTXOs ain't all that clever but it's a progress, Ethereum introduced a better way to compute user balance which i think is way better than UTXOs.

 ### How it works
 UTXOs are stored on BadgerDB and specific commands were provided to handle this but Note, UTXOS are created from the blockchain starting from the genesis block and it is computed everytime a new transaction is carried out.


### Merkle Tree

### Networking (P2P/Distributed System)


### Project Setup


### Interfacing with Blockchain Node (JSON-RPC)


Create Wallet

Example 

    curl -X POST -H "Content-Type: application/json" -d '{"id": 1, "method": "API.CreateWallet", "params": []}' http://localhost:5000/_jsonrpc


Get Balance

Example 

    curl -X POST -H "Content-Type: application/json" -d '{"id": 1, "method": "API.GetBalance", "params": [{"Address":"1EWXfMkVj3dAytVuUEHUdoAKdEfAH99rxa"}]}' http://localhost:5000/_jsonrpc



Get Blockchain

Example 

    curl -X POST -H "Content-Type: application/json" -d '{"id": 1,"method": "API.GetBlockchain", "params": []}' http://localhost:5000/_jsonrpc


Get Block by Height

Example 

    curl -X POST -H "Content-Type: application/json" -d '{"id": 1,"method": "API.GetBlockByHeight", "params": ["Height":1]}' http://localhost:5000/_jsonrpc


Send

Example

    curl -X POST -H "Content-Type: application/json" -d '{"id": 1 , "method": "API.Send", "params": [{"sendFrom":"1D214Jcep7x7zPphLGsLdS1hHaxnwTatCW","sendTo": "15ViKshPBH6SzKun1UwmHpbAKD2mKZNtBU", "amount":0.50, "mine": true}]}' http://localhost:5000/_jsonrpc

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

## References && Credits
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
- [Libp2p Overview](https://simpleaswater.com/libp2p-glossary/)
- [bitcoin memory pool](https://99bitcoins.com/bitcoin/mempool/)
- [Big data and business intelligence](https://subscription.packtpub.com/book/big_data_and_business_intelligence/9781787125445/1/ch01lvl1sec8/the-history-of-blockchain)