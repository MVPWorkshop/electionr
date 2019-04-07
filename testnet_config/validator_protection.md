# Validator Protection

Whenever a cycle gains a majority vote, it's validator elects gain 30 day protection period.  
During this time these newly elected validators cannot be removed from validator set.  
If a validator who is still in his protection period decides to leave however, he is no longer protected.

In order to make sure that validator protection works we can perform a simple test - 
we can try to add a new validator with higher stake. 
If the newly elected validator wasn't protected he would have been removed from validator set.

### Check whether validator protection works

Currently we have 3 validators:
- operator0 who has staked 5000000stakecoin
- operator1 with 4000000stakecoin staked
- operator2 (our newly elected validator) with 1000000stakecoin staked

Maximum number of validators is 3, so if we were to stake 2000000stakecoin with operator3, 
by Cosmos rules, he should replace operator2. However operator2 is being protected by Electionr rules, 
and validator set will remain the same.

First lets see current validators with the following command:  
`electionrcli query staking validators --chain-id electionr-chain --home ~/.electionrd/node0`  
You should get a list of operators0, 1 and 2 and all of them should be **_bonded_**, 
which means that they are in the validator set.

Now lets try to make operator3 a validator by executing this command:
  
```
electionrcli tx staking create-validator --from operator3 --amount 2000000stakecoin --moniker node3 \
--pubkey cosmosvalconspub1zcjduepqw29jdgktm22u5fajnwdatf6acfa9af4v8ze03hvjrmnelnsudwhqcu2mqp \
--commission-max-change-rate 0 --commission-max-rate 0 --commission-rate 0 --min-self-delegation 1 \
--chain-id electionr-chain
```

Get the list of all validators one more time:  
`electionrcli query staking validators --chain-id electionr-chain --home ~/.electionrd/node0`  
You should now see that operator2 is still **_bonded_** although his stake is lower than the stake of operator3, 
who remains **_unbonded_** and out of validator set.

## FAQ

1) Does validator protection period starts from the moment his cycle gains a majority vote, 
or from the moment he decides to join with `create-validator` command?  
Validator protection period starts from the moment the cycle he is in gains a majority vote.