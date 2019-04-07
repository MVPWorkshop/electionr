# Electionr

Electionr is a custom blockchain implementation written in Golang.  
It is implemented on top of [Tendermint](https://github.com/tendermint/tendermint) by extending [Cosmos SDK](https://github.com/cosmos/cosmos-sdk).

During the first year, validator nodes are elected using Proof of Determination (PoD).  
This process is described in [Legaler Whitepaper](https://github.com/Legaler/Whitepapers/blob/master/Proof%20of%20Determination.md).

### Prerequisites

- `Golang v1.12.1+`

### Installation

- Navigate to your electionr cloned repository
- Install Golang packages: `go mod tidy`
- Install Electionr: `make install`

#### Follow our [guide](./testnet_config/README.md) to setup and start Electionr blockchain.