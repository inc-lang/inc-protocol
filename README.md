# inc-protocol

Protocol Buffer and gRPC service definitions for the INC language runtime.

## Overview

inc-protocol defines the **Bridge** service, which provides remote code execution for the INC language. It supports two execution modes:

- **CallStateless** — single request-response execution with no shared state
- **CallStateful** — bidirectional streaming where state persists across calls

## Prerequisites

- [Go](https://go.dev/) 1.25+
- [Buf](https://buf.build/docs/installation) CLI

## Usage

Add as a Go module dependency:

```sh
go get github.com/inc-lang/inc-protocol
```

Import in your Go code:

```go
import inccore "github.com/inc-lang/inc-protocol/gen/go"
```

## Development

```sh
# Regenerate Go code from inc.proto
make generate

# Lint the proto definition
make lint

# Check for breaking changes against master
make breaking

# Run Go tests
make test

# Run all checks (lint + breaking + test)
make check
```

## Protocol

### Messages

| Message | Description |
|---------|-------------|
| `ExecutionRequest` | Code, input variables, and runtime config |
| `ExecutionResponse` | Return value, stdout/stderr, and exit code |
| `Variable` | Dynamically-typed value (int, float, string, bool, bytes, list, map, null) |
| `ListValue` | Ordered sequence of Variables |
| `MapValue` | String-keyed collection of Variables |

### Example

```protobuf
// Request
{
  request_id: "abc-123",
  code: "return x + y",
  variables: {"x": {int_val: 10}, "y": {int_val: 20}},
  config: {"memory_limit": "128M"}
}

// Response
{
  request_id: "abc-123",
  result: {int_val: 30},
  stdout: "",
  stderr: "",
  exit_code: 0
}
```

## Contributing

1. Edit `inc.proto`
2. Run `make generate` to regenerate Go code
3. Run `make check` to verify everything passes
4. Submit a pull request
