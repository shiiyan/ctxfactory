# ctxfactory

[![CI](https://github.com/shiiyan/ctxfactory/actions/workflows/ci.yaml/badge.svg)](https://github.com/shiiyan/ctxfactory/actions/workflows/ci.yaml)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/shiiyan/ctxfactory.svg)](https://pkg.go.dev/github.com/shiiyan/ctxfactory)


Minimal Factory-style context builder for Go tests.  
Easily set defaults, overrides, and skips for `context.Context`.

## Install

```bash
go get github.com/shiiyan/ctxfactory
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/shiiyan/ctxfactory"
)

type ctxKey string

var userKey = ctxKey("user")

type User struct {
    ID       int
    Name     string
    Location string
}

func main() {
    f := ctxfactory.New(map[any]any{
        userKey:  User{ID: 0, Name: "guest", Location: "unknown"},
        "traceID": "default-trace", // string key is allowed but less safe
    })

    // Default context (uses context.Background())
    ctx1 := f.Build()
    u := ctx1.Value(userKey).(User)

    // Override a default (fluent With then Build)
    ctx2 := f.With(map[any]any{
        userKey: User{ID: 42, Name: "Alice", Location: "SF"},
    }).Build()
    u2 := ctx2.Value(userKey).(User)

    // Skip a default key when building
    ctx3 := f.Skip("traceID").Build()

    // Build using an existing base context
    base := context.WithValue(context.Background(), "requestID", "req-123")
    ctx4 := f.BuildWith(base)

    // Build with timeout (caller must call cancel)
    ctx5, cancel := f.BuildWithTimeout(nil, 200*time.Millisecond)
    defer cancel()
    _ = ctx5

    // Build with cancel
    ctx6, cancel2 := f.BuildWithCancel(nil)
    _ = ctx6
    cancel2()
}
```
