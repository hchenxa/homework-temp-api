# CSGHub Lite API Test Suite

An API integration test suite for [CSGHub Lite](https://github.com/OpenCSGs/csghub-lite) REST API, built with [Ginkgo](https://github.com/onsi/ginkgo) and [Gomega](https://github.com/onsi/gomega).

## Coverage

Covers all 9 endpoints across 3 categories, totaling **20 test cases**:

| Category | Endpoint | Cases |
|----------|----------|-------|
| Inference | `POST /api/chat` | 4 |
| Inference | `POST /api/generate` | 4 |
| Model Management | `GET /api/tags` | 3 |
| Model Management | `POST /api/show` | 2 |
| Model Management | `POST /api/pull` | 1 |
| Model Management | `DELETE /api/delete` | 1 |
| Service Management | `GET /api/health` | 2 |
| Service Management | `GET /api/ps` | 2 |
| Service Management | `POST /api/stop` | 1 |

## Prerequisites

- Go 1.26+
- A running CSGHub Lite server (default `http://localhost:11435`)

## Quick Start

```bash
# Install dependencies
go mod tidy

# Run all tests
go test ./pkg/tests/ -v

# Or use the ginkgo CLI
go run github.com/onsi/ginkgo/v2/ginkgo ./pkg/tests/
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `API_BASE_URL` | `http://localhost:11435` | API server address |
| `TEST_MODEL` | `Qwen/Qwen3-0.6B-GGUF` | Model name used in tests |

```bash
# Custom address and model
API_BASE_URL=http://192.168.1.100:11435 TEST_MODEL=qwen2.5:0.5b go test ./pkg/tests/ -v
```

## Label Filtering

Every spec is tagged with Ginkgo labels for selective execution:

```bash
# Run only inference tests
ginkgo --label-filter "inference" ./pkg/tests/

# Run only streaming tests
ginkgo --label-filter "streaming" ./pkg/tests/

# Run tests that don't need a model (health, tags, ps)
ginkgo --label-filter "!requires-model" ./pkg/tests/

# Run only service management tests
ginkgo --label-filter "service-management" ./pkg/tests/

# Combined filter: inference + non-streaming
ginkgo --label-filter "inference && non-streaming" ./pkg/tests/
```

Available labels:

- **Endpoint**: `health`, `tags`, `chat`, `generate`, `show`, `pull`, `delete`, `ps`, `stop`
- **Category**: `inference`, `model-management`, `service-management`
- **Mode**: `streaming`, `non-streaming`
- **HTTP Method**: `get`, `post`, `delete`
- **Dependency**: `requires-model`

## Project Structure

```
├── go.mod
├── go.sum
└── pkg/
    ├── utils/
    │   └── client.go       # HTTP client, request/response types, config
    └── tests/
        ├── suite_test.go              # Ginkgo entry point + BeforeSuite
        ├── health_test.go             # GET /api/health
        ├── tags_test.go               # GET /api/tags
        ├── inference_test.go          # POST /api/chat and /api/generate
        ├── model_management_test.go   # POST /api/show, /api/pull, DELETE /api/delete
        └── service_management_test.go # GET /api/ps, POST /api/stop
```

## Adding New Tests

1. Create a new `_test.go` file in `pkg/tests/`
2. Use `Describe` + `It` to define specs, and `utils.Get/Post/Delete` to make requests
3. Add `Label` for filtering
4. Register new request/response types in `pkg/utils/client.go` if needed
