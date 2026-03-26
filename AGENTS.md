# AGENTS.md - Alert Webhook 开发指南

## 项目概述
Prometheus Alertmanager 到飞书的通知转发服务，接收告警并格式化为飞书卡片消息推送。

## 构建和运行

```bash
# 本地构建
go mod init alert-webhook
go mod tidy
go build -o alert-webhook ./cmd/webhook

# Docker 构建
docker build -t alert-webhook .

# 本地运行
./alert-webhook

# Docker 运行
docker run -d -p 8080:8080 alert-webhook
```

## 测试

```bash
# 当前无测试文件，建议添加：
# _test.go 文件，使用标准 testing 包

# 运行所有测试
go test ./...

# 运行单个测试
go test -run TestFunctionName

# 带覆盖率
go test -cover ./...
```

## 代码风格指南

### 命名约定
- **包名**: 小写单数（`main`, `config`, `internal`）
- **结构体**: PascalCase（`AlertmanagerWebhook`）
- **函数**: PascalCase（导出）、camelCase（私有）
- **变量**: camelCase（`feishuWebhook`）
- **常量**: PascalCase 或 UPPER_SNAKE_CASE
- **接口**: PascalCase，以 `er` 结尾（如有）

### 导入规范
- 标准库在前，第三方库在后
- 按字母排序分组
- 每组间空行分隔

```go
import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)
```

### 结构体标签
- JSON 标签使用 snake_case：`json:"alert_name"`
- 必填字段使用 omitempty：`json:"status,omitempty"`

### 错误处理
```go
// 标准 if err != nil 模式
data, err := io.ReadAll(r.Body)
if err != nil {
    log.Println("read body error:", err)
    http.Error(w, err.Error(), 400)
    return
}

// defer 关闭资源
defer resp.Body.Close()
```

### 日志规范
- 使用 `log` 包记录操作日志
- 关键节点使用 `====` 分隔标识
- 包含时间戳和状态信息
- 使用中文说明
- 可在日志中适当使用 emoji 增强可读性

```go
log.Println("==== Received Alertmanager Webhook ====")
log.Println("FEISHU_WEBHOOK:", feishuWebhook)
log.Println("Feishu adapter running on :8080")
```

### 时间处理
```go
// 使用 RFC3339 标准格式
t.In(loc).Format("2006-01-02 15:04:05")

// 时区处理
loc, _ := time.LoadLocation("Asia/Shanghai")
```

### 类型使用
- JSON 解析使用 `map[string]interface{}` 或定义结构体
- HTTP 请求使用标准库 `net/http`
- 编码/解码使用 `encoding/json`

### 函数组织
- 小型函数，单一职责
- 辅助函数如 `formatTime`, `severityEmoji`, `statusEmoji`
- 核心逻辑函数如 `buildCard`, `handleAlert`

### 代码注释
- 关键逻辑使用中文注释
- 导出函数应有文档注释
- 复杂业务逻辑需说明

```go
// 构建飞书卡片消息
func buildCard(p AlertmanagerWebhook) map[string]interface{} {
```

## 项目结构
```
alert-webhook/
├── cmd/           # 命令行程序（如有多入口）
├── config/        # 配置管理
├── internal/      # 内部包
│   └── utils/     # 工具函数
├── main.go        # 主入口
└── Dockerfile     # 容器化配置
```

## 开发注意事项
- 端口：8080
- 环境变量：`FEISHU_WEBHOOK`（当前硬编码，建议改为环境变量）
- 端点：`/alert`（POST 接收 Alertmanager webhook）
- 输入格式：Alertmanager webhook JSON
- 输出格式：飞书交互式卡片

## 代码质量工具（建议添加）
```bash
# gofmt - 自动格式化
gofmt -s -w .

# golint - 代码风格检查
golangci-lint run

# go vet - 静态分析
go vet ./...
```
