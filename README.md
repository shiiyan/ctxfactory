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
// use github.com/shiiyan/ctxfactory

factory := ctxfactory.NewContextFactory(map[interface{}]interface{}{
  "userID":  0,
  "traceID": "default-trace",
  "isAdmin": false,
})

// Default context
ctx1 := factory.With(nil, nil)

// Override userID
ctx2 := factory.With(map[interface{}]interface{}{"userID": 42}, nil)

// Skip traceID
ctx3 := factory.With(map[interface{}]interface{}{"isAdmin": true}, []interface{}{"traceID"})
```
