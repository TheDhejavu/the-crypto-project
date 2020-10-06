![logo](https://github.com/TheDhejavu/the-crypto-project/blob/master/public/cover.png)
# The Crypto Project

This is a blockchain project that implements some of the major feature of popular cryptocurrency project like Bitcoin and ethereum using go programming language. This an experimental project for learning purposes and it contains detailed overview of how blockchain works, **most importantly how this project works**. This project was inspired by [Go Ethereum](https://geth.ethereum.org/docs/) and [bitcoin](https://bitcoin.org)

# Flow Diagram
![flow diagram](https://github.com/TheDhejavu/the-crypto-project/blob/master/public/flowdiagram.jpg)

## Prerequisite
- Programming language: [Golang](https://golang.org/)
- Networking: [libp2p-go ](https://docs.libp2p.io/)
- Database: [BadgerDB](https://github.com/dgraph-io/badger) 

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
Nodes can be defined as any kind of device(mostly computers), phones, laptops, large data centers that uses [graphics processing unit(GPU)](https://en.wikipedia.org/wiki/Graphics_processing_unit) , [Tensor Processing Unit (TPU)](https://en.wikipedia.org/wiki/Tensor_Processing_Unit) E.T.C for expensive and overhead computations. Nodes form the basic infrastructure of a blockchain network, without a node there is no network. All nodes on a blockchain are connected to each other and they constantly exchange the latest blockchain data with each other so that all nodes stay up to date. The main purpose of nodes includes but not limited to: storage of blockchain data, verifying of new transactions and blocks, helping of new and existing nodes stay upto date E.t.c. 

This blockchain comprises of three types of nodes which includes:
#### Mining Nodes

This is the most important node in the blockchain network, they keep the network running, they foster the minting of new coins, they verify transactions , verifies and mine new blocks.

#### Full Nodes

This type of nodes fosters the verifications of new transactions, manage memory pool (unconfirmed transactions) for miners and also verifies new blocks. 

#### Ordinary Nodes
This type of nodes are part of the network to keep the network running, they mostly verify new blocks on the network

The-crypto-project only supports **1 fullnode** and **1 mining node** with **infinite number of ordinary node** due to underlining issues with the memorypool and the mining implementation as explained here

### Consensus mechanism,Mining, Blocks & Proof Of Work (POW)
Consensus  mechanism means to reach agreements among network nodes or systems. It fosters consistency of information accross multiple Nodes. Most financial institution today are centralized with lot's of restrictions and regulations, blockchian helps remove that barrier and consensus mechanism is an essential part of the blockchain network  because it allows every nodes in the network to maintain an identical copy of the database. Otherwise, we might end up with conflicting information, undermining the entire purpose of the blockchain network.  Bitcoin was the first cryptocurrency to solve the problem of distributed consensus in a trustless network by using the idea behind [Hashcash](http://www.hashcash.org/). Hashcash is a proof-of-work algorithm, which has been used as a denial-of-service (Dos)counter measure technique in a number of systems. Proof of work fosters minting of new digital currency in blockchain network by allowing Nodes to perfrorm expensive computer calculation, also called **mining**, that needs to be performed in order to create a new group of trustless transactions that forms a **block** on a distributed ledger called **blockchain**. The key purpose of this is to prevent [double spending](https://en.wikipedia.org/wiki/Double-spending), [distributed denial-of-service attack (DDoS)](https://en.wikipedia.org/wiki/Denial-of-service_attack) E.T.C. There are different kinds of consensus mechanism algorithms which work on different principles E.G [Proof of Capacity (POC)](https://www.investopedia.com/terms/c/consensus-mechanism-cryptocurrency.asp) and  [proof of stake (POS)](https://www.investopedia.com/terms/p/proof-stake-pos.asp) but this project implements the Proof of work algorithm used in bitcoin & litecoin

#### How we know that a block is valid ?
We basically check for two things.
1. We Check if the previous block referenced by the block exists and is valid.

2. We Check that the proof of work done on the block is valid.

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
A transaction is a transfer of value between wallets that gets included in the block chain as defined by [bitcoin.org](https://bitcoin.org/en/how-it-works). It comprises of the transaction Inputs and outputs, the transaction inputs comprises of an array of spent coins gotten from the outputs  while transaction outputs comprises of unspent coins. This transactions are signed with a secret called private key that can be found in the user wallet to proof that a user is indeed the owner of the coins, this transactions is initialized and sent to the network which in turn under-goes a series of verification by the network nodes to confirm the validity of the transaction using the user's public key.

####  Memory pool
This is also know as the transaction pool, this is the waiting area for unconfirmed transactions. When a transaction is carried out by a user, it is sent out to all the avialaible **full nodes** in the network, this full nodes verifies the transaction before adding it to their memory pool while waiting for **mining nodes** to pick it up and includes it in the next block.

 [What is the Bitcoin Mempool? A Beginner's Explanation (2020 Updated)](https://99bitcoins.com/bitcoin/mempool/)

### Uspent Transaction Output (UTXO) Model

This concept became really popular due to the bitcoin blockhain and this can be defined as an output of a blockchain transaction that has not been spent

They are available to be used in new transactions (as long as you can unlock them with your private key), which makes them useful. UTXOs is used as inputs when a user tries to send X amount of coin to Y person given that the amount of UTXOs that the user can unlock is enough to be used as an input. Calculating a wallet address balance can be gotten by accumulating all the unspent transaction outputs that are locked to the particular address

#### Why do we need this ?
 
 Blockchain data are quite verbose, it can range from hundrends to billions of data and computing user wallet balance from a blockchian of that size is computationally expensive in which UTXOs came in as a rescue to reduce overhead. UTXOs ain't all that clever but it's a progress, and every idea has it's tradeoff's. [Ethereum introduced it's own way to compute user balance ](https://github.com/ethereum/wiki/wiki/Design-Rationale#accounts-and-not-utxos)

 ### How it works (the-crypto-project context)
 UTXOs are stored on BadgerDB and specific commands were provided to handle this but Note, UTXOS are created from the blockchain starting from the genesis block and it is computed everytime a new transaction is carried out and when a new block is added as opposed everytime a user checks his/her balance.


### Merkle Tree

A Merkle tree can be simply defined as a binary hash tree data structure , composed of a set of nodes with a large number of leaf nodes at the bottom of the tree containing the underlying data, a set of intermediate nodes where each node is the hash of its two children, and finally a single root node, also formed from the hash of its two children, representing the "top" of the tree called the merkle root, which enables the quick verification of blockchain data, as well as quick movement of large amounts of data from one computer node to the other on the blockchain network. The transactions are executed on a merkle tree algorithm to generate a single hash which is a string of numbers and letters that can be used to verify that a given set of data is the same as the original set of transactions. 

### Networking (Peer-2-Peer)

The Blockchain protocol operates on top of the Internet, on a P2P network of computers that run the protocol and hold an identical copy of the ledger of transactions , enabling P2P value transactions without a middleman through consensus mechanism. In computing, p2p is a network of peers that enables the storage and sharing of files with equal power( might differs in terms of computation) and functionality. They can act as both a client and a server, exchanging information in realtime and when a node acts as a client, they download files from other network nodes. But when they are working as a server, they are the source from which other nodes can download files. P2P networks don’t have a single point of failure and enables a system to continue operating properly in the event of the failure know as **fault tolerance**. The P2P network is an essential part of the blockchain network because it allows the distribution of blockchain data across mulitple node/peer, it prevents the Denial-of-Service (DoS) attacks that plague numerous systems, and also renders them resistant to censorship by central authorities. The major limitation of P2p is the ability to maintain consistent in data across all peers (subjective) and also proof of work is way too computationally expensive for a less powerful computer and this will only get worse as the Blockchain gets bigger with an increase in difficulty which means nodes who has less computational power to participate will eventually leave, but on the bright side, P2P makes decentralization possible and provides overall security for the blockchain.

the-crypto-project achieved 100% decentralization via the use of  [libp2p-go ](https://docs.libp2p.io/) networking libraries used by popular project like [Ipfs](https://ipfs.io/), [filecoin ](https://filecoin.io/) and most recently Ethereum 2.0.

### Project Setup


### Node JSON-RPC server


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

- Web Application
- Basic Virtual Machine (Maybe) 
- Test Coverage
- Improved Error Handling

## References and Credits
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
- [Merkle Tree by investopedia](https://www.investopedia.com/terms/m/merkle-tree.asp)
- [P2P Explained](https://academy.binance.com/en/articles/peer-to-peer-networks-explained)
- [Big data and business intelligence](https://subscription.packtpub.com/book/big_data_and_business_intelligence/9781787125445/1/ch01lvl1sec8/the-history-of-blockchain)