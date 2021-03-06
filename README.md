# Course for gRPC usage in Go

## Motivation

Learn about gRPC, protocol buffers and improve Go skills.

## About the project

This project is a part of [course](https://www.udemy.com/course/grpc-golang/). Project contains 1 service which exposes 4 different RPC calls:

- Unary call - simply says hi.
- Server streaming call - says hi for 10 times.
- Client streaming call - receives hi for many times.
- BiDi streaming call - receives hi many times and says hi many times.

Also, in this project more advanced topics investigated:

- gRPC errors and error handling.
- Deadlines.
- Reflection.

## Local development

### gRPC Server

- Start Go gRPC server.

```bash
  make run-server
```

### gRPC Client

- Start server command mentioned in **gRPC Server** section or run it using docker.

```bash
  make run-server-docker
```

- Start Go gRPC client.

```bash
  make run-client
```
