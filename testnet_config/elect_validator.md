# Validator Election

In order to elect a new validator, his election cycle needs to obtain majority of votes (more than 2/3).  
In [Tendermint](https://github.com/tendermint/tendermint) validator votes are counted based on their voting power, 
which is proportional to the amount of tokens they are staking.  
The more validator stakes, higher his voting power.

### Generate insert validator elects transaction

First we need to generate insert validator elects transaction which we can sign in order to broadcast it to Tendermint.  
We can do so by sending HTTP POST request to Electionr REST server.  
Lets first send a request as an operator of the first node:
```
curl http://127.0.0.1:1317/election/validator_elects \
--header "Content-Type: application/json" \
--request POST \
--data @- <<'EOF' > tx.json
{
    "base_req": {
        "from": "cosmos1ge5t2s0k78y4zws2gpv5wsa0dw83fmflcpw0qn",
        "chain_id": "legaler-chain"
    },
    "cycle_number": "1",
    "elected_validators": [
        {
            "operator_addr": "cosmosvaloper1hgvneywkqgm03hvjl4mtj5ee6ymy60guuqyv59",
            "cons_pub_key": "cosmosvalconspub1zcjduepqw29jdgktm22u5fajnwdatf6acfa9af4g9lc6svm6an5ltjns4vwqjjqflu",
            "place": "1"
        }
    ]
}
EOF
```
Transaction `tx.json` is now ready for signing.  
It should look something like this:
```
{
    "type": "auth/StdTx",
    "value": {
        "msg": [
            {
                "type": "electionr/MsgInsertValidatorElects",
                "value": {
                    "elected_validators": [
                        {
                            "operator_addr": "cosmosvaloper1hgvneywkqgm03hvjl4mtj5ee6ymy60guuqyv59",
                            "cons_pub_key": "cosmosvalconspub1zcjduepqw29jdgktm22u5fajnwdatf6acfa9af4g9lc6svm6an5ltjns4vwqjjqflu",
                            "place": "1"
                        }
                    ],
                    "initiator_address": "cosmosvaloper1ge5t2s0k78y4zws2gpv5wsa0dw83fmfla466vq",
                    "cycle_number": "1"
                }
            }
        ],
        "fee": {
            "amount": null,
            "gas": "200000"
        },
        "signatures": null,
        "memo": ""
    }
}
```

### Sign the transaction

In order to sign transaction do the following (from the root of the repo):

`./testnet_config/sign <path_to_tx.json> cosmos1ge5t2s0k78y4zws2gpv5wsa0dw83fmflcpw0qn`

Output represents signature of the first node's operator which will be used to broadcast transaction.  
Now change `"signatures"` field in `tx.json` by modifying `null` to `[]` and pasting signature inside of brackets.

### Broadcast transaction

Copy `value` field of `tx.json` and paste it inside `tx` field, then execute following command:

```
curl http://127.0.0.1:1317/txs \
--header "Content-Type: application/json" \
--request POST \
--data @- <<'EOF'
{
    "tx": <insert_value>,
    "return": "block"
}
EOF
```

Transaction is now broadcasted to other nodes. However, since this validator doesn't have enough power to achieve 
majority [second validator needs to vote](./elect_validator_consensus.md) as well, in order to reach network consensus.