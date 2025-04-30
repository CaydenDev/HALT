# HALT - Human Authentication Layer Tool

HALT is a lightweight reverse proxy that implements a cryptographic challenge mechanism to detect and block AI traffic. It works by presenting mathematical challenges that are easy for humans to solve but difficult for AI systems to process.

## Features

- Lightweight HTTP reverse proxy
- Simple mathematical challenge-based authentication
- Session management using secure cookies
- Configurable target service

## Requirements

- Go 1.21 or higher

## Installation

```bash
go get github.com/cayde/halt
```

## Usage

1. Start your target service (the service you want to protect)
2. Run HALT proxy:

```bash
go run main.go
```

By default, HALT will:
- Listen on port 3000
- Forward requests to `http://localhost:8080`

## How It Works

1. When a user first accesses the protected service, they are presented with a mathematical challenge
2. Upon solving the challenge correctly, a secure cookie is set
3. Subsequent requests with a valid cookie are forwarded to the target service
4. Challenges expire after 5 minutes for security

## Security Features

- Secure cookie settings (HTTPOnly, Secure, SameSite)
- Randomized challenge generation
- Time-based challenge expiration
- Single-use challenges

## Configuration

The following parameters can be modified in `main.go`:
- Target service URL
- Proxy port
- Challenge timeout duration
- Cookie settings

## License

MIT