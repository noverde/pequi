# Pequi

URL shortener with multi database support.

## Supported Databases

- [x] In Memory
- [ ] Redis (open a Pull Request)
- [x] Firestore
- [ ] DynamoDB (open a Pull Request)
- [ ] MySQL (open a Pull Request)
- [ ] Postgres (open a Pull Request)
- [ ] SQLite (open a Pull Request)

![Build](https://github.com/noverde/pequi/workflows/Build/badge.svg?branch=master&event=push)
![Test](https://github.com/noverde/pequi/workflows/Test/badge.svg?branch=master&event=push)
[![codecov](https://codecov.io/gh/noverde/pequi/branch/master/graph/badge.svg)](https://codecov.io/gh/noverde/pequi)
[![Go Report Card](https://goreportcard.com/badge/github.com/noverde/pequi)](https://goreportcard.com/badge/github.com/noverde/pequi)
[![Docker pulls](https://img.shields.io/docker/pulls/ecanuto/pequi.svg)](https://hub.docker.com/r/ecanuto/pequi/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## setup

```
go mod tidy
```

## run

**var envs:**
- HTTP_PORT _default: 8080_
- FIRESTORE_PROJECT
- FIRESTORE_COLLECTION

```
go run main.go
```
