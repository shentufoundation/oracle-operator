# Oracle Operator

Oracle Operator listens to the `create_task` event from Shentu Chain, queries the primitive APIs and pushes the aggregated result back to Shentu Chain.

## How to Config and Run

1. Register the operator on Shentu Chain (through CLI or RESTful API) and lock a certain amount of `CTK`.
  ```bash
  $ shentud tx oracle create-operator <account address> <collateral> --name <operator name> --from <account> --fees 5000uctk --chain-id <chainid> -y -b block
  ```
2. Create the oracle operator configuration file in `shentud` home (default `.shentud/config/oracle-operator.toml`). See template at [oracle-operator.toml](oracle-operator.toml):
  - `type`: Aggregation type, e.g. `linear`. Check [Strategy](STRATEGY.md).
  - `primitive_type`: security primitive type.
  - `weight`: the weight of the result from the corresponding primitive to the final result.
3. Run the oracle operator by the following command.
  ```bash
  $ oracle-operator start --home ~/.shentud --log_level debug --from <account> --chain-id <chainid>
  ```

A sample shell script of running Oracle Operator:

```bash
shentud tx oracle create-operator $(shentud keys show alice -a) 100000uctk --from alice --fees 5000uctk --chain-id yulei-4 -y -b block
oracle-operator start --home ~/.shentud --log_level debug --from alice --chain-id yulei-4
```

## Support of Multiple Client Chain

Contract addresses in security oracle tasks are prefixed with the chain identifier, e.g. `eth:0xabc`(Ethereum), `bsc:0xdef`(Binance Smart Chain). To enable oracle operator handle tasks for multiple chains, the configuration file can be specified as:

```toml
[strategy.eth]
type = "linear"
[[stragety.eth.primitive]]
primitive_type = "whitelist"
weight = 0.1
[[stragety.eth.primitive]]
primitive_type = "bytecode"
weight = 0.1

[strategy.bsc]
type = "linear"
[[stragety.bsc.primitive]]
primitive_type = "sourcecode"
weight = 0.1
[[stragety.bsc.primitive]]
primitive_type = "blacklist"
weight = 0.1
```

## Modules

### Chain `x/oracle`

The `x/oracle` module in chain handles operator registry and task management (create, response, delete, etc...)
