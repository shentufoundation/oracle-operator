### How to deploy SecurityPrimitive contract on CertiK Chain

1. Revise function `getEndpointUrl` in [SecurityPrimitive.sol](SecurityPrimitive.sol) 
to set your own url pattern.
    
    e.g. `return string(abi.encodePacked(_endpoint, "?address=", contractAddress, "&functionSignature=", functionSignature));`

    Or custom your personal score evaluation method and change the first return value of
    `getInsight` to true, which represents `isUrl`.
    
2. Run `solc` to obtain binary and ABI files for your SecurityPrimitive contract.
```bash
solc SecurityPrimitive.sol --abi --bin -o .
```

3. Deploy your SecurityPrimitive contract on CertiK Chain. 

```bash
certik tx cvm deploy SecurityPrimitive.bin --abi SecurityPrimitive.abi --args <your/person/base/endpoint/url> --from node0 --gas-prices 0.025uctk --gas-adjustment 2.0 --gas auto --chain-id <chainid> -y
```

4. Record your Primitive Contract address `new-contract-address` from screen output.

5. Check your Primitive Contract by query `getInsight` function.
```bash
certikcli tx cvm call <your/primitive/contract/address> "getInsight" "0x00000000000000000000" "0x0100" --from node0 --gas-prices 0.025uctk --gas-adjustment 2.0 --gas auto -y -b block
```

6. Set your contract address in the oracle-operator configuration file (<home>/.certik/config/oracle-operator.toml).
See `primitive` in template at [oracle-operator.toml](../oracle-operator.toml).
