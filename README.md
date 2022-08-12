# Trade Executor

This is a simple project for order maching engine.

## Architecture

- Exchange API

    Implemented the simple Binance websocket api and mock api.

- Order Executor

    For every order book update, it is trying to match orders from the pool.

- Internal API

    It provides gRPC API and common CLIs.

## CLIs

```bash
# serve the main service including gRPC server
trade-executor serve --cfg <config file path> 

# get the last order book for the specific symbol
trade-executor price --symbol <symbol name>

# apply order through the CLI
trade-executor order apply --symbol <symbol name> --type <order dir> --amount <order amount> --price <order price>

# get the executed resutl
trade-executor order result --id <order id>
```

## Further

- Currently, the matching engine is processed when getting a new order book. It can be a bottleneck when getting massive orders. We can build it as an independent process separating from the order book stream.

- It supports only `Limit Order` right now, we can update it to provide more order kinds like `Market Order`, `Stop Loss`, `Take Profit`, and so on.

- We can implement `FIX` api or stream api using websocket instead of gRPC api.