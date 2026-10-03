# chat-mute

禁言：违规次数递增时长、封顶、解除与请求幂等，只用 Go 标准库。

```
go run ./cmd/chat-mute --sample=escalate
go run ./cmd/chat-mute --sample=cap
go run ./cmd/chat-mute --sample=lift
go run ./cmd/chat-mute --sample=idem
go run ./cmd/chat-mute --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **递增**：第 n 次违规的禁言时长是 `baseMs * 2^(n-1)`。
- **封顶**：时长不超过 `capMs`。
- **到期**：当前时刻不小于禁言结束时刻就不再禁言。
- **解除**：`Lift` 清掉当前禁言（但保留违规次数）；没有禁言时返回 false。
- **幂等**：同一个 `reqID` 只记一次违规，重复调用不增加次数也不延长时长。
- **不变量**：违规次数只增不减；禁言时长不超过封顶；同一串违规序列重复执行结果相同。
- **代价**：`Muted` 不许遍历全部玩家，`Scanned` 不随玩家数线性增长。
- **规模**：10 万玩家、每秒 10 万次查询，单次摊还 O(1)。

## 输出契约（不改格式）

```
third=..
fifth=..
muted=..
second=..
scanned=..
```

## 目录

```
mute.go              递增、封顶、解除与幂等
mute_test.go         单元测试
cmd/chat-mute/       命令行入口
```
