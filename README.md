# TradeDemo

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
