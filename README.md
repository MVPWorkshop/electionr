# Electionr

Electionr is created to solve the problem of launching a PoS network and issuing block rewards to validators in a fair and secure way with no pre-mined sum. 

As the starting validators in a PoS network have 100% of the tokens without a permissionless way to nominate validators, then the starting validators may be a colluding cartel or a greedy founder who can dominate the network - especially in the early years of high inflation. Even if a cartel or greedy founder is not present, then massive coordination is required between genesis validators to start the network, which sets up channels of communication which may be exploited later by cartels. If tokens are not issued to non-validators, then non-validators can never enter the validator set which makes the network governance easily compromised.

In order to solve for this during the first year, validator nodes are elected using a process called Proof of Determination (PoD).

Once the minimum number of validators have been elected, they begin the network. New validators are also added in a permissionless way, until the maximum is achieved. PoS tokens are also distributed into a liquidity pool that allows anyone to acquire them. Once the network has been bootstrapped, the election and bridge contracts can be decommissioned safely.

This process was first described and used in [Legaler Whitepaper](https://github.com/Legaler/Whitepapers/blob/master/Proof%20of%20Determination.md).

Electionr is a custom blockchain implementation written in Golang.  
It is implemented on top of [Tendermint](https://github.com/tendermint/tendermint) by extending [Cosmos SDK](https://github.com/cosmos/cosmos-sdk).

### Prerequisites

- `Golang v1.12.1+`

### Installation & Setup

- Navigate to your electionr cloned repository
- Install Golang packages: `go mod tidy`
- Install Electionr: `make install`
- Follow our [guide](./testnet_config/README.md) to setup and start Electionr blockchain with 2 validators and 1 regular node.

### Versioning

- We use [SemVer](https://semver.org) for versioning of all our open source projects. 

### Authors

- [Vuksan Simunović](https://www.linkedin.com/in/vuksan-simunovi%C4%87-bb39286a/), Software Engineer @ [MVP Workshop](https://mvpworkshop.co)
- [Filip Petrović](https://www.linkedin.com/in/filip-petrovi%C4%87-160076129/), Software Engineer @ [MVP Workshop](https://mvpworkshop.co)
- [Tomislav Ranđić](https://www.linkedin.com/in/tomislav-randjic-2601b4b/), Engineering Director @ [MVP Workshop](https://mvpworkshop.co)

### Contributors

- [Mališa Pušonja](https://www.linkedin.com/in/malisapusonja/), CTO @ [MVP Workshop](https://mvpworkshop.co)

### Acknowledgments

- [Legaler](https://www.legaler.com/), for the Proof of Determination (PoD) idea and concept.
