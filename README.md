# WIRE 
Fast,minimalist web framework for [go](https://golang.com)

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/AhmedAshraf780/wire/internals/wire"
)

type responseType struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}
type User struct {
	Name  string `json:"name"`
	Email string `json:"email""`
}

func main() {
	app := wire.NewApplication()
	wire.USE(app, func(req *wire.Request[wire.EmptyBody], resp *wire.Response[responseType]) error {
		token , ok := req.Headers["token"]
		if !ok {
			return resp.Write(http.StatusUnauthorized,responseType{
				Ok : false,
				Message: "Missing token",
            })
        }
		req.Context["token"] = token
		return wire.ErrNext
	})
	wire.POST(app, "/users/:id", func(req *wire.Request[User], resp *wire.Response[responseType]) error {
		name := req.Body.Name
		email := req.Body.Email
		token := req.Context["token"]
		msg := fmt.Sprintf("Hello %s your email is %s and token is %s",name,email,token)
		return resp.Write(http.StatusOK, responseType{Ok: true, Message: msg})
	})
	app.Listen(3000)
}
```

# Installation
Before installing, [download and install go](https://go.dev/dl/)

Installation is done using `go get` command:

```shell
go get github.com/AhmedAshraf780/wire
```

# Features
- Easy to setup
- Robust routing 
- Focus on high performance
- Specialized for HTTP APIs


# Docs

## Table of Contents

- [Getting Started](#getting-started)
- [Creating an Application](#creating-an-application)
- [Registering Routes](#registering-routes)
- [Handlers](#handlers)
- [Request Body & EmptyBody](#request-body--emptybody)
- [Path Parameters](#path-parameters)
- [Query Parameters](#query-parameters)
- [Middleware](#middleware)
  - [Global Middleware](#global-middleware)
  - [Route-Specific Middleware](#route-specific-middleware)
- [ErrNext](#errext)
- [Context](#context)
- [Response](#response)
- [Setting Headers](#setting-headers)
- [Full Example](#full-example)

---

## Getting Started

Wire is a minimalist, generic HTTP API framework for Go. It uses a raw TCP server (not `net/http`) and Go generics to give you strongly typed request bodies and response bodies.

### Installation

```shell
go get github.com/AhmedAshraf780/wire
```

---

## Creating an Application

Every Wire app starts by creating an `Application` instance:

```go
import "github.com/AhmedAshraf780/wire/internals/wire"

app := wire.NewApplication()
```

Then start the server on a given port:

```go
app.Listen(3000)
```

---

## Registering Routes

Wire provides five functions for registering routes, one per HTTP method:

```go
wire.GET(app, path, callbacks...)
wire.POST(app, path, callbacks...)
wire.PUT(app, path, callbacks...)
wire.PATCH(app, path, callbacks...)
wire.DELETE(app, path, callbacks...)
```

Each function takes:
1. `app` -- the application instance
2. `path` -- the route path (supports `:param` segments)
3. `callbacks` -- a variadic list of handler/middleware functions. The **last** callback is the final handler. Any preceding callbacks are treated as route-specific middlewares.

### Static vs Dynamic Routes

- **Static routes** have no parameters (e.g., `/users`, `/about`). They are matched in O(1) via map lookup.
- **Dynamic routes** contain `:param` segments (e.g., `/users/:id`). They are matched by iterating segments at runtime.

---

## Handlers

Every handler is a function with this signature:

```go
func(req *wire.Request[TReq], resp *wire.Response[TRes]) error
```

Where `TReq` is the Go type for the request body (parsed from JSON), and `TRes` is the Go type for the response body.

### Returning a Response

Call `resp.Write(statusCode, body)` to set the status code and response body, then return `nil`:

```go
func getUser(req *wire.Request[wire.EmptyBody], resp *wire.Response[UserResponse]) error {
    return resp.Write(200, UserResponse{Ok: true, Name: "Alice"})
}
```

### Error Handling

If the handler encounters an error, it can write an error response directly and return `nil`, or return the error:

```go
func getUser(req *wire.Request[wire.EmptyBody], resp *wire.Response[UserResponse]) error {
    id, ok := req.Params["id"]
    if !ok {
        return resp.Write(400, UserResponse{Ok: false, Message: "Missing id"})
    }
    // ... use id
    return resp.Write(200, UserResponse{Ok: true, Name: id})
}
```

---

## Request Body & EmptyBody

Wire uses Go generics to automatically parse the JSON request body into the type you specify as `TReq`.

### When Your Handler Needs a Body

If your handler expects a JSON body (e.g., `POST`, `PUT`, `PATCH`), specify the body type directly:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUser(req *wire.Request[CreateUserRequest], resp *wire.Response[ResponseMsg]) error {
    name := req.Body.Name
    email := req.Body.Email
    return resp.Write(201, ResponseMsg{Ok: true, Message: "Created " + name})
}
```

### When Your Handler Has No Body (use `wire.EmptyBody`)

**If your handler does not parse a request body, you MUST use `wire.EmptyBody` as the `TReq` type parameter.** This tells Wire to skip JSON parsing for that handler.

`wire.EmptyBody` is an empty struct:

```go
type EmptyBody struct{}
```

You must use it explicitly in two scenarios:

**1. GET / DELETE handlers (no request body expected):**

```go
func getUsers(req *wire.Request[wire.EmptyBody], resp *wire.Response[UserList]) error {
    return resp.Write(200, UserList{Users: []string{"Alice", "Bob"}})
}
```

**2. Middleware functions (always use `wire.EmptyBody`):**

```go
func authMiddleware(req *wire.Request[wire.EmptyBody], resp *wire.Response[SomeResponse]) error {
    token := req.Headers["token"]
    if token == "" {
        return resp.Write(401, SomeResponse{Ok: false, Message: "Unauthorized"})
    }
    req.Context["token"] = token
    return wire.ErrNext
}
```

**Why is this required?** Internally, Wire checks if `TReq` is `wire.EmptyBody` using a type assertion. If it is, JSON parsing is skipped entirely. If you use any other type without a valid JSON body, Wire will return a `400 Bad Request` error.

### Summary

| Scenario | `TReq` Type |
|---|---|
| Handler expects a JSON body | Your struct type (e.g., `CreateUserRequest`) |
| Handler does **not** expect a body (GET, DELETE) | `wire.EmptyBody` |
| Any middleware | `wire.EmptyBody` |

---

## Path Parameters

Define path parameters with a `:` prefix:

```go
wire.GET(app, "/users/:id", func(req *wire.Request[wire.EmptyBody], resp *wire.Response[UserResponse]) error {
    id := req.Params["id"]
    return resp.Write(200, UserResponse{Ok: true, Name: id})
})
```

You can have multiple parameters:

```go
wire.GET(app, "/users/:id/posts/:postId", handler)
// req.Params["id"]     -> "42"
// req.Params["postId"] -> "7"
```

---

## Query Parameters

Query parameters are parsed from the URL automatically. Access them via `req.Query`:

```
GET /search?q=wire&page=2
```

```go
func search(req *wire.Request[wire.EmptyBody], resp *wire.Response[SearchResult]) error {
    q := req.Query["q"]       // "wire"
    page := req.Query["page"] // "2"
    return resp.Write(200, SearchResult{Query: q, Page: page})
}
```

---

## Middleware

Wire supports two kinds of middleware: **global** and **route-specific**.

### Global Middleware

Global middleware runs on **every incoming request**, before any route handler. Register it with `wire.USE`:

```go
wire.USE(app, func(req *wire.Request[wire.EmptyBody], resp *wire.Response[SomeResponse]) error {
    // This runs for ALL routes
    return wire.ErrNext
})
```

### Route-Specific Middleware

Route-specific middleware runs only for a particular route. Pass them as extra arguments before the final handler:

```go
wire.GET(app, "/admin/dashboard", authMiddleware, roleCheck, dashboardHandler)
```

Execution order: `authMiddleware` -> `roleCheck` -> `dashboardHandler`.

### Middleware Rules

1. **Every middleware MUST use `wire.EmptyBody` as its `TReq` type.** Middleware does not consume a request body.

2. **Every middleware MUST return `wire.ErrNext`** to pass control to the next handler/middleware in the chain. If it returns anything else (e.g., after calling `resp.Write`), the chain stops.

3. You can write a response directly from a middleware to short-circuit the chain (e.g., for authentication failures):

```go
func authMiddleware(req *wire.Request[wire.EmptyBody], resp *wire.Response[AuthResponse]) error {
    token := req.Headers["token"]
    if token == "" {
        return resp.Write(401, AuthResponse{Ok: false, Message: "Missing token"})
        // Chain stops here -- handler is never called
    }
    req.Context["token"] = token
    return wire.ErrNext
    // Chain continues to next middleware/handler
}
```

---

## ErrNext

`wire.ErrNext` is a sentinel error that signals "continue to the next handler in the chain":

```go
var ErrNext = errors.New("next")
```

### When to Return `wire.ErrNext`

- **In every middleware** that wants to pass control forward. The middleware has done its work (e.g., authentication, logging) and wants the request to continue.

### When NOT to Return `wire.ErrNext`

- **In the final handler** -- the handler should call `resp.Write(...)` and return `nil`.
- **In any middleware that wants to stop the chain** -- write a response with `resp.Write(...)` and return its result (or return `nil` after writing).

### How It Works Internally

When a middleware returns `wire.ErrNext`:
1. Any changes to `req.Context` are propagated to the next handler.
2. The framework continues to the next callback in the chain.

When a middleware returns anything else:
1. The chain stops. No further handlers are called.
2. The response written via `resp.Write()` (if any) is sent to the client.

---

## Context

`req.Context` is a `map[string]interface{}` that lets middleware pass data to handlers:

```go
// In middleware
req.Context["userID"] = "12345"
return wire.ErrNext

// In handler
userID := req.Context["userID"].(string)
```

Context values are **propagated** forward through the chain when a middleware returns `wire.ErrNext`. This is the standard way to share data between middleware and handlers (e.g., authenticated user info, parsed tokens).

---

## Response

The `Response[T]` struct gives you three methods:

```go
resp.Write(statusCode int, body T) error   // Set status code and body
resp.SetHeader(key string, value string)    // Set a custom response header
resp.Header(key string) string              // Read a response header
```

### Setting Custom Headers

```go
func handler(req *wire.Request[wire.EmptyBody], resp *wire.Response[Data]) error {
    resp.SetHeader("X-Custom-Header", "hello")
    return resp.Write(200, Data{Value: "ok"})
}
```

### Default Headers

Wire automatically sets these headers on every response:
- `Content-Type: application/json`
- `Connection: keep-alive`
- `Content-Length`
- `Server: Wire/1.0`

---

## Full Example

```go
package main

import (
    "net/http"
    "github.com/AhmedAshraf780/wire/internals/wire"
)

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type ResponseMsg struct {
    Ok      bool   `json:"ok"`
    Message string `json:"message"`
}

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    app := wire.NewApplication()

    // Global middleware -- runs on every route
    // MUST use wire.EmptyBody and return wire.ErrNext
    wire.USE(app, func(req *wire.Request[wire.EmptyBody], resp *wire.Response[ResponseMsg]) error {
        token, ok := req.Headers["token"]
        if !ok {
            return resp.Write(http.StatusUnauthorized, ResponseMsg{
                Ok: false, Message: "Missing token",
            })
        }
        req.Context["token"] = token
        return wire.ErrNext
    })

    // GET route with dynamic param -- no body, so use wire.EmptyBody
    wire.GET(app, "/users/:id", func(req *wire.Request[wire.EmptyBody], resp *wire.Response[User]) error {
        id := req.Params["id"]
        token := req.Context["token"].(string)
        return resp.Write(200, User{Name: id, Email: token + "@example.com"})
    })

    // POST route with JSON body
    wire.POST(app, "/users", func(req *wire.Request[CreateUserRequest], resp *wire.Response[ResponseMsg]) error {
        name := req.Body.Name
        return resp.Write(201, ResponseMsg{Ok: true, Message: "Created user: " + name})
    })

    // Route with middleware chain
    wire.GET(app, "/users/:id/profile", authMiddleware, getProfile)

    app.Listen(3000)
}

// Middleware -- always EmptyBody, always returns ErrNext (or writes a response to stop)
func authMiddleware(req *wire.Request[wire.EmptyBody], resp *wire.Response[ResponseMsg]) error {
    if req.Headers["token"] == "" {
        return resp.Write(401, ResponseMsg{Ok: false, Message: "Unauthorized"})
    }
    return wire.ErrNext
}

// Final handler
func getProfile(req *wire.Request[wire.EmptyBody], resp *wire.Response[User]) error {
    id := req.Params["id"]
    return resp.Write(200, User{Name: id, Email: id + "@example.com"})
}
```

