# ZerologixHomework

# Zerologix coding assignment
Imagine you have a trade engine that accepts orders via the protocol (or triggering)
you defined. An order request at least has this information (buy or sell, quantity,
market price or limit price).
The engine matches buy and sell orders that have the same price. Orders have the
same price determined by their timestamp (FIFO). Pending orders queue up in your
system until they are filled or killed. Your system will have multiple traders executing
orders at the same time.

## What is expected?
- SOLID design accommodates requirements to change in the future.
- Testable, compilable and runnable code.
- Concurrent and thread-safe considered.
- Document your design which can be any human-readable form. For example,
README file.
- A git repository can be accessed publicly.
- Golang are preferred but not required.

# 啟動
啟動 docker
```
make up-test
```

啟動程式
```
go run . server
```

1. 到 swagger http://localhost:1234/docs/index.html#/Order/post_trade_order_test
2. 呼叫多次執行測試
3. 到 redis ui http://localhost:38081/
4. 會看到有 trade:LastOrderID 存在

執行邏輯是如果配對成功就不會生成新的 order id
所以執行後的 order id 如果小於執行測試次數是正常的

# 測試指令
```
make test
```

# 文檔

## 買/賣訂單API

### 架構圖

```mermaid
sequenceDiagram
    participant 使用者
    participant API
    participant Redis

    activate 使用者 
    使用者->>API: 賣訂單
        activate API
        API->>Redis: 存取資料
        Redis->>API: 存取資料
        API->>使用者: 完成
        deactivate API
    deactivate 使用者
```

### 流程圖

```mermaid
flowchart TB
    流程開始([流程開始])
    讀取redis訂單資料[/讀取redis訂單資料/]
    是否配對完成{是否配對完成}
    修改redis剩餘訂單[/修改redis剩餘訂單/]
    存入redis此筆訂單[/存入redis此筆訂單/]
    完成([完成])

    流程開始 --> 讀取redis訂單資料
    
    讀取redis訂單資料 --> 配對訂單
    配對訂單 --> 是否配對完成

    是否配對完成 -->|是| 修改redis剩餘訂單
    是否配對完成 -->|否| 存入redis此筆訂單
    
    修改redis剩餘訂單 --> 完成
    存入redis此筆訂單 --> 完成
```