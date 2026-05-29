# CSGHub Lite API Test Suite

基于 [Ginkgo](https://github.com/onsi/ginkgo) 和 [Gomega](https://github.com/onsi/gomega) 的 API 集成测试套件，用于测试 [CSGHub Lite](https://github.com/OpenCSGs/csghub-lite) 的 REST API。

## 测试覆盖

覆盖 API 文档中所有 3 类共 9 个端点，共 **20 个测试用例**：

| 类别 | 端点 | 用例数 |
|------|------|--------|
| 推理 | `POST /api/chat` | 4 |
| 推理 | `POST /api/generate` | 4 |
| 模型管理 | `GET /api/tags` | 3 |
| 模型管理 | `POST /api/show` | 2 |
| 模型管理 | `POST /api/pull` | 1 |
| 模型管理 | `DELETE /api/delete` | 1 |
| 服务管理 | `GET /api/health` | 2 |
| 服务管理 | `GET /api/ps` | 2 |
| 服务管理 | `POST /api/stop` | 1 |

## 前置条件

- Go 1.26+
- 运行中的 CSGHub Lite 服务（默认 `http://localhost:11435`）

## 快速开始

```bash
# 安装依赖
go mod tidy

# 运行所有测试
go test ./pkg/tests/ -v

# 或使用 ginkgo CLI
go run github.com/onsi/ginkgo/v2/ginkgo ./pkg/tests/
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `API_BASE_URL` | `http://localhost:11435` | API 服务地址 |
| `TEST_MODEL` | `Qwen/Qwen3-0.6B-GGUF` | 测试使用的模型名称 |

```bash
# 指定自定义地址和模型
API_BASE_URL=http://192.168.1.100:11435 TEST_MODEL=qwen2.5:0.5b go test ./pkg/tests/ -v
```

## 按标签筛选

每个测试用例都标记了 Ginkgo Label，可按类别、端点、模式等维度筛选运行：

```bash
# 只跑推理相关用例
ginkgo --label-filter "inference" ./pkg/tests/

# 只跑流式响应测试
ginkgo --label-filter "streaming" ./pkg/tests/

# 只跑不需要模型的（健康检查、tags、ps）
ginkgo --label-filter "!requires-model" ./pkg/tests/

# 只跑服务管理类
ginkgo --label-filter "service-management" ./pkg/tests/

# 组合筛选：推理 + 非流式
ginkgo --label-filter "inference && non-streaming" ./pkg/tests/
```

可用标签列表：

- **端点**: `health`, `tags`, `chat`, `generate`, `show`, `pull`, `delete`, `ps`, `stop`
- **分类**: `inference`, `model-management`, `service-management`
- **模式**: `streaming`, `non-streaming`
- **HTTP 方法**: `get`, `post`, `delete`
- **模型依赖**: `requires-model`

## 项目结构

```
├── go.mod
├── go.sum
└── pkg/
    ├── utils/
    │   └── client.go       # HTTP 客户端、请求/响应类型、配置
    └── tests/
        ├── suite_test.go              # Ginkgo 入口 + BeforeSuite
        ├── health_test.go             # GET /api/health
        ├── tags_test.go               # GET /api/tags
        ├── inference_test.go          # POST /api/chat 和 /api/generate
        ├── model_management_test.go   # POST /api/show, /api/pull, DELETE /api/delete
        └── service_management_test.go # GET /api/ps, POST /api/stop
```

## 新增测试

1. 在 `pkg/tests/` 下创建新的 `_test.go` 文件
2. 使用 `Describe` + `It` 定义用例，通过 `utils.Get/Post/Delete` 发起请求
3. 添加 `Label` 以便筛选
4. 首次运行前在 `pkg/utils/client.go` 中注册新的请求/响应类型
