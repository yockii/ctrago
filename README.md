# ctrago

[中文文档（README_zh.md）](./README_zh.md)

A Go client library for cTrader OpenAPI, supporting both WebSocket and TCP communication. This library allows you to interact with cTrader's trading API, including authentication, account management, and trading operations.

## Features
- Connect to cTrader OpenAPI via WebSocket or TCP
- Application and account authentication
- Query account list and details
- Token refresh support
- Modular design for account operations (orders, symbols, traders)

## Installation

```
go get github.com/yockii/ctrago
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/yockii/ctrago"
)

func main() {
    client, err := ctrago.NewClient(
        false, // isLive: false for demo, true for live
        "<clientId>",
        "<clientSecret>",
        "<accessToken>",
        30*time.Second, // heartbeat interval
    )
    if err != nil {
        panic(err)
    }
    ctx := context.Background()
    if authResp, err := client.ApplicationAuth(ctx); err != nil {
        panic(err)
    } else {
        fmt.Printf("authResp: %s\n", authResp.String())
    }
    if accountList, err := client.GetAccountList(ctx); err != nil {
        panic(err)
    } else {
        fmt.Printf("accountList: %s\n", accountList.String())
    }
}
```

## Protobuf Code Generation
See [protobuf.md](./protobuf.md) for instructions on generating Go code from proto files.

## Dependencies
- [google.golang.org/protobuf](https://pkg.go.dev/google.golang.org/protobuf)
- [github.com/gorilla/websocket](https://pkg.go.dev/github.com/gorilla/websocket)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[Specify your license here]
