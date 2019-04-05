# Quick Start

This tutorial will help you setup and start Electionr blockchain.

We will start the blockchain with 2 validator nodes, adding new election in first election cycle.

### Prerequisites

- `Golang v1.11+`

### Installation

- Navigate to your electionr cloned repository
- Install Golang packages: `go mod tidy`
- Install Electionr: `make install`

### Initialize nodes configuration files

- First we need to initialize configuration files for each node:
  - `electionrd init node0 --home ~/.electionrd/node0 --chain-id electionr-chain`
  - `electionrd init node1 --home ~/.electionrd/node1 --chain-id electionr-chain`
  - `electionrd init node2 --home ~/.electionrd/node2 --chain-id electionr-chain`

- Replace generated genesis.json files in each node's configuration directory with genesis.json from testnet_config

- Now recover node keys via mnemonics:
  - `electionrcli keys add operator0 --recover`
  - `electionrcli keys add operator1 --recover`
  - `electionrcli keys add operator2 --recover`  
Note: After each command you will be prompted for mnemonic.  
Enter appropriate mnemonic from mnemonics directory.

### Start Electionr blockchain

In order to start Electionr blockchain run following daemons:

- `electionrd start --home ~/.electionrd/node0`
- `electionrd start --home ~/.electionrd/node1`
- `electionrd start --home ~/.electionrd/node2`

This should result in empty blocks being minted.

### Start REST server

You can start REST server with this command:

- `electionrcli rest-server --chain-id electionr-chain`

## FAQ

1) How can I get cosmos validator public key in different formats?  
Use gaiadebug: `gaiadebug pubkey <pub_key>`

2) How can I get node id needed for persistent peers in node configuration files?  
Execute command: `electionrd tendermint show-node-id --home ~/.electionrd/node0`

3) How can I reset Electionr blockchain?  
Execute these commands:  
`legalerd unsafe-reset-all --home ~/.electionrd/node0`  
`legalerd unsafe-reset-all --home ~/.electionrd/node1`  
`legalerd unsafe-reset-all --home ~/.electionrd/node2`