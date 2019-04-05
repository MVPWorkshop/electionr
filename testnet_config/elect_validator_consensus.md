# Second Validator Vote

In order to achieve network consensus, second validator needs to vote as well.  
We just need to repeat the process using second operator keys.

### Generate insert validator elects transaction

Send a request as an operator of the second node:
```
curl http://127.0.0.1:1317/election/validator_elects \
--header "Content-Type: application/json" \
--request POST \
--data @- <<'EOF' > tx1.json
{
    "base_req": {
        "from": "cosmos14l0n0qwhkf7p5uvuraeyvxmks5n8yt3f9vq462",
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
Transaction `tx1.json` is now ready for signing.  
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
                    "initiator_address": "cosmosvaloper14l0n0qwhkf7p5uvuraeyvxmks5n8yt3fqc5qke",
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

`./testnet_config/sign <path_to_tx1.json> cosmos14l0n0qwhkf7p5uvuraeyvxmks5n8yt3f9vq462`

Output represents signature of the second node's operator which will be used to broadcast transaction.  
Now change `"signatures"` field in `tx1.json` by modifying `null` to `[]` and pasting signature inside of brackets.

### Broadcast transaction

Copy `value` field of `tx1.json` and paste it inside `tx` field, then execute following command:

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

Transaction is now broadcasted to other nodes, and it should achieve network consensus.  
Validator elect has now been added successfully.

You can make sure that election cycle has been completed by sending following HTTP GET request:
```
curl http://127.0.0.1:1317/election/cycle/1 --header "Content-Type: application/json"
```