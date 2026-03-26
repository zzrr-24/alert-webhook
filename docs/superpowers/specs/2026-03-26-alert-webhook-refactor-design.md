# Alert Webhook 重构设计文档

## 项目概述
将现有的单文件 Alert Webhook 重构为符合 Go 标准项目结构的分层架构，提升代码可维护性、可测试性和扩展性。

## 目标
- 代码组织和可维护性：将大文件拆分成职责清晰的小模块
- 为扩展功能做准备：良好的架构支持未来功能添加
- 符合 Go 最佳实践：遵循 Standard Go Project Layout

## 技术选型
- **HTTP 框架**：`gin-gonic/gin`（高性能 Web 框架）
- **日志库**：`uber-go/zap`（高性能结构化日志）
- **配置库**：`spf13/viper`（简化环境变量和文件读取）
- **测试**：标准库 `testing` + `testify/assert`（断言库）

## 架构设计

### 整体架构
项目采用分层架构，从上到下依次为：
- **Handler 层**：处理 HTTP 请求，参数解析，调用 Service 层
- **Service 层**：核心业务逻辑，包括飞书客户端调用、卡片构建
- **Model 层**：数据结构定义（Alertmanager、Feishu 相关）
- **Config 层**：配置管理，支持环境变量和 YAML 文件

各层之间通过接口通信，降低耦合，便于测试。

### 目录结构

```
alert-webhook/
├── cmd/
│   └── webhook/
│       └── main.go                 # 程序入口：初始化配置、日志、路由、启动服务
├── internal/
│   ├── handler/                   # HTTP 处理层
│   │   ├── handler.go             # Alert 处理器，定义接口和方法
│   │   └── handler_test.go        # 单元测试（mock Service）
│   ├── service/                   # 业务逻辑层
│   │   ├── feishu.go              # 飞书客户端：发送消息、重试逻辑
│   │   ├── feishu_test.go         # 单元测试（mock HTTP）
│   │   ├── card.go                # 卡片构建：formatTime, severityEmoji, buildCard
│   │   └── card_test.go           # 单元测试
│   ├── model/                     # 数据模型
│   │   └── alert.go               # AlertmanagerWebhook 结构定义
│   └── config/                    # 配置管理
│       ├── config.go              # Config 结构定义（Server, Feishu, Log）
│       ├── loader.go              # 加载配置：读取 YAML，覆盖环境变量
│       └── loader_test.go         # 单元测试
├── configs/
│   └── config.yaml.example        # 配置文件模板（不提交实际配置）
├── pkg/
│   └── logger/                    # 结构化日志封装（基于 zap）
│       ├── logger.go              # 初始化 logger，提供 Info/Error 等方法
│       └── logger_test.go         # 单元测试
├── go.mod
├── go.sum
├── Dockerfile
├── Makefile                       # 构建、测试、运行命令
├── AGENTS.md
└── .gitignore
```

## 数据流设计

### 请求处理流程
1. **接收 Alertmanager Webhook**
   - Alertmanager → Gin Handler (`POST /alert`)
   - Handler 使用 `ShouldBindJSON` 解析请求体为 `AlertmanagerWebhook` model

2. **构建飞书卡片**
   - Handler 调用 Service 层 `CardService`
   - `CardService` 将告警数据转换为飞书交互式卡片

3. **发送飞书消息**
   - Handler 调用 `FeishuService`
   - `FeishuService` 使用 HTTP 客户端发送到飞书 webhook
   - 记录请求/响应日志，处理错误

4. **返回响应**
   - 成功返回 `{"status": "ok"}`
   - 失败返回相应的 HTTP 状态码和错误信息

## 配置设计

### 配置项示例

```yaml
server:
  port: 8080
  mode: debug  # debug/release

feishu:
  webhook_url: "https://open.feishu.cn/open-apis/bot/v2/hook/xxx"
  timeout: 10s

log:
  level: info  # debug/info/warn/error
  format: json  # json/console
```

### 环境变量覆盖规则
- `FEISHU_WEBHOOK_URL` 覆盖 `feishu.webhook_url`
- `SERVER_PORT` 覆盖 `server.port`
- `LOG_LEVEL` 覆盖 `log.level`

## 错误处理策略

### 分层错误处理
- **Handler 层**：捕获 panic，返回 HTTP 错误响应（使用 Gin 的 recovery 中间件）
- **Service 层**：返回业务错误，由 Handler 处理
- **日志记录**：所有错误都记录结构化日志
- **用户响应**：统一错误格式，避免敏感信息泄露

### 错误响应格式
```json
{
  "error": "错误描述",
  "code": "错误码"
}
```

## 测试策略

### 单元测试
- **handler**：测试 HTTP 请求处理，mock Service 层
- **service**：测试业务逻辑，mock HTTP 客户端
- **config**：测试配置加载和环境变量覆盖
- **logger**：测试日志初始化和输出

### 测试覆盖率目标
- 核心业务逻辑（service 层）：≥ 80%
- Handler 层：≥ 70%
- 其他模块：≥ 60%

## 依赖管理

### 第三方依赖
- `github.com/gin-gonic/gin` - Web 框架
- `go.uber.org/zap` - 日志库
- `github.com/spf13/viper` - 配置库
- `github.com/stretchr/testify` - 测试断言库

### 最小依赖原则
尽可能使用标准库，减少外部依赖。

## 实施计划
本设计文档完成后，将创建详细的实施计划，包括：
1. 目录结构创建
2. 依赖安装
3. 模块迁移顺序（main → config → logger → model → service → handler）
4. 每个步骤的验证方法
