# HALT - Human Authentication Layer Tool

HALT is a lightweight reverse proxy that implements a cryptographic challenge mechanism to detect and block AI traffic. It works by presenting mathematical challenges that are easy for humans to solve but difficult for AI systems to process.

## How It Works

1. When a user first accesses the protected service, they are presented with a mathematical challenge
2. Upon solving the challenge correctly, a secure cookie is set
3. Subsequent requests with a valid cookie are forwarded to the target service
4. Challenges expire after 5 minutes for security