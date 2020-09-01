![logo](https://github.com/TheDhejavu/the-crypto-project/blob/master/public/cover.png)
# The Crypto Project
This is a blockchain project that implements some of the major feature of popular cryptocurrency project like Bitcoin and ethereum using go programming language.


# High-level Flow Diagram

## PereQuisite
- [Golang](https://golang.org/)
Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.
- [libp2p-go ](https://docs.libp2p.io/)
- [BadgerDB](https://github.com/dgraph-io/badger)
BadgerDB is an embeddable, persistent and fast key-value (KV) database written in pure Go. 

## Terms
- Blockchain
- Consensus, Blocks & Proof Of Work (POW)
- Wallet
- Transactions
- Uspent Transaction Output (UTXO)
- Merkle Tree
- Networking (P2P/Distributed System)

### Blockchain

#### Interacting with the blockchain
    - Via CLI

    - REStFul API

### Consensus, Blocks & Proof Of Work (POW)


###  Wallet
The wallet system, comparable to a bank account, contains a pair of public and private cryptographic keys. The keys can be used to track ownership, receive or spend cryptocurrencies. A public key allows for other wallets to make payments to the wallet's address, whereas a private key enables the spending of cryptocurrency from that address.
## NB: you can't spend your digital currency without your private key and once your private key is compromise, moving your money to a new wallet address is the best thing to do. 

### Transactions


### Uspent Transaction Output (UTXO) Model
    This concept is particularly famous with the bitcoin blockhain which is defined as an output of a blockchain transaction that has not been spent

### Merkle Tree

### Networking (P2P/Distributed System)


## TODO
- Data visualization
- Node visualization
- gRPC implementation for accessing blockchain data

## References
- [MIT 6.824 Distributed Systems (Spring 2020)](https://www.youtube.com/playlist?list=PLrw6a1wE39_tb2fErI4-WkMbsvGQk9_UB)
- [Code a simple P2P blockchain in Go!](https://medium.com/@mycoralhealth/code-a-simple-p2p-blockchain-in-go-46662601f417)
- [Advanced Blockchain Concepts for Beginners](https://medium.com/@mycoralhealth/advanced-blockchain-concepts-for-beginners-32887202afad)
- [Tensor go programming ](https://www.youtube.com/playlist?list=PLJbE2Yu2zumCe9cO3SIyragJ8pLmVv0z9)
- [MerkleTree](https://brilliant.org/wiki/merkle-tree/)
- [BadgerDB](https://github.com/dgraph-io/badger)
- [Waller](https://en.wikipedia.org/wiki/Cryptocurrency_wallet)