# Electionr

Electionr is created to solve the problem of launching a PoS network in a fair and secure way with no pre-mined sum, and block rewards to be issued to validators. As the starting validators have 100% of the tokens without a permissionless way to nominate validators, then the starting validators may be a colluding cartel or a greedy founder who can dominate the network - especially in the early years of high inflation. Even if a cartel or greedy founder is not present, then massive coordination is required between genesis validators to start the network, which sets up channels of communication which may be exploited later by cartels. If tokens are not issued to non-validators, then non-validators can never enter the validator set. Such a network is doomed.

In order to solve for this during the first year, validator nodes are elected using Proof of Determination (PoD).

Once the minimum number of validators have been elected, they begin the network. New validators are also added in a permissionless way, until the maximum is achieved. PoS tokens are also distributed into a liquidity pool that allows anyone to acquire them. Once the network has been bootstrapped, the election and bridge contracts can be decommissioned safely.

This process was first described and used in [Legaler Whitepaper](https://github.com/Legaler/Whitepapers/blob/master/Proof%20of%20Determination.md).

Electionr is a custom blockchain implementation written in Golang.  
It is implemented on top of [Tendermint](https://github.com/tendermint/tendermint) by extending [Cosmos SDK](https://github.com/cosmos/cosmos-sdk).

### Prerequisites

- `Golang v1.12.1+`

### Installation

- Navigate to your electionr cloned repository
- Install Golang packages: `go mod tidy`
- Install Electionr: `make install`

### Follow our [guide](./testnet_config/README.md) to setup and start Electionr blockchain.

### Versioning

- We use [SemVer](https://semver.org) for versioning of all our open source projects. 

### Authors

- Vuksan Simunović (https://mvpworkshop.co)
- Filip Petrović (https://mvpworkshop.co)
- Tomislav Ranđić (https://mvpworkshop.co)

### Contributors

- Mališa Pušonja (https://mvpworkshop.co)

### Acknowledgments

- Legaler (https://www.legaler.com/), for the Proof of Determination (PoD) concept.
