# HTTP from TCP in Go

An educational project to build an HTTP server from scratch using raw TCP sockets in Go.

## Goal

> The purpose of this project is purely educational. The main goal is to learn how HTTP works on a lower level, without relying on high-level libraries that abstract away the protocol details.

## HTTP Version

This implementation is based on the specification for **HTTP/1.1**. While newer versions like HTTP/2 and HTTP/3 have different on-the-wire implementations (e.g., binary framing instead of plaintext), they largely share the same core semantics (methods, status codes, headers) defined in the HTTP Semantics RFC.

## Features

- [x] **Request-line parsing**: `GET /path HTTP/1.1`
- [x] **Header parsing**: `Key: Value` pairs.

## Roadmap

- [ ] Body parsing
- [ ] Response generation
- [ ] Handling different HTTP methods (POST, PUT, DELETE, etc.)
- [ ] Concurrent connection handling

## How to Run

1.  Clone the repository.
2.  Run the TCP listener:
    ```bash
    go run ./cmd/tcplistener/main.go
    ```
3.  Send an HTTP request to `localhost:42069`. You can use `curl` or `telnet`.
    ```bash
    curl -v localhost:42069
    ```
4.  As of this moment the parsed ouput of the curl command is displayed in the terminal.

## Relevant RFCs

For those interested in the technical details, this project loosely follows the specifications laid out in:

- [RFC 9112: HTTP/1.1](https://www.rfc-editor.org/rfc/rfc9112.html)
- [RFC 9110: HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110.html)
