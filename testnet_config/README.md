# Quick Start

This tutorial will help you setup and start Electionr blockchain.  
We will start the blockchain with 2 validators and 1 regular node.

### Initialize nodes configuration files

Copy nodes configurations to your home directory:
- `cp -r testnet_config/config/.electionrd ~`
- `cp -r testnet_config/config/.electionrcli ~`  
(Run both commands from the root of the repo)

### Start Electionr blockchain

You can start Electionr blockchain by running following daemons:

- `electionrd start --home ~/.electionrd/node0`
- `electionrd start --home ~/.electionrd/node1`
- `electionrd start --home ~/.electionrd/node2`

This should result in empty blocks being minted.

### Start REST server

In order to start REST server run these commands:

- `make update_api_docs` (from the root of the repo)
- `electionrcli rest-server --chain-id electionr-chain`

API documentation will be available [here](http://127.0.0.1:1317/swagger-ui/).

#### Next, learn how to [elect a new validator](./elect_validator.md).

## FAQ

1) How to initialize new node?  
`electionrd init <node_name> --home ~/.electionrd/<node_name> --chain-id electionr-chain`

2) How can I get cosmos validator public key in different formats?  
Use gaiadebug: `gaiadebug pubkey <pub_key>`  

3) How can I get node id needed for persistent peers in node configuration files?  
Execute command: `electionrd tendermint show-node-id --home ~/.electionrd/<node_name>`

4) How can I reset Electionr blockchain?  
Execute these commands:  
`legalerd unsafe-reset-all --home ~/.electionrd/node0`  
`legalerd unsafe-reset-all --home ~/.electionrd/node1`  
`legalerd unsafe-reset-all --home ~/.electionrd/node2`

5) What are operator passwords?  
Password for operators 1, 2 and 3 is: `supersifra`

6) How to change operator password?  
Run following command: `electionrcli keys add <operator_name> --recover`  
When prompted enter appropriate mnemonic from [mnemonics directory](./mnemonics). 