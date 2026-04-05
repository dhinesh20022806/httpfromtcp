# HTTP From TCP

An HTTP/1.1 server built from scratch in Go, parsing raw TCP connections into HTTP requests.

## Overview

This project implements an HTTP server at the TCP level, manually parsing the HTTP protocol rather than relying on Go's `net/http` package. It reads raw bytes from TCP connections and parses them into structured HTTP requests including request lines and headers.

## Project Structure

```
├── cmd/tcplistener/     # Main entry point — TCP server listening on port 42069
├── internal/
│   ├── request/         # HTTP request line parser (method, target, version)
│   └── headers/         # HTTP header parser
├── go.mod
└── message.txt
```

## Features

- Raw TCP listener on port `:42069`
- HTTP/1.1 request line parsing (method, target, version)
- Header field parsing with validation
- Incremental streaming parser using a state machine

## Getting Started

### Prerequisites

- Go 1.22+

### Run the server

```sh
go run cmd/tcplistener/main.go
```

### Send a test request

```sh
curl http://localhost:42069/
```

### Run tests

```sh
go test ./...
```
