# ctrago

一个用于 cTrader OpenAPI 的 Go 客户端库，支持 WebSocket 和 TCP 通信。该库可用于与 cTrader 交易 API 进行交互，包括鉴权、账户管理和交易操作。

## 功能特性
- 通过 WebSocket 或 TCP 连接 cTrader OpenAPI
- 应用和账户鉴权
- 查询账户列表及详情
- 支持刷新 Token
- 账户操作模块化（订单、品种、交易员等）

## 安装方法

```
go get github.com/yockii/ctrago
```

## 使用示例

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
        false, // isLive: 演示环境为 false，实盘为 true
        "<clientId>",
        "<clientSecret>",
        "<accessToken>",
        30*time.Second, // 心跳间隔
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

## Protobuf 代码生成
请参考 [protobuf.md](./protobuf.md) 了解如何根据 proto 文件生成 Go 代码。

## 依赖
- [google.golang.org/protobuf](https://pkg.go.dev/google.golang.org/protobuf)
- [github.com/gorilla/websocket](https://pkg.go.dev/github.com/gorilla/websocket)

## 贡献
欢迎提交 Pull Request。如有重大更改，请先提交 Issue 进行讨论。

## 许可证
[请在此处注明您的许可证]
