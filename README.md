# KV - Key-Value Store

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-zoox/kv)](https://pkg.go.dev/github.com/go-zoox/kv)
[![Build Status](https://github.com/go-zoox/kv/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/go-zoox/kv/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-zoox/kv)](https://goreportcard.com/report/github.com/go-zoox/kv)
[![Coverage Status](https://coveralls.io/repos/github/go-zoox/kv/badge.svg?branch=master)](https://coveralls.io/github/go-zoox/kv?branch=master)
[![GitHub issues](https://img.shields.io/github/issues/go-zoox/kv.svg)](https://github.com/go-zoox/kv/issues)
[![Release](https://img.shields.io/github/tag/go-zoox/kv.svg?label=Release)](https://github.com/go-zoox/kv/tags)

## Installation
To install the package, run:
```bash
go get github.com/go-zoox/kv
```

## Getting Started

```go
func TestMemoryKV(t *testing.T) {
	m := kv.NewMemory()
	if m.Size() != 0 {
		t.Errorf("Expected size 0, got %d", m.Size())
	}

	m.Set("key", "value")
	if m.Get("key") != "value" {
		t.Error("Expected value to be 'value'")
	}

	if m.Size() != 1 {
		t.Errorf("Expected size 1, got %d", m.Size())
	}
}
```

## Engines
* [x] Memory
* [x] Redis
* [ ] MongoDB
* [ ] SQLite
* [ ] PostgreSQL
* [ ] MySQL
* [ ] DynamoDB

## Inspired by
* [srfrog/dict](https://github.com/srfrog/dict) - Python-like dictionaries for Go

## License
GoZoox is released under the [MIT License](./LICENSE).