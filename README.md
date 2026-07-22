# Demo: Using mini-protoc & mini-rpc in Your Project

This is a standalone demo project that shows how to use **mini-protoc** (protobuf compiler) and **mini-rpc** (RPC runtime) together in any Go project to build a gRPC-like service — without actually using gRPC.

## What This Demo Builds

A simple **Todo Service** with two RPCs:

- `CreateTodo` — creates a todo with a title and description, returns it with an auto-generated ID
- `GetTodo` — fetches a todo by ID

## Prerequisites

- Go 1.25.3+

That's it. No cloning required — everything installs via `go get` and `go install`.

## Step-by-Step Guide

### 1. Install mini-protoc

```sh
go install github.com/akansha204/mini-protoc/cmd/mini-protoc@v1.0.0
```

This gives you the `mini-protoc` CLI globally.

### 2. Add mini-rpc to Your Project

```sh
go get github.com/akansha204/mini-rpc@v1.0.0
```

### 3. Define Your Schema

Create a `.proto` file describing your messages and service:

```protobuf
// todo.proto
syntax = "proto3";

package todo;

message TodoRequest {
    string title = 1;
    string description = 2;
    int32 id = 3;
}

message TodoResponse {
    int32 id = 1;
    string title = 2;
    string description = 3;
    bool completed = 4;
}

service TodoService {
    rpc CreateTodo(TodoRequest) returns (TodoResponse);
    rpc GetTodo(TodoRequest) returns (TodoResponse);
}
```

### 4. Generate Go Code

```sh
mini-protoc ./todo.proto
```

This generates `todo.pb.go` containing:

- `TodoRequest` and `TodoResponse` structs
- `TodoService` interface (the contract your implementation must satisfy)
- `TodoServiceClient` struct and constructor (placeholder — not wired to the network yet)
- `RegisterTodoService(server, service)` — the glue that connects your implementation to the RPC runtime

### 5. Implement the Service

Write a file that satisfies the generated `TodoService` interface:

```go
// server/service.go
package main

import (
    "fmt"
    "sync"
    todo "github.com/akansha204/demo-in-go"
)

type TodoServiceImpl struct {
    mu     sync.Mutex
    todos  map[int32]todo.TodoResponse
    nextID int32
}

func NewTodoServiceImpl() *TodoServiceImpl {
    return &TodoServiceImpl{
        todos:  make(map[int32]todo.TodoResponse),
        nextID: 1,
    }
}

func (s *TodoServiceImpl) CreateTodo(req todo.TodoRequest) (todo.TodoResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    t := todo.TodoResponse{
        Id:          s.nextID,
        Title:       req.Title,
        Description: req.Description,
        Completed:   false,
    }
    s.todos[t.Id] = t
    s.nextID++
    return t, nil
}

func (s *TodoServiceImpl) GetTodo(req todo.TodoRequest) (todo.TodoResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    t, ok := s.todos[req.Id]
    if !ok {
        return todo.TodoResponse{}, fmt.Errorf("todo %d not found", req.Id)
    }
    return t, nil
}
```

### 6. Wire Up the Server

```go
// server/main.go
package main

import (
    "log"
    todo "github.com/akansha204/demo-in-go"
    "github.com/akansha204/mini-rpc/rpc"
)

func main() {
    server := rpc.NewDefaultServer()

    svc := NewTodoServiceImpl()
    todo.RegisterTodoService(server, svc)

    log.Println("Todo server listening on :8080")
    if err := server.Serve(":8080"); err != nil {
        log.Fatal(err)
    }
}
```

What happens here:
1. `rpc.NewDefaultServer()` creates an RPC server with JSON codec and empty registry
2. `todo.RegisterTodoService(server, svc)` registers each RPC method (like `TodoService/CreateTodo`) into the server's registry with decode/call/encode handlers
3. `server.Serve(":8080")` starts the TCP listener and handles connections

### 7. Write the Client

```go
// client/client.go
package main

import (
    "fmt"
    "log"
    todo "github.com/akansha204/demo-in-go"
    "github.com/akansha204/mini-rpc/rpc"
)

func main() {
    client, err := rpc.NewDefaultClient(":8080")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create a todo
    req := todo.TodoRequest{
        Title:       "Learn mini-gRPC",
        Description: "Build a demo using mini-protoc and mini-rpc",
    }
    var resp todo.TodoResponse
    if err := client.Call("TodoService/CreateTodo", req, &resp); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created: id=%d title=%q\n", resp.Id, resp.Title)

    // Fetch it back
    getReq := todo.TodoRequest{Id: resp.Id}
    var getResp todo.TodoResponse
    if err := client.Call("TodoService/GetTodo", getReq, &getResp); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Fetched: id=%d title=%q description=%q\n",
        getResp.Id, getResp.Title, getResp.Description)
}
```

### 8. Run It

Terminal 1 — start the server:

```sh
go run ./server/
```

Terminal 2 — run the client:

```sh
go run ./client/
```

Expected output:

```
[client] created: id=1 title="Learn mini-gRPC" completed=false
[client] fetched: id=1 title="Learn mini-gRPC" description="Build a demo using mini-protoc and mini-rpc"
```

## Project Structure

```
demo-in-go/
├── go.mod          # depends on github.com/akansha204/mini-rpc v1.0.0
├── go.sum
├── todo.proto      # your schema
├── todo.pb.go      # generated by mini-protoc (do not edit manually)
├── server/
│   ├── main.go     # server entry point — 10 lines
│   └── service.go  # your TodoService implementation
└── client/
    └── client.go   # client that calls the RPC methods
```

## go.mod

```go
module github.com/akansha204/demo-in-go

go 1.25.3

require github.com/akansha204/mini-rpc v1.0.0
```

Only `mini-rpc` is a runtime dependency. `mini-protoc` is a build tool — you install it once globally, not per-project.

## What mini-protoc Generates

For each `.proto` file, mini-protoc generates:

| Artifact | What It Is |
|----------|-----------|
| Message structs | Go structs with exported fields matching your proto message fields |
| Service interface | An interface your implementation must satisfy |
| Client struct + constructor | Placeholder client (methods currently panic — use `rpc.Client` directly) |
| Register function | The bridge: takes your implementation, decodes payloads, calls your methods, encodes responses |

## How the Request Flows

```
Client                          Server
  │                               │
  ├─ Call("TodoService/CreateTodo", req, &resp)
  │   │                           │
  │   ├─ encode req → JSON bytes  │
  │   ├─ wrap in protocol.Request │
  │   ├─ encode envelope → bytes  │
  │   ├─ TCP send (4-byte frame)  │
  │   │──────────────────────────>│
  │   │                           ├─ TCP receive frame
  │   │                           ├─ decode protocol.Request
  │   │                           ├─ lookup "TodoService/CreateTodo" in registry
  │   │                           ├─ decode payload → TodoRequest
  │   │                           ├─ call service.CreateTodo(req)
  │   │                           ├─ encode TodoResponse → bytes
  │   │                           ├─ wrap in protocol.Response
  │   │                           ├─ encode envelope → bytes
  │   │<──────────────────────────┤
  │   ├─ TCP receive frame        │
  │   ├─ decode protocol.Response │
  │   ├─ decode payload → resp    │
  │   │                           │
  │   └─ return nil               │
```

## Quick Reference

| Task | Command |
|------|---------|
| Install compiler | `go install github.com/akansha204/mini-protoc/cmd/mini-protoc@v1.0.0` |
| Add runtime to project | `go get github.com/akansha204/mini-rpc@v1.0.0` |
| Generate code | `mini-protoc ./your-file.proto` |
| Start server | `go run ./server/` |
| Run client | `go run ./client/` |

## The Only Import You Need

Everything lives behind one package:

```go
import "github.com/akansha204/mini-rpc/rpc"
```

The generated code also imports this package internally for `server.Encode`, `server.Decode`, and the `rpc.Server` type in the `Register` function signature. You never need to touch `internal/codec`, `internal/transport`, or `internal/protocol` directly.
