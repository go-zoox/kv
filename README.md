# KV - Key-Value Store

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-zoox/jmap)](https://pkg.go.dev/github.com/go-zoox/jmap)
[![Build Status](https://github.com/go-zoox/jmap/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/go-zoox/jmap/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-zoox/jmap)](https://goreportcard.com/report/github.com/go-zoox/jmap)
[![Coverage Status](https://coveralls.io/repos/github/go-zoox/jmap/badge.svg?branch=master)](https://coveralls.io/github/go-zoox/jmap?branch=master)
[![GitHub issues](https://img.shields.io/github/issues/go-zoox/jmap.svg)](https://github.com/go-zoox/jmap/issues)
[![Release](https://img.shields.io/github/tag/go-zoox/jmap.svg?label=Release)](https://github.com/go-zoox/jmap/tags)

## Installation
To install the package, run:
```bash
go get github.com/go-zoox/jmap
```

## Getting Started

```go
func TestMapGetSet(t *testing.T) {
	m := Map{}
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

## Inspired by
* [srfrog/dict](https://github.com/srfrog/dict) - Python-like dictionaries for Go

## License
GoZoox is released under the [MIT License](./LICENSE).